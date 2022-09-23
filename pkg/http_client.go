// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/viper"
)

//go:generate counterfeiter . HTTPClient
type HTTPClient interface {
	Get(requestURL *url.URL) (*http.Response, error)
	Post(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error)
	PostForm(requestURL *url.URL, content url.Values) (resp *http.Response, err error)
	PostJSON(requestURL *url.URL, content interface{}) (*http.Response, error)
	Put(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error)
	SendRequest(method string, requestURL *url.URL, headers map[string]string, content io.Reader) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

//go:generate counterfeiter . PerformRequestFunc
type PerformRequestFunc func(req *http.Request) (*http.Response, error)

type DebuggingClient struct {
	Logger               *log.Logger
	PrintRequests        bool
	PrintRequestPayloads bool
	PrintResposePayloads bool
	requestID            int
	PerformRequest       PerformRequestFunc
}

func NewClient(output io.Writer, printRequests, printRequestPayloads, printResponsePayloads bool) *DebuggingClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	if viper.GetBool("skip_ssl_validation") {
		transport := cleanhttp.DefaultPooledTransport()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		retryClient.HTTPClient.Transport = transport
	}

	return &DebuggingClient{
		Logger:               log.New(output, "", log.LstdFlags),
		PrintRequests:        printRequests,
		PrintRequestPayloads: printRequestPayloads,
		PrintResposePayloads: printResponsePayloads,
		requestID:            0,
		PerformRequest:       retryClient.StandardClient().Do,
	}
}

func (c *DebuggingClient) printRequest(req *http.Request) int {
	requestID := c.requestID
	c.requestID++
	if c.PrintRequests {
		c.Logger.Printf("Request #%d: %s %s\n", requestID, req.Method, req.URL.String())
	}
	if c.PrintRequestPayloads && req.ContentLength > 0 {
		req.Body = c.printPayload(fmt.Sprintf("request #%d body", requestID), req.Body)
	}

	return requestID
}

func (c *DebuggingClient) printResponse(requestID int, resp *http.Response) {
	if c.PrintRequests && resp != nil {
		c.Logger.Printf("Request #%d Response: %s", requestID, resp.Status)
		if c.PrintResposePayloads {
			resp.Body = c.printPayload(fmt.Sprintf("request #%d response body", requestID), resp.Body)
		}
	}
}

func (c *DebuggingClient) printPayload(name string, payload io.ReadCloser) io.ReadCloser {
	c.Logger.Printf("--- Start of %s payload ---", name)
	content, _ := io.ReadAll(payload)
	c.Logger.Println(string(content))
	c.Logger.Printf("--- End of %s payload ---", name)

	return io.NopCloser(bytes.NewReader(content))
}

func (c *DebuggingClient) Get(requestURL *url.URL) (*http.Response, error) {
	return c.SendRequest("GET", requestURL, map[string]string{}, nil)
}

func (c *DebuggingClient) Post(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": contentType,
	}
	return c.SendRequest("POST", requestURL, headers, content)
}

func (c *DebuggingClient) PostForm(requestURL *url.URL, content url.Values) (resp *http.Response, err error) {
	return c.Post(requestURL, strings.NewReader(content.Encode()), "application/x-www-form-urlencoded")
}

func (c *DebuggingClient) PostJSON(requestURL *url.URL, content interface{}) (*http.Response, error) {
	encoded, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request payload: %w", err)
	}

	return c.Post(requestURL, bytes.NewReader(encoded), "application/json")
}

func (c *DebuggingClient) Put(requestURL *url.URL, content io.Reader, contentType string) (*http.Response, error) {
	headers := map[string]string{}
	if contentType != "" {
		headers["Content-Type"] = contentType
	}
	return c.SendRequest("PUT", requestURL, headers, content)
}

func (c *DebuggingClient) SendRequest(method string, requestURL *url.URL, headers map[string]string, content io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, requestURL.String(), content)
	if err != nil {
		return nil, fmt.Errorf("failed to build %s request: %w", requestURL.String(), err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Accept", "application/json")
	CSPAPIToken := viper.GetString("csp.refresh-token")
	if CSPAPIToken != "" {
		req.Header.Add("csp-auth-token", viper.GetString("csp.refresh-token"))
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *DebuggingClient) Do(req *http.Request) (*http.Response, error) {
	requestID := c.printRequest(req)
	resp, err := c.PerformRequest(req)
	c.printResponse(requestID, resp)
	return resp, err
}

type QueryStringParameter interface {
	QueryString() string
}

func MakeURL(host, path string, values url.Values) *url.URL {
	queryString := ""
	if values != nil {
		queryString = values.Encode()
	}
	return &url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     path,
		RawQuery: queryString,
	}
}

func ApplyParameters(url *url.URL, parameters ...QueryStringParameter) {
	var queryStrings []string
	if url.RawQuery != "" {
		queryStrings = append(queryStrings, url.RawQuery)
	}
	for _, param := range parameters {
		queryStrings = append(queryStrings, param.QueryString())
	}
	url.RawQuery = strings.Join(queryStrings, "&")
}
