// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"helm.sh/helm/v3/pkg/chart/loader"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func LoadChart(chartPath string) (*models.ChartVersion, error) {
	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read chart: %w", err)
	}

	return &models.ChartVersion{
		Version: chart.Metadata.Version,
		Repo: &models.Repo{
			Name: chart.Name(),
		},
		IsExternalUrl: false,
	}, nil
}

func (m *Marketplace) DownloadChart(chartURL *url.URL) (*models.ChartVersion, error) {
	req, err := http.NewRequest("GET", chartURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to download chart: %w", err)
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download chart: %w", err)
	}

	chartFile, err := os.CreateTemp("", "chart-*.tgz")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary local chart: %w", err)
	}

	_, err = io.Copy(chartFile, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save local chart: %w", err)
	}

	_ = resp.Body.Close()
	_ = chartFile.Close()

	chart, err := LoadChart(chartFile.Name())
	if err != nil {
		return nil, err
	}
	chart.IsExternalUrl = true
	chart.HelmTarUrl = chartURL.String()
	return chart, nil
}
