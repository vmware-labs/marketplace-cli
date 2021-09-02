// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type CredentialsResponse struct {
	AccessID     string    `json:"accessId"`
	AccessKey    string    `json:"accessKey"`
	SessionToken string    `json:"sessionToken"`
	Expiration   time.Time `json:"expiration"`
}

func (m *Marketplace) GetUploadCredentials() (*CredentialsResponse, error) {
	requestURL := m.MakeURL("/aws/credentials/generate", url.Values{})
	requestURL.Host = m.APIHost
	response, err := m.Get(requestURL)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch credentials: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	credsResponse := &CredentialsResponse{}
	err = json.Unmarshal(body, credsResponse)
	if err != nil {
		return nil, err
	}

	return credsResponse, nil
}
