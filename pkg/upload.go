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

func (c *CredentialsResponse) AWSCredentials() aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     c.AccessID,
		SecretAccessKey: c.AccessKey,
		SessionToken:    c.SessionToken,
		Expires:         c.Expiration,
	}
}

func (m *Marketplace) GetUploadCredentials() (*CredentialsResponse, error) {
	requestURL := MakeURL(m.GetAPIHost(), "/aws/credentials/generate", nil)
	response, err := m.Client.Get(requestURL)
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

func (m *Marketplace) GetUploader(orgID string) (internal.Uploader, error) {
	if m.uploader == nil {
		credentials, err := m.GetUploadCredentials()
		if err != nil {
			return nil, fmt.Errorf("failed to get upload credentials: %w", err)
		}
		client := internal.NewS3Client(m.StorageRegion, credentials.AWSCredentials())
		return internal.NewS3Uploader(m.StorageBucket, m.StorageRegion, orgID, client, m.Output), nil
	}
	return m.uploader, nil
}

func (m *Marketplace) SetUploader(uploader internal.Uploader) {
	m.uploader = uploader
}
