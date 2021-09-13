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
	writer          io.Writer
	marketplaceHost string
}

func NewTableOutput(writer io.Writer, marketplaceHost string) *TableOutput {
	return &TableOutput{
		writer:          writer,
		marketplaceHost: marketplaceHost,
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
	_, _ = fmt.Fprintln(t.writer, product.DisplayName)
	_, _ = fmt.Fprintln(t.writer, product.Description.Summary)
	_, _ = fmt.Fprintf(t.writer, "https://%s/services/details/%s?slug=true\n", t.marketplaceHost, product.Slug)
	_, _ = fmt.Fprintln(t.writer, "\nProduct Details:")
	table := t.NewTable("Slug", "Type", "Latest Version")
	table.Append([]string{product.Slug, product.SolutionType, product.GetLatestVersion().Number})
	table.Render()
	_, _ = fmt.Fprintf(t.writer, "\nDescription:\n%s\n", product.Description.Description)
	return nil
}

func (t *TableOutput) RenderProducts(products []*models.Product) error {
	table := t.NewTable("Slug", "Name", "Type", "Latest Version")
	for _, product := range products {
		table.Append([]string{product.Slug, product.DisplayName, product.SolutionType, product.GetLatestVersion().Number})
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

	models.Sort(product.AllVersions)
	for _, version := range product.AllVersions {
		table.Append([]string{version.Number, version.Status})
	}
	table.Render()
	return nil
}

func (t *TableOutput) RenderChart(chart *models.ChartVersion) error {
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

func (t *TableOutput) RenderCharts(charts []*models.ChartVersion) error {
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

func (t *TableOutput) RenderContainerImage(image *models.DockerURLDetails) error {
	_, _ = fmt.Fprintf(t.writer, "%s\n", image.Url)
	_, _ = fmt.Fprintln(t.writer, "Tags:")

	footnotes := ""
	table := t.NewTable("Tag", "Type", "Downloads")
	for _, tag := range image.ImageTags {
		downloads := "N/A*"
		if tag.IsUpdatedInMarketplaceRegistry {
			downloads = strconv.FormatInt(tag.DownloadCount, 10)
		} else {
			footnotes += fmt.Sprintf("* %s\n", tag.ProcessingError)
		}
		table.Append([]string{tag.Tag, tag.Type, downloads})
	}
	table.Render()
	if footnotes != "" {
		_, _ = fmt.Fprintln(t.writer, footnotes)
	}

	_, _ = fmt.Fprintln(t.writer, "\nDeployment instructions:")
	_, _ = fmt.Fprintln(t.writer, image.DeploymentInstruction)

	return nil
}

func (t *TableOutput) RenderContainerImages(images *models.DockerVersionList) error {
	table := t.NewTable("Image", "Tags", "Downloads")
	for _, docker := range images.DockerURLs {
		var tagList []string
		var downloads int64 = 0
		downloadable := true
		for _, tag := range docker.ImageTags {
			if downloadable && tag.IsUpdatedInMarketplaceRegistry {
				downloads += tag.DownloadCount
			} else {
				downloadable = false
			}
			tagList = append(tagList, tag.Tag)
		}
		if downloadable {
			table.Append([]string{docker.Url, strings.Join(tagList, ", "), strconv.FormatInt(downloads, 10)})
		} else {
			table.Append([]string{docker.Url, strings.Join(tagList, ", "), "Err"})
		}
	}
	table.Render()
	if images.DeploymentInstruction != "" {
		_, _ = fmt.Fprintln(t.writer, "Deployment instructions:")
		_, _ = fmt.Fprintln(t.writer, images.DeploymentInstruction)
	}
	return nil
}

func (t *TableOutput) RenderOVA(file *models.ProductDeploymentFile) error {
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

func (t *TableOutput) RenderOVAs(ovas []*models.ProductDeploymentFile) error {
	table := t.NewTable("ID", "Name", "Status", "Size", "Type", "Files")
	for _, ova := range ovas {

		if ova.ItemJson != "" {
			details := &models.ProductItemDetails{}
			err := json.Unmarshal([]byte(ova.ItemJson), details)
			if err != nil {
				return fmt.Errorf("failed to parse the list of OVA files: %w", err)
			}

			var size int64 = 0
			for _, file := range details.Files {
				size += int64(file.Size)
			}
			table.Append([]string{ova.FileID, details.Name, ova.Status, FormatSize(size), details.Type, strconv.Itoa(len(details.Files))})
		} else {
			table.Append([]string{ova.FileID, ova.Name, ova.Status, "unknown", "unknown", "unknown"})
		}

	}
	table.Render()
	return nil
}
