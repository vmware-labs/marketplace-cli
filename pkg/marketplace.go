// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

//go:generate counterfeiter . HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Marketplace struct {
	Host          string
	APIHost       string
	StorageBucket string
	StorageRegion string
	Client        HTTPClient
}

var (
	ProductionConfig *Marketplace
	StagingConfig    *Marketplace
)

func init() {
	ProductionConfig = &Marketplace{
		Host:          "gtw.marketplace.cloud.vmware.com",
		APIHost:       "api.marketplace.cloud.vmware.com",
		StorageBucket: "cspmarketplaceprd",
		StorageRegion: "us-west-2",
		Client:        http.DefaultClient,
	}
	StagingConfig = &Marketplace{
		Host:          "gtwstg.market.csp.vmware.com",
		APIHost:       "apistg.market.csp.vmware.com",
		StorageBucket: "cspmarketplacestage",
		StorageRegion: "us-east-2",
		Client:        http.DefaultClient,
	}
}

func (m *Marketplace) MakeURL(path string, params url.Values) *url.URL {
	return &url.URL{
		Scheme:   "https",
		Host:     m.Host,
		Path:     path,
		RawQuery: params.Encode(),
	}
}

func (m *Marketplace) Get(requestURL *url.URL) (*http.Response, error) {
	return m.SendRequest("GET", requestURL, map[string]string{}, nil)
}

func (m *Marketplace) Post(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": contentType,
	}
	return m.SendRequest("POST", requestURL, headers, content)
}

func (m *Marketplace) Put(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": contentType,
	}
	return m.SendRequest("PUT", requestURL, headers, content)
}

func (m *Marketplace) SendRequest(method string, requestURL *url.URL, headers map[string]string, content io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, requestURL.String(), content)
	if err != nil {
		return nil, fmt.Errorf("failed to build %s request: %w", requestURL.String(), err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("csp-auth-token", viper.GetString("csp.refresh-token"))

	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("marketplace request failed: %w", err)
	}

	return resp, nil
}
