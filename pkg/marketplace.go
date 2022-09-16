// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

//go:generate counterfeiter . MarketplaceInterface
type MarketplaceInterface interface {
	EnableStrictDecoding()
	DecodeJson(input io.Reader, output interface{}) error

	GetHost() string
	GetAPIHost() string
	GetUIHost() string

	ListProducts(filter *ListProductFilter) ([]*models.Product, error)
	GetProduct(slug string) (*models.Product, error)
	GetProductWithVersion(slug, version string) (*models.Product, *models.Version, error)
	PutProduct(product *models.Product, versionUpdate bool) (*models.Product, error)

	GetUploader(orgID string) (internal.Uploader, error)
	SetUploader(uploader internal.Uploader)

	Download(filename string, payload *DownloadRequestPayload) error

	DownloadChart(chartURL *url.URL) (*models.ChartVersion, error)
	AttachLocalChart(chartPath, instructions string, product *models.Product, version *models.Version) (*models.Product, error)
	AttachPublicChart(chartPath *url.URL, instructions string, product *models.Product, version *models.Version) (*models.Product, error)

	AttachLocalContainerImage(imageFile, image, tag, tagType, instructions string, product *models.Product, version *models.Version) (*models.Product, error)
	AttachPublicContainerImage(image, tag, tagType, instructions string, product *models.Product, version *models.Version) (*models.Product, error)

	AttachMetaFile(metafile, metafileType, metafileVersion string, product *models.Product, version *models.Version) (*models.Product, error)

	AttachOtherFile(file string, product *models.Product, version *models.Version) (*models.Product, error)

	UploadVM(vmFile string, product *models.Product, version *models.Version) (*models.Product, error)
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

func (m *Marketplace) GetHost() string {
	return m.Host
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
