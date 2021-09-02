// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mitchellh/ioprogress"
)

type DownloadRequestPayload struct {
	DockerlinkVersionID string `json:"dockerlinkVersionId,omitempty"`
	DockerUrlId         string `json:"dockerUrlId,omitempty"`
	ImageTagId          string `json:"imageTagId,omitempty"`
	DeploymentFileId    string `json:"deploymentFileId,omitempty"`
	AppVersion          string `json:"appVersion,omitempty"`
	ChartVersion        string `json:"chartVersion,omitempty"`
	IsAddonFile         string `json:"isAddonFile,omitempty"`
	AddonFileId         string `json:"addonFileId,omitempty"`
	EulaAccepted        bool   `json:"eulaAccepted"`
}

type DownloadResponse struct {
	Response struct {
		PreSignedURL string `json:"presignedurl"`
		Message      string `json:"message"`
		StatusCode   int    `json:"statuscode"`
	} `json:"response"`
}

func (m *Marketplace) Download(productId string, filename string, payload *DownloadRequestPayload, output io.Writer) error {
	requestURL := m.MakeURL(fmt.Sprintf("/api/v1/products/%s/download", productId), nil)
	encoded, _ := json.Marshal(payload)

	resp, err := m.Post(requestURL, bytes.NewReader(encoded), "application/json")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return fmt.Errorf("failed to fetch download link: %s\n%s", resp.Status, string(body))
		}
		return fmt.Errorf("failed to fetch download link: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	downloadResponse := &DownloadResponse{}
	err = json.Unmarshal(body, downloadResponse)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return m.downloadFile(filename, downloadResponse.Response.PreSignedURL, output)
}

func (m *Marketplace) downloadFile(filename string, fileDownloadURL string, output io.Writer) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("GET", fileDownloadURL, nil)
	if err != nil {
		return err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	progressBody := &ioprogress.Reader{
		Reader:   resp.Body,
		Size:     resp.ContentLength,
		DrawFunc: ioprogress.DrawTerminalf(output, ioprogress.DrawTextFormatBytes),
	}

	_, err = io.Copy(file, progressBody)
	return err
}
