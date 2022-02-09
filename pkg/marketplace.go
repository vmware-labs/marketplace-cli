// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

//go:generate counterfeiter . MarketplaceInterface
type MarketplaceInterface interface {
	EnableDebugging(printRequestPayloads bool, writer io.Writer)
	EnableStrictDecoding()
	DecodeJson(input io.Reader, output interface{}) error

	MakeURL(path string, params url.Values) *url.URL
	SendRequest(method string, requestURL *url.URL, headers map[string]string, content io.Reader) (*http.Response, error)

	GetAPIHost() string
	GetUIHost() string

	ListProducts(allOrgs bool, searchTerm string) ([]*models.Product, error)
	Get(requestURL *url.URL) (*http.Response, error)
	GetProduct(slug string) (*models.Product, error)
	GetProductWithVersion(slug, version string) (*models.Product, *models.Version, error)

	Post(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error)

	Put(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error)
	PutProduct(product *models.Product, versionUpdate bool) (*models.Product, error)

	GetUploadCredentials() (*CredentialsResponse, error)
	GetUploader(orgID string, credentials aws.Credentials) internal.Uploader

	Download(productId string, filename string, payload *DownloadRequestPayload) error
	DownloadChart(chartURL *url.URL) (*models.ChartVersion, error)
}

type Marketplace struct {
	Host           string
	APIHost        string
	UIHost         string
	StorageBucket  string
	StorageRegion  string
	Client         HTTPClient
	Output         io.Writer
	uploader       internal.Uploader
	strictDecoding bool
}

func (m *Marketplace) EnableStrictDecoding() {
	m.strictDecoding = true
}

func (m *Marketplace) EnableDebugging(printRequestPayloads bool, writer io.Writer) {
	m.Client = &DebuggingClient{
		client:               m.Client,
		logger:               log.New(writer, "", log.LstdFlags),
		printRequestPayloads: printRequestPayloads,
		printResposePayloads: false,
		requestID:            0,
	}
}

func (m *Marketplace) GetAPIHost() string {
	return m.APIHost
}

func (m *Marketplace) GetUIHost() string {
	return m.UIHost
}

func (m *Marketplace) DecodeJson(input io.Reader, output interface{}) error {
	d := json.NewDecoder(input)
	if m.strictDecoding {
		d.DisallowUnknownFields()
	}
	return d.Decode(output)
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
