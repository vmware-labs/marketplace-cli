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

func (o *HumanOutput) PrintHeader(message string) {
	o.Println(message)
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
	footnote := ""
	downloads := ""
	if chart.IsUpdatedInMarketplaceRegistry {
		downloads = strconv.FormatInt(chart.DownloadCount, 10)
	} else {
		if chart.ProcessingError != "" {
			downloads = "Error*"
			footnote += fmt.Sprintf("* %s\n", chart.ProcessingError)
		} else {
			downloads = "Not yet available"
		}
	}

	table := o.NewTable("ID", "Version", "URL", "Repository", "Downloads")
	table.Append([]string{chart.Id, chart.Version, chart.TarUrl, chart.Repo.Name + " " + chart.Repo.Url, downloads})
	table.Render()

	if footnote != "" {
		o.Println()
		o.Println(footnote)
	}
	return nil
}

func (o *HumanOutput) RenderCharts(charts []*models.ChartVersion) error {
	footnotes := ""
	table := o.NewTable("ID", "Version", "URL", "Repository", "Downloads")
	for _, chart := range charts {
		downloads := ""
		if chart.IsUpdatedInMarketplaceRegistry {
			downloads = strconv.FormatInt(chart.DownloadCount, 10)
		} else {
			if chart.ProcessingError != "" {
				downloads = "Error*"
				footnotes += fmt.Sprintf("* %s\n", chart.ProcessingError)
			} else {
				downloads = "Not yet available"
			}
		}
		table.Append([]string{chart.Id, chart.Version, chart.TarUrl, chart.Repo.Name + " " + chart.Repo.Url, downloads})
	}
	table.Render()
	o.Printf("Total count: %d\n", len(charts))

	if footnotes != "" {
		o.Println()
		o.Println(footnotes)
	}
	return nil
}

func (o *HumanOutput) RenderContainerImage(image *models.DockerURLDetails) error {
	o.Println(image.Url)
	o.Println("Tags:")

	footnotes := ""
	table := o.NewTable("Tag", "Type", "Downloads")
	for _, tag := range image.ImageTags {
		downloads := ""
		if tag.IsUpdatedInMarketplaceRegistry {
			downloads = strconv.FormatInt(tag.DownloadCount, 10)
		} else {
			if tag.ProcessingError != "" {
				downloads = "Error*"
				footnotes += fmt.Sprintf("* %s\n", tag.ProcessingError)
			} else {
				downloads = "Not yet available"
			}
		}
		table.Append([]string{tag.Tag, tag.Type, downloads})
	}
	table.Render()
	o.Println()
	if footnotes != "" {
		o.Println(footnotes)
	}

	o.Println("Deployment instructions:")
	o.Println(image.DeploymentInstruction)

	return nil
}

func (o *HumanOutput) RenderContainerImages(images *models.DockerVersionList) error {
	var imageList []*models.DockerURLDetails
	if images != nil {
		imageList = images.DockerURLs
	}
	footnote := ""
	table := o.NewTable("Image", "Tags", "Downloads")
	for _, docker := range imageList {
		var tagList []string
		var downloads int64 = 0
		downloadable := true
		problem := false
		for _, tag := range docker.ImageTags {
			if downloadable && tag.IsUpdatedInMarketplaceRegistry {
				downloads += tag.DownloadCount
			} else {
				downloadable = false
				if tag.ProcessingError != "" {
					problem = true
				}
			}
			tagList = append(tagList, tag.Tag)
		}
		if downloadable {
			table.Append([]string{docker.Url, strings.Join(tagList, ", "), strconv.FormatInt(downloads, 10)})
		} else if problem {
			table.Append([]string{docker.Url, strings.Join(tagList, ", "), "Err*"})
			footnote = "* There is an error with this image."
		} else {
			table.Append([]string{docker.Url, strings.Join(tagList, ", "), "N/A"})
		}
	}
	table.Render()
	o.Println()
	o.Printf("Total count: %d\n", len(imageList))

	if footnote != "" {
		o.Println(footnote)
		o.Println()
	}

	if images != nil && images.DeploymentInstruction != "" {
		o.Println("Deployment instructions:")
		o.Println(images.DeploymentInstruction)
	}
	return nil
}

func (o *HumanOutput) RenderFile(file *models.ProductDeploymentFile) error {
	footnotes := ""
	table := o.NewTable("ID", "Name", "Status", "Size", "Type", "Files", "Downloads")
	downloads := ""
	if file.Status == "INACTIVE" {
		downloads = "Error*"
		footnotes += fmt.Sprintf("* %s\n", file.Comment)
	} else {
		downloads = strconv.FormatInt(file.DownloadCount, 10)
	}

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
		table.Append([]string{file.FileID, details.Name, file.Status, FormatSize(size), details.Type, strconv.Itoa(len(details.Files)), downloads})
	} else {
		table.Append([]string{file.FileID, file.Name, file.Status, "unknown*", "unknown*", "unknown*", downloads})
	}
	table.Render()

	if footnotes != "" {
		o.Println()
		o.Println(footnotes)
	}
	return nil
}

func (o *HumanOutput) RenderFiles(files []*models.ProductDeploymentFile) error {
	footnotes := ""
	table := o.NewTable("ID", "Name", "Status", "Size", "Type", "Files", "Downloads")
	for _, file := range files {
		downloads := ""
		if file.Status == "INACTIVE" {
			downloads = "Error*"
			footnotes += fmt.Sprintf("* %s\n", file.Comment)
		} else {
			downloads = strconv.FormatInt(file.DownloadCount, 10)
		}

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
			table.Append([]string{file.FileID, details.Name, file.Status, FormatSize(size), details.Type, strconv.Itoa(len(details.Files)), downloads})
		} else {
			table.Append([]string{file.FileID, file.Name, file.Status, "unknown*", "unknown*", "unknown*", downloads})
		}
	}
	table.Render()
	o.Printf("Total count: %d\n", len(files))

	if footnotes != "" {
		o.Println()
		o.Println(footnotes)
	}
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
