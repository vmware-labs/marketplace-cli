// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output

import (
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

const (
	FormatHuman = "human"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
)

var SupportedOutputs = []string{FormatHuman, FormatJSON, FormatYAML}

//go:generate counterfeiter . Format
type Format interface {
	PrintHeader(message string)

	RenderProduct(product *models.Product) error
	RenderProducts(products []*models.Product) error
	RenderVersions(product *models.Product) error
	RenderChart(chart *models.ChartVersion) error
	RenderCharts(charts []*models.ChartVersion) error
	RenderContainerImage(image *models.DockerURLDetails) error
	RenderContainerImages(images *models.DockerVersionList) error
	RenderFile(file *models.ProductDeploymentFile) error
	RenderFiles(files []*models.ProductDeploymentFile) error

	RenderSubscriptions(subscriptions []*models.Subscription) error
}
