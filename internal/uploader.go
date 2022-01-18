// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func MakeUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s-%d%s", base, makeTimestamp(), ext)
}

func makeUniqueFileID() string {
	return fmt.Sprintf("fileuploader%d.url", makeTimestamp())
}

//go:generate counterfeiter . Uploader
type Uploader interface {
	Hash(filePath string) (string, string, error)
	Upload(filePath string) (string, error)
	UploadFile(filePath string) (*models.ProductDeploymentFile, error)
}

type S3Uploader struct {
	bucket        string
	region        string
	hash          hash.Hash
	hashAlgorithm string
	credentials   aws.Credentials
	orgID         string
}

func NewS3Uploader(bucket, region, hashAlgorithm, orgID string, credentials aws.Credentials) *S3Uploader {
	var hashAlgo hash.Hash
	if hashAlgorithm == models.HashAlgoSHA1 {
		hashAlgo = sha1.New()
	} else if hashAlgorithm == models.HashAlgoSHA256 {
		hashAlgo = sha256.New()
	}

	return &S3Uploader{
		bucket:        bucket,
		region:        region,
		hash:          hashAlgo,
		hashAlgorithm: hashAlgorithm,
		credentials:   credentials,
		orgID:         orgID,
	}
}

func (u *S3Uploader) makeKey(filename string) string {
	return path.Join(u.orgID, "marketplace-product-files", filename)
}

func (u *S3Uploader) makeURL(filename string) string {
	key := u.makeKey(filename)
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)
}

func (u *S3Uploader) Hash(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open %s: %w", filePath, err)
	}

	_, err = io.Copy(u.hash, file)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate the hash for %s: %w", filePath, err)
	}

	err = file.Close()
	if err != nil {
		return "", "", fmt.Errorf("failed to close the file %s: %w", filePath, err)
	}

	return hex.EncodeToString(u.hash.Sum(nil)), u.hashAlgorithm, nil
}

func (u *S3Uploader) Upload(filePath string) (string, error) {
	s3Config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: u.credentials,
		}),
		config.WithRegion(u.region),
	)
	if err != nil {
		return "", err
	}

	filename := MakeUniqueFilename(filepath.Base(filePath))
	client := s3.NewFromConfig(s3Config)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	stat, _ := file.Stat()

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(u.bucket),
		Key:           aws.String(u.makeKey(filename)),
		Body:          file,
		ContentLength: stat.Size(),
	})
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close file: %w", err)
	}

	return u.makeURL(filename), nil
}

func (u *S3Uploader) UploadFile(filePath string) (*models.ProductDeploymentFile, error) {
	hashString, hashAlgorithm, err := u.Hash(filePath)
	if err != nil {
		return nil, err
	}

	s3Config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: u.credentials,
		}),
		config.WithRegion(u.region),
	)
	if err != nil {
		return nil, err
	}

	filename := MakeUniqueFilename(filepath.Base(filePath))
	client := s3.NewFromConfig(s3Config)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get info for %s: %w", filePath, err)
	}

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(u.bucket),
		Key:           aws.String(u.makeKey(filename)),
		Body:          file,
		ContentLength: stat.Size(),
	})
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close file: %w", err)
	}

	return &models.ProductDeploymentFile{
		Name:          filename,
		Url:           u.makeURL(filename),
		HashAlgo:      hashAlgorithm,
		HashDigest:    hashString,
		IsRedirectUrl: false,
		UniqueFileID:  makeUniqueFileID(),
		VersionList:   []string{},
	}, nil
}
