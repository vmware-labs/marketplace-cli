// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

type JSONOutput struct {
	writer io.Writer
}

func NewJSONOutput(writer io.Writer) *JSONOutput {
	return &JSONOutput{
		writer: writer,
	}
}

func (j *JSONOutput) PrintJSON(object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(j.writer, string(data))
	return err
}

func (j *JSONOutput) RenderProduct(product *models.Product) error {
	return j.PrintJSON(product)
}

func (j *JSONOutput) RenderProducts(products []*models.Product) error {
	return j.PrintJSON(products)
}

func (j *JSONOutput) RenderVersion(product *models.Product, version string) error {
	return j.PrintJSON(product.GetVersion(version))
}

func (j *JSONOutput) RenderVersions(product *models.Product) error {
	return j.PrintJSON(product.AllVersions)
}

func (j *JSONOutput) RenderChart(chart *models.ChartVersion) error {
	return j.PrintJSON(chart)
}

func (j *JSONOutput) RenderCharts(charts []*models.ChartVersion) error {
	return j.PrintJSON(charts)
}

func (j *JSONOutput) RenderContainerImage(image *models.DockerURLDetails) error {
	return j.PrintJSON(image)
}

func (j *JSONOutput) RenderContainerImages(images *models.DockerVersionList) error {
	return j.PrintJSON(images)
}

func (j *JSONOutput) RenderOVA(file *models.ProductDeploymentFile) error {
	return j.PrintJSON(file)
}

func (j *JSONOutput) RenderOVAs(files []*models.ProductDeploymentFile) error {
	return j.PrintJSON(files)
}
