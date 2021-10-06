// Copyright 2021 VMware, Inc.
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

type Uploader interface {
	Upload(filePath string) (*models.ProductDeploymentFile, error)
}

const (
	HashAlgoSHA1   = "SHA1"
	HashAlgoSHA256 = "SAH256"
)

type S3Uploader struct {
	region        string
	hashAlgorithm hash.Hash
	credentials   aws.Credentials
	orgID         string
}

func NewS3Uploader(region, hashAlgorithm, orgID string, credentials aws.Credentials) *S3Uploader {
	var hashAlgo hash.Hash
	if hashAlgorithm == HashAlgoSHA1 {
		hashAlgo = sha1.New()
	} else if hashAlgorithm == HashAlgoSHA256 {
		hashAlgo = sha256.New()
	}

	return &S3Uploader{
		region:        region,
		hashAlgorithm: hashAlgo,
		credentials:   credentials,
		orgID:         orgID,
	}
}

func (u *S3Uploader) Upload(bucket, filePath string) (*models.ProductDeploymentFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	stat, _ := file.Stat()

	_, err = io.Copy(u.hashAlgorithm, file)
	if err != nil {
		return nil, fmt.Errorf("failed to generate the hash for %s: %w", filePath, err)
	}
	_, _ = file.Seek(0, io.SeekStart)

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

	key := path.Join(u.orgID, "marketplace-product-files", filename)
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
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
		Url:           fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, u.region, key),
		HashAlgo:      models.HashAlgoSHA1,
		HashDigest:    hex.EncodeToString(u.hashAlgorithm.Sum(nil)),
		IsRedirectUrl: false,
		UniqueFileID:  makeUniqueFileID(),
		VersionList:   []string{},
	}, nil
}
