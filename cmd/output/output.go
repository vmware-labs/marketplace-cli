// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output

import (
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

const (
	FormatJSON  = "json"
	FormatTable = "table"
)

var SupportedOutputs = []string{FormatJSON, FormatTable}

//go:generate counterfeiter . Format
type Format interface {
	RenderProduct(product *models.Product) error
	RenderProducts(products []*models.Product) error
	RenderVersion(product *models.Product, version string) error
	RenderVersions(product *models.Product) error
	RenderChart(chart *models.ChartVersion) error
	RenderCharts(charts []*models.ChartVersion) error
	RenderContainerImage(image *models.DockerURLDetails) error
	RenderContainerImages(*models.DockerVersionList) error
	RenderOVA(file *models.ProductDeploymentFile) error
	RenderOVAs(files []*models.ProductDeploymentFile) error
}
