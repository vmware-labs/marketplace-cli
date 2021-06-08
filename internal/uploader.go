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
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func MakeUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s-%d%s", base, makeTimestamp(), ext)
}

type Uploader interface {
	Upload(filePath string) (string, string, error)
}

const (
	HashAlgoSHA1   = "SHA1"
	HashAlgoSHA256 = "SAH256"
)

type S3Uploader struct {
	region        string
	hashAlgorithm hash.Hash
	credentials   aws.Credentials
	orgId         string
}

func NewS3Uploader(region, hashAlgorithm, orgId string, credentials aws.Credentials) *S3Uploader {
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
		orgId:         orgId,
	}
}

func (u *S3Uploader) Upload(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open %s: %w", filePath, err)
	}

	_, err = io.Copy(u.hashAlgorithm, file)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate the hash for %s: %w", filePath, err)
	}

	s3Config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: u.credentials,
		}),
		config.WithRegion(u.region),
	)
	if err != nil {
		return "", "", err
	}

	filename := MakeUniqueFilename(filepath.Base(filePath))
	client := s3.NewFromConfig(s3Config)
	input := &s3.PutObjectInput{
		Bucket: aws.String("cspmarketplacestage"),
		Key:    aws.String(path.Join(u.orgId, "marketplace-product-files", filename)),
		Body:   file,
	}

	_, err = client.PutObject(context.Background(), input)
	if err != nil {
		return "", "", err
	}

	err = file.Close()
	if err != nil {
		return "", "", fmt.Errorf("failed to close file: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", *input.Bucket, u.region, *input.Key)
	fileHash := hex.EncodeToString(u.hashAlgorithm.Sum(nil))
	return fileURL, fileHash, nil
}
