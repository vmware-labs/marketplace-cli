// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
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
	IsAddonFile         string `json:"isAddonFile,omitempty"`
	AddonFileId         string `json:"addonFileId,omitempty"`
	MetaFileID          string `json:"metafileid,omitempty"`
	MetaFileObjectID    string `json:"metafileobjectid,omitempty"`
}

type DownloadResponse struct {
	Response struct {
		PreSignedURL string `json:"presignedurl"`
		Message      string `json:"message"`
		StatusCode   int    `json:"statuscode"`
	} `json:"response"`
}

func (m *Marketplace) Download(productId string, filename string, payload *DownloadRequestPayload) error {
	requestURL := m.MakeURL(fmt.Sprintf("/api/v1/products/%s/download", productId), nil)
	resp, err := m.PostJSON(requestURL, payload)
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

	progressBar := m.makeDownloadProgressBar(resp.ContentLength, filename)
	_, err = io.Copy(io.MultiWriter(file, progressBar), resp.Body)
	return err
}

func (m *Marketplace) makeDownloadProgressBar(length int64, filename string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		length,
		progressbar.OptionSetDescription(fmt.Sprintf("Downloading to %s", filename)),
		progressbar.OptionSetWriter(m.Output),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprintln(m.Output, "")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	_ = bar.RenderBlank()
	return bar
}
