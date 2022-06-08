// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

type DownloadRequestPayload struct {
	ProductId           string `json:"productid,omitempty"`
	AppVersion          string `json:"appVersion,omitempty"`
	EulaAccepted        bool   `json:"eulaAccepted"`
	DockerlinkVersionID string `json:"dockerlinkVersionId,omitempty"`
	DockerUrlId         string `json:"dockerUrlId,omitempty"`
	ImageTagId          string `json:"imageTagId,omitempty"`
	DeploymentFileId    string `json:"deploymentFileId,omitempty"`
	ChartVersion        string `json:"chartVersion,omitempty"`
	IsAddonFile         bool   `json:"isAddonFile,omitempty"`
	AddonFileId         string `json:"addonFileId,omitempty"`
	MetaFileID          string `json:"metafileid,omitempty"`
	MetaFileObjectID    string `json:"metafileobjectid,omitempty"`
}

type DownloadResponseBody struct {
	PreSignedURL string `json:"presignedurl"`
	Message      string `json:"message"`
	StatusCode   int    `json:"statuscode"`
}
type DownloadResponse struct {
	Response *DownloadResponseBody `json:"response"`
}

func (m *Marketplace) Download(filename string, payload *DownloadRequestPayload) error {
	requestURL := m.MakeURL(fmt.Sprintf("/api/v1/products/%s/download", payload.ProductId), nil)
	resp, err := m.PostJSON(requestURL, payload)
	if err != nil {
		return fmt.Errorf("failed to get download link: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return fmt.Errorf("failed to fetch download link: %s\n%s", resp.Status, string(body))
		}
		return fmt.Errorf("failed to fetch download link: %s", resp.Status)
	}

	downloadResponse := &DownloadResponse{}
	err = m.DecodeJson(resp.Body, downloadResponse)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return m.downloadFile(filename, downloadResponse.Response.PreSignedURL)
}

func (m *Marketplace) downloadFile(filename string, fileDownloadURL string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file for download: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequest("GET", fileDownloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download file request: %w", err)
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	progressBar := internal.MakeProgressBar(fmt.Sprintf("Downloading %s", filename), resp.ContentLength, m.Output)
	_, err = io.Copy(progressBar.WrapWriter(file), resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download file to disk: %w", err)
	}
	return nil
}
