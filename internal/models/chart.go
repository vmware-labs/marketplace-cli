// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type Repo struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type ChartVersion struct {
	Id            string `json:"id,omitempty"`
	Version       string `json:"version,omitempty"`
	AppVersion    string `json:"appversion"`
	Details       string `json:"details,omitempty"`
	Readme        string `json:"readme,omitempty"`
	Repo          *Repo  `json:"repo,omitempty"`
	Values        string `json:"values,omitempty"`
	Digest        string `json:"digest,omitempty"`
	HashDigest    string `json:"hashdigest,omitempty"`
	HashAlgorithm string `json:"hashalgo,omitempty"`
	Status        string `json:"status,omitempty"`
	Size          int64  `json:"size,omitempty"`

	TarUrl                         string `json:"tarurl"` // to use during imgprocessor update & download from UI/API
	IsExternalUrl                  bool   `json:"isexternalurl"`
	HelmTarUrl                     string `json:"helmtarurl"` // to use during UI/API create & update product
	IsUpdatedInMarketplaceRegistry bool   `json:"isupdatedinmarketplaceregistry"`
	ProcessingError                string `json:"processingerror"`
	DownloadCount                  int64  `json:"downloadcount"`
	ValidationStatus               string `json:"validationstatus"`
	InstallOptions                 string `json:"installoptions"`
}

func (product *Product) GetChartsForVersion(version string) []*ChartVersion {
	var charts []*ChartVersion
	for _, chart := range product.ChartVersions {
		if chart.AppVersion == product.GetVersion(version).Number {
			charts = append(charts, chart)
		}
	}
	return charts
}

func (product *Product) GetChart(chartId string) *ChartVersion {
	for _, chart := range product.ChartVersions {
		if chart.Id == chartId {
			return chart
		}
	}
	return nil
}

func (product *Product) AddChart(chart *ChartVersion) {
	product.ChartVersions = append(product.ChartVersions, chart)
}
