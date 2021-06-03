package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/models"
)

var outputSupportsColor = false

//
//func init() {
//	fileInfo, _ := os.Stdout.Stat()
//	outputSupportsColor = (fileInfo.Mode() & os.ModeCharDevice) != 0
//}

func NewTableWriter(output io.Writer, headers ...string) *tablewriter.Table {
	table := tablewriter.NewWriter(output)
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
	if outputSupportsColor {
		var colors []tablewriter.Colors
		for range headers {
			colors = append(colors, []int{tablewriter.Bold})
		}
		table.SetHeaderColor(colors...)
	}
	return table
}

func RenderVersions(format string, product *models.Product, output io.Writer) error {
	if format == FormatTable {
		_, _ = fmt.Fprintln(output, "\nVersions:")
		table := NewTableWriter(output, "Number", "Status")
		for _, version := range product.Versions {
			table.Append([]string{version.Number, version.Status})
		}
		table.Render()

		for _, version := range product.Versions {
			err := RenderVersion(format, version.Number, product, output)
			if err != nil {
				return err
			}
		}
	} else if format == FormatJSON {
		return PrintJson(output, product.Versions)
	}
	return nil
}

func RenderVersion(format string, version string, product *models.Product, output io.Writer) error {
	if format == FormatTable {
		_, _ = fmt.Fprintf(output, "\nVersion %s:\n", version)
		dockerList := product.GetDockerImagesForVersion(version)
		if dockerList != nil {
			err := RenderContainerImages(format, dockerList, output)
			if err != nil {
				return err
			}
		}
		charts := product.GetChartsForVersion(version)
		if len(charts) > 0 {
			err := RenderCharts(format, charts, output)
			if err != nil {
				return err
			}
		}

	} else if format == FormatJSON {
		return PrintJson(output, product.Versions)
	}
	return nil
}

func RenderContainerImages(format string, images *models.DockerVersionList, output io.Writer) error {
	if format == FormatTable {
		table := NewTableWriter(output, "Image", "Tags")
		for _, docker := range images.DockerURLs {
			var tagList []string
			for _, tags := range docker.ImageTags {
				tagList = append(tagList, tags.Tag)
			}
			table.Append([]string{docker.Url, strings.Join(tagList, ", ")})
		}
		table.Render()
		_, _ = fmt.Fprintln(output, "Deployment instructions:")
		_, _ = fmt.Fprintln(output, images.DeploymentInstruction)
	} else if format == FormatJSON {
		return PrintJson(output, images)
	}
	return nil
}

func RenderContainerImage(format string, image *models.DockerURLDetails, output io.Writer) error {
	if format == FormatTable {
		table := NewTableWriter(output, "Tag", "Type")
		for _, tag := range image.ImageTags {
			table.Append([]string{tag.Tag, tag.Type})
		}
		table.Render()
	} else if format == FormatJSON {
		return PrintJson(output, image)
	}
	return nil
}

func RenderCharts(format string, charts []*models.ChartVersion, output io.Writer) error {
	if format == FormatTable {
		table := NewTableWriter(output, "Id", "Version", "URL", "Repository")
		for _, chart := range charts {
			table.Append([]string{
				chart.Id,
				chart.Version,
				chart.TarUrl,
				chart.Repo.Name + " " + chart.Repo.Url,
			})
		}
		table.Render()
	} else if format == FormatJSON {
		return PrintJson(output, charts)
	}
	return nil
}

func RenderProductList(format string, products []*models.Product, output io.Writer) error {
	if format == FormatTable {
		table := NewTableWriter(output, "Slug", "Name")
		for _, product := range products {
			table.Append([]string{product.Slug, product.DisplayName})
		}
		table.SetFooter([]string{"", "", fmt.Sprintf("Total count: %d", len(products))})
		table.SetAutoMergeCells(true)
		table.Render()
	} else if format == FormatJSON {
		return PrintJson(output, products)
	}
	return nil
}

func RenderProduct(format string, response *models.Product, output io.Writer) error {
	if format == FormatTable {
		_, _ = fmt.Fprintln(output, "Product Details:")
		table := NewTableWriter(output, "Slug", "Name")
		table.Append([]string{response.Slug, response.DisplayName})
		table.Render()
		return RenderVersions(format, response, output)
	} else if format == FormatJSON {
		return PrintJson(output, response)
	}
	return nil
}

func PrintJson(output io.Writer, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(output, string(data))
	return err
}
