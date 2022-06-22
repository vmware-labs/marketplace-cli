// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//go:generate counterfeiter . HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient() HTTPClient {
	return &http.Client{}
}

type DebuggingClient struct {
	client               HTTPClient
	logger               *log.Logger
	printRequestPayloads bool
	printResposePayloads bool
	requestID            int
}

func (c *DebuggingClient) printRequest(req *http.Request) int {
	requestID := c.requestID
	c.requestID++
	c.logger.Printf("Request #%d: %s %s\n", requestID, req.Method, req.URL.String())
	if c.printRequestPayloads && req.ContentLength > 0 {
		req.Body = c.printPayload(fmt.Sprintf("request #%d body", requestID), req.Body)
	}

	return requestID
}

func (c *DebuggingClient) printResponse(requestID int, resp *http.Response) {
	if resp != nil {
		c.logger.Printf("Request #%d Response: %s", requestID, resp.Status)
		if c.printResposePayloads {
			resp.Body = c.printPayload(fmt.Sprintf("request #%d response body", requestID), resp.Body)
		}
	}
}

func (c *DebuggingClient) Do(req *http.Request) (*http.Response, error) {
	requestID := c.printRequest(req)
	resp, err := c.client.Do(req)
	c.printResponse(requestID, resp)
	return resp, err
}

func (c *DebuggingClient) printPayload(name string, payload io.ReadCloser) io.ReadCloser {
	c.logger.Printf("--- Start of %s payload ---", name)
	content, _ := ioutil.ReadAll(payload)
	c.logger.Println(string(content))
	c.logger.Printf("--- End of %s payload ---", name)

	return io.NopCloser(bytes.NewReader(content))
}

type QueryStringParameter interface {
	QueryString() string
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
