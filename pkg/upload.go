// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

type CredentialsResponse struct {
	AccessID     string    `json:"accessId"`
	AccessKey    string    `json:"accessKey"`
	SessionToken string    `json:"sessionToken"`
	Expiration   time.Time `json:"expiration"`
}

func (m *Marketplace) GetUploadCredentials() (*CredentialsResponse, error) {
	requestURL := m.MakeURL("/aws/credentials/generate", nil)
	requestURL.Host = m.APIHost
	response, err := m.Get(requestURL)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch credentials: %d", response.StatusCode)
	}

	credsResponse := &CredentialsResponse{}
	d := json.NewDecoder(response.Body)
	if m.strictDecoding {
		d.DisallowUnknownFields()
	}
	err = d.Decode(credsResponse)
	if err != nil {
		return nil, err
	}

	return credsResponse, nil
}

func (m *Marketplace) GetUploader(orgID, hashAlgorithm string, credentials aws.Credentials) internal.Uploader {
	if m.uploader == nil {
		return internal.NewS3Uploader(m.StorageBucket, m.StorageRegion, hashAlgorithm, orgID, credentials)
	}
	return m.uploader
}
