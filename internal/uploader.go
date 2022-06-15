// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	FolderMediaFiles   = "media-files"
	FolderMetaFiles    = "meta-files"
	FolderProductFiles = "marketplace-product-files"
)

func now() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

func MakeUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s-%s%s", base, now(), ext)
}

//go:generate counterfeiter . Uploader
type Uploader interface {
	UploadMediaFile(filePath string) (string, string, error)
	UploadMetaFile(filePath string) (string, string, error)
	UploadProductFile(filePath string) (string, string, error)
}

type S3Uploader struct {
	bucket string
	region string
	orgID  string
	client S3Client
	output io.Writer
}

//go:generate counterfeiter . S3Client
type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func NewS3Client(region string, creds aws.Credentials) S3Client {
	s3Config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: creds,
		}),
		config.WithRegion(region),
	)
	if err != nil {
		return nil
	}

	return s3.NewFromConfig(s3Config)
}

func NewS3Uploader(bucket, region, orgID string, client S3Client, output io.Writer) *S3Uploader {
	return &S3Uploader{
		bucket: bucket,
		orgID:  orgID,
		region: region,
		client: client,
		output: output,
	}
}

func (u *S3Uploader) UploadMediaFile(filePath string) (string, string, error) {
	filename := filepath.Base(filePath)
	key := path.Join(u.orgID, FolderMediaFiles, now(), filename)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)
	err := u.upload(filePath, key, types.ObjectCannedACLPublicRead)
	return filename, url, err
}

func (u *S3Uploader) UploadMetaFile(filePath string) (string, string, error) {
	filename := filepath.Base(filePath)
	datestamp := now()
	key := path.Join(u.orgID, FolderMetaFiles, datestamp, filename)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)
	err := u.upload(filePath, key, types.ObjectCannedACLPrivate)
	return filename, url, err
}

func (u *S3Uploader) UploadProductFile(filePath string) (string, string, error) {
	filename := MakeUniqueFilename(filepath.Base(filePath))
	key := path.Join(u.orgID, FolderProductFiles, filename)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)
	err := u.upload(filePath, key, types.ObjectCannedACLPrivate)
	return filename, url, err
}

func (u *S3Uploader) upload(filePath, key string, acl types.ObjectCannedACL) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get info for %s: %w", filePath, err)
	}

	progressBar := MakeProgressBar(fmt.Sprintf("Uploading %s", path.Base(file.Name())), stat.Size(), u.output)
	_, err = u.client.PutObject(context.Background(), &s3.PutObjectInput{
		ACL:           acl,
		Bucket:        aws.String(u.bucket),
		Key:           aws.String(key),
		Body:          progressBar.WrapReader(file),
		ContentLength: stat.Size(),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	return nil
}
