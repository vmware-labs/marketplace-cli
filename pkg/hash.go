// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func Hash(filePath, hashAlgorithm string) (string, error) {
	var hashAlgo hash.Hash
	if hashAlgorithm == models.HashAlgoSHA1 {
		hashAlgo = sha1.New()
	} else if hashAlgorithm == models.HashAlgoSHA256 {
		hashAlgo = sha256.New()
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", filePath, err)
	}

	_, err = io.Copy(hashAlgo, file)
	if err != nil {
		return "", fmt.Errorf("failed to generate the hash for %s: %w", filePath, err)
	}

	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close the file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hashAlgo.Sum(nil)), nil
}
