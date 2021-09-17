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

type HumanOutput struct {
	writer          io.Writer
	marketplaceHost string
}

func (o *HumanOutput) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(o.writer, format, a...)
}
func (o *HumanOutput) Println(a ...interface{}) {
	_, _ = fmt.Fprintln(o.writer, a...)
}

func NewHumanOutput(writer io.Writer, marketplaceHost string) *HumanOutput {
	return &HumanOutput{
		writer:          writer,
		marketplaceHost: marketplaceHost,
	}
}

func (o *HumanOutput) NewTable(headers ...string) *tablewriter.Table {
	table := tablewriter.NewWriter(o.writer)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(true)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetColWidth(100)
	table.SetHeader(headers)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetRowSeparator("")
	table.SetTablePadding("\t\t")
	//if outputSupportsColor {
	//	var colors []tablewriter.Colors
	//	for range headers {
	//		colors = append(colors, []int{tablewriter.Bold})
	//	}
	//	table.SetHeaderColor(colors...)
	//}
	return table
}

func (o *HumanOutput) RenderProduct(product *models.Product) error {
	o.Printf("Name:      %s\n", product.DisplayName)
	o.Printf("Publisher: %s\n", product.PublisherDetails.OrgDisplayName)
	o.Println()
	o.Println(product.Description.Summary)
	o.Printf("https://%s/services/details/%s?slug=true\n", o.marketplaceHost, product.Slug)
	o.Println()
	o.Println("Product Details:")
	table := o.NewTable("Product ID", "Slug", "Type", "Latest Version")
	table.Append([]string{product.ProductId, product.Slug, product.SolutionType, product.GetLatestVersion().Number})
	table.Render()
	o.Println()
	o.Println("Description:")
	o.Println(product.Description.Description)
	return nil
}

func (o *HumanOutput) RenderProducts(products []*models.Product) error {
	table := o.NewTable("Slug", "Name", "Publisher", "Type", "Latest Version")
	for _, product := range products {
		table.Append([]string{product.Slug, product.DisplayName, product.PublisherDetails.OrgDisplayName, product.SolutionType, product.GetLatestVersion().Number})
	}
	table.Render()
	o.Printf("Total count: %d\n", len(products))
	return nil
}

func (o *HumanOutput) RenderVersion(product *models.Product, version string) error {
	o.Printf("Version %s\n", version)
	return nil
}

func (o *HumanOutput) RenderVersions(product *models.Product) error {
	o.Println("Versions:")
	table := o.NewTable("Number", "Status")

	models.Sort(product.AllVersions)
	for _, version := range product.AllVersions {
		table.Append([]string{version.Number, version.Status})
	}
	table.Render()
	return nil
}

func (o *HumanOutput) RenderChart(chart *models.ChartVersion) error {
	table := o.NewTable("ID", "Version", "URL", "Repository")
	table.Append([]string{
		chart.Id,
		chart.Version,
		chart.TarUrl,
		chart.Repo.Name + " " + chart.Repo.Url,
	})
	table.Render()
	return nil
}

func (o *HumanOutput) RenderCharts(charts []*models.ChartVersion) error {
	table := o.NewTable("ID", "Version", "URL", "Repository")
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

func (o *HumanOutput) RenderContainerImage(image *models.DockerURLDetails) error {
	o.Printf("%s\n", image.Url)
	o.Println("Tags:")

	footnotes := ""
	table := o.NewTable("Tag", "Type", "Downloads")
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
		o.Println(footnotes)
	}

	o.Println("\nDeployment instructions:")
	o.Println(image.DeploymentInstruction)

	return nil
}

func (o *HumanOutput) RenderContainerImages(images *models.DockerVersionList) error {
	table := o.NewTable("Image", "Tags", "Downloads")
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
		o.Println("Deployment instructions:")
		o.Println(images.DeploymentInstruction)
	}
	return nil
}

func (o *HumanOutput) RenderOVA(file *models.ProductDeploymentFile) error {
	table := o.NewTable("ID", "Name", "Status", "Size", "Type", "Files")
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

func (o *HumanOutput) RenderOVAs(ovas []*models.ProductDeploymentFile) error {
	table := o.NewTable("ID", "Name", "Status", "Size", "Type", "Files")
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

func (o *HumanOutput) RenderSubscriptions(subscriptions []*models.Subscription) error {
	table := o.NewTable("ID", "Product ID", "Product Name", "Status")
	for _, subscription := range subscriptions {
		table.Append([]string{strconv.Itoa(subscription.ID), subscription.ProductID, subscription.ProductName, subscription.StatusText})
	}
	table.Render()
	o.Printf("Total count: %d\n", len(subscriptions))
	return nil
}
