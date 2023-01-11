// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"gopkg.in/yaml.v3"
)

type Encoder func(v interface{}) ([]byte, error)

type EncodedOutput struct {
	Marshall Encoder
	writer   io.Writer
}

func NewJSONOutput(writer io.Writer) *EncodedOutput {
	return &EncodedOutput{
		Marshall: json.Marshal,
		writer:   writer,
	}
}

func NewYAMLOutput(writer io.Writer) *EncodedOutput {
	return &EncodedOutput{
		Marshall: yaml.Marshal,
		writer:   writer,
	}
}

func (o *EncodedOutput) Print(object interface{}) error {
	data, err := o.Marshall(object)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(o.writer, string(data))
	return err
}

// PrintHeader is a no-op for encoded output. This output only prints the data
func (o *EncodedOutput) PrintHeader(message string) {}

func (o *EncodedOutput) RenderProduct(product *models.Product, _ *models.Version) error {
	return o.Print(product)
}

func (o *EncodedOutput) RenderProducts(products []*models.Product) error {
	return o.Print(products)
}

func (o *EncodedOutput) RenderVersions(product *models.Product) error {
	return o.Print(product.AllVersions)
}

func (o *EncodedOutput) RenderChart(chart *models.ChartVersion) error {
	return o.Print(chart)
}

func (o *EncodedOutput) RenderCharts(charts []*models.ChartVersion) error {
	return o.Print(charts)
}

func (o *EncodedOutput) RenderContainerImages(images []*models.DockerVersionList) error {
	return o.Print(images)
}

func (o *EncodedOutput) RenderFile(file *models.ProductDeploymentFile) error {
	return o.Print(file)
}

func (o *EncodedOutput) RenderFiles(files []*models.ProductDeploymentFile) error {
	return o.Print(files)
}

func (o *EncodedOutput) RenderAssets(assets []*pkg.Asset) error {
	return o.Print(assets)
}
