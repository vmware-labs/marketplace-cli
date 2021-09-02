// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

type TableOutput struct {
	writer io.Writer
}

func NewTableOutput(writer io.Writer) *TableOutput {
	return &TableOutput{
		writer: writer,
	}
}

func (t *TableOutput) NewTable(headers ...string) *tablewriter.Table {
	table := tablewriter.NewWriter(t.writer)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetColWidth(300)
	table.SetTablePadding("\t\t")
	table.SetHeader(headers)
	//if outputSupportsColor {
	//	var colors []tablewriter.Colors
	//	for range headers {
	//		colors = append(colors, []int{tablewriter.Bold})
	//	}
	//	table.SetHeaderColor(colors...)
	//}
	return table
}

func (t *TableOutput) RenderProduct(product *models.Product) error {
	_, _ = fmt.Fprintln(t.writer, "Product Details:")
	table := t.NewTable("Slug", "Name", "Type")
	table.Append([]string{product.Slug, product.DisplayName, product.SolutionType})
	table.Render()
	return t.RenderVersions(product)
}

func (t *TableOutput) RenderProducts(products []*models.Product) error {
	table := t.NewTable("Slug", "Name", "Type", "Latest Version")
	for _, product := range products {
		latestVersion := "N/A"
		if len(product.AllVersions) > 0 {
			latestVersion = product.AllVersions[len(product.AllVersions)-1].Number
		}
		table.Append([]string{product.Slug, product.DisplayName, product.SolutionType, latestVersion})
	}
	table.SetFooter([]string{"", "", "", fmt.Sprintf("Total count: %d", len(products))})
	table.Render()
	return nil
}

func (t *TableOutput) RenderVersion(product *models.Product, version string) error {
	_, _ = fmt.Fprintf(t.writer, "Version %s\n", version)
	return nil
}

func (t *TableOutput) RenderVersions(product *models.Product) error {
	_, _ = fmt.Fprintln(t.writer, "Versions:")
	table := t.NewTable("Number", "Status")
	for _, version := range product.AllVersions {
		table.Append([]string{version.Number, version.Status})
	}
	table.Render()
	return nil
}

func (t *TableOutput) RenderChart(product *models.Product, version string, chart *models.ChartVersion) error {
	table := t.NewTable("ID", "Version", "URL", "Repository")
	table.Append([]string{
		chart.Id,
		chart.Version,
		chart.TarUrl,
		chart.Repo.Name + " " + chart.Repo.Url,
	})
	table.Render()
	return nil
}

func (t *TableOutput) RenderCharts(product *models.Product, version string) error {
	charts := product.GetChartsForVersion(version)
	if len(charts) == 0 {
		_, _ = fmt.Fprintf(t.writer, "%s %s does not have any charts\n", product.Slug, version)
		return nil
	}

	table := t.NewTable("ID", "Version", "URL", "Repository")
	for _, chart := range charts {
		table.Append([]string{
			chart.Id,
			chart.Version,
			chart.TarUrl,
			chart.Repo.Name + " " + chart.Repo.Url,
		})
	}
	table.Render()
	return nil
}

func (t *TableOutput) RenderContainerImage(product *models.Product, version string, image *models.DockerURLDetails) error {
	table := t.NewTable("Tag", "Type")
	for _, tag := range image.ImageTags {
		table.Append([]string{tag.Tag, tag.Type})
	}
	table.Render()
	return nil
}

func (t *TableOutput) RenderContainerImages(product *models.Product, version string) error {
	images := product.GetContainerImagesForVersion(version)
	if images == nil || len(images.DockerURLs) == 0 {
		_, _ = fmt.Fprintf(t.writer, "%s %s does not have any container images\n", product.Slug, version)
		return nil
	}

	table := t.NewTable("Image", "Tags")
	for _, docker := range images.DockerURLs {
		var tagList []string
		for _, tags := range docker.ImageTags {
			tagList = append(tagList, tags.Tag)
		}
		table.Append([]string{docker.Url, strings.Join(tagList, ", ")})
	}
	table.Render()
	_, _ = fmt.Fprintln(t.writer, "Deployment instructions:")
	_, _ = fmt.Fprintln(t.writer, images.DeploymentInstruction)
	return nil
}

func (t *TableOutput) RenderOVA(product *models.Product, version string, file *models.ProductDeploymentFile) error {
	table := t.NewTable("ID", "Name", "Status", "Size", "Type", "Files")
	if file.ItemJson != "" {
		details := &models.ProductItemDetails{}
		err := json.Unmarshal([]byte(file.ItemJson), details)
		if err != nil {
			return fmt.Errorf("failed to parse the list of OVA files: %w", err)
		}

		var size int64 = 0
		for _, file := range details.Files {
			size += int64(file.Size)
		}
		table.Append([]string{file.FileID, details.Name, file.Status, FormatSize(size), details.Type, strconv.Itoa(len(details.Files))})
	} else {
		table.Append([]string{file.FileID, "unknown", file.Status, "unknown", "unknown", "unknown"})
	}

	table.Render()
	return nil
}

func (t *TableOutput) RenderOVAs(product *models.Product, version string) error {
	ovas := product.GetFilesForVersion(version)
	if len(ovas) == 0 {
		_, _ = fmt.Fprintf(t.writer, "product \"%s\" %s does not have any OVAs\n", product.Slug, version)
		return nil
	}

	table := t.NewTable("ID", "Name", "Status", "Size", "Type", "Files")
	for _, ova := range ovas {

		if ova.ItemJson != "" {
			details := &models.ProductItemDetails{}
			err := json.Unmarshal([]byte(ova.ItemJson), details)
			if err != nil {
				return fmt.Errorf("failed to parse the list of OVA files: %w", err)
			}

			size := 0
			for _, file := range details.Files {
				size += file.Size
			}
			table.Append([]string{ova.FileID, details.Name, ova.Status, strconv.Itoa(size), details.Type, strconv.Itoa(len(details.Files))})
		} else {
			table.Append([]string{ova.FileID, "unknown", ova.Status, "unknown", "unknown", "unknown"})
		}

	}
	table.Render()
	return nil
}
