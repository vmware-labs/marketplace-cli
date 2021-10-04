// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"gopkg.in/yaml.v2"
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

func (o *EncodedOutput) RenderProduct(product *models.Product) error {
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

func (o *EncodedOutput) RenderContainerImage(image *models.DockerURLDetails) error {
	return o.Print(image)
}

func (o *EncodedOutput) RenderContainerImages(images *models.DockerVersionList) error {
	return o.Print(images)
}

func (o *EncodedOutput) RenderOVA(file *models.ProductDeploymentFile) error {
	return o.Print(file)
}

func (o *EncodedOutput) RenderOVAs(files []*models.ProductDeploymentFile) error {
	return o.Print(files)
}

func (o *EncodedOutput) RenderSubscriptions(subscriptions []*models.Subscription) error {
	return o.Print(subscriptions)
}
