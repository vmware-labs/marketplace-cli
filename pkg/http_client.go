// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

func (c *DebuggingClient) Do(req *http.Request) (*http.Response, error) {
	currentRequest := c.requestID
	c.requestID++
	c.logger.Printf("Request #%d: %s %s\n", currentRequest, req.Method, req.URL.String())
	if c.printRequestPayloads && req.ContentLength > 0 {
		req.Body = c.printPayload(fmt.Sprintf("request #%d body", currentRequest), req.Body)
	}

	resp, err := c.client.Do(req)

	c.logger.Printf("Request #%d Response: %s", currentRequest, resp.Status)
	if c.printResposePayloads {
		resp.Body = c.printPayload(fmt.Sprintf("request #%d response body", currentRequest), resp.Body)
	}

	return resp, err
}

func (c *DebuggingClient) printPayload(name string, payload io.ReadCloser) io.ReadCloser {
	c.logger.Printf("--- Start of %s payload ---", name)
	content, _ := ioutil.ReadAll(payload)
	c.logger.Println(string(content))
	c.logger.Printf("--- End of %s payload ---", name)

	return io.NopCloser(bytes.NewReader(content))
}

func EnableDebugging(printRequestPayloads bool, existingClient HTTPClient, writer io.Writer) HTTPClient {
	return &DebuggingClient{
		client:               existingClient,
		logger:               log.New(writer, "", log.LstdFlags),
		printRequestPayloads: printRequestPayloads,
		printResposePayloads: false,
		requestID:            0,
	}
}
