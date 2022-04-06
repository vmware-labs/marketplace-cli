// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	ChartID             string
	ChartReadme         string
	ChartProductSlug    string
	ChartProductVersion string
	ChartURL            string
)

func init() {
	rootCmd.AddCommand(ChartCmd)
	ChartCmd.AddCommand(ListChartsCmd)
	ChartCmd.AddCommand(GetChartCmd)
	ChartCmd.AddCommand(AttachChartCmd)

	ListChartsCmd.Flags().StringVarP(&ChartProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListChartsCmd.MarkFlagRequired("product")
	ListChartsCmd.Flags().StringVarP(&ChartProductVersion, "product-version", "v", "", "Product version (default to latest version)")

	GetChartCmd.Flags().StringVarP(&ChartProductSlug, "product", "p", "", "Product slug (required)")
	_ = GetChartCmd.MarkFlagRequired("product")
	GetChartCmd.Flags().StringVarP(&ChartProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	GetChartCmd.Flags().StringVar(&ChartID, "chart-id", "", "chart ID")

	AttachChartCmd.Flags().StringVarP(&ChartProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachChartCmd.MarkFlagRequired("product")
	AttachChartCmd.Flags().StringVarP(&ChartProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	//AttachChartCmd.Flags().StringVarP(&ChartVersion, "chart-version", "e", "", "chart version (required)")
	//_ = AttachChartCmd.MarkFlagRequired("chart-version")
	AttachChartCmd.Flags().StringVarP(&ChartURL, "chart", "c", "", "path to to chart, either local tgz or public URL (required)")
	_ = AttachChartCmd.MarkFlagRequired("chart")
	//AttachChartCmd.Flags().StringVar(&ChartRepositoryURL, "repository-url", "", "chart public repository url")
	//_ = AttachChartCmd.MarkFlagRequired("repository-url")
	//AttachChartCmd.Flags().StringVar(&ChartRepositoryName, "repository-name", "", "chart public repository name")
	//_ = AttachChartCmd.MarkFlagRequired("repository-name")
	AttachChartCmd.Flags().StringVar(&ChartReadme, "readme", "", "readme information")
	_ = AttachChartCmd.MarkFlagRequired("readme")
}

var ChartCmd = &cobra.Command{
	Use:       "chart",
	Aliases:   []string{"charts"},
	Short:     "List and manage Helm charts attached to a product",
	Long:      "List and manage Helm charts attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListChartsCmd.Use, GetChartCmd.Use, AttachChartCmd.Use},
}

var ListChartsCmd = &cobra.Command{
	Use:     "list",
	Short:   "List product charts",
	Long:    "Prints the list of Helm charts attached to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ChartProductSlug, ChartProductVersion)
		if err != nil {
			return err
		}

		charts := product.GetChartsForVersion(version.Number)

		Output.PrintHeader(fmt.Sprintf("Charts for %s %s:", product.DisplayName, version.Number))
		return Output.RenderCharts(charts)
	},
}

var GetChartCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get details for a chart",
	Long:    "Prints detailed information about a Helm chart attached to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ChartProductSlug, ChartProductVersion)
		if err != nil {
			return err
		}

		var chart *models.ChartVersion
		if ChartID != "" {
			chart = product.GetChart(ChartID)
			if chart == nil {
				return fmt.Errorf("%s %s does not have the chart \"%s\"", ChartProductSlug, version.Number, ChartID)
			}
		} else {
			charts := product.GetChartsForVersion(version.Number)
			if len(charts) == 0 {
				return fmt.Errorf("%s %s does not have any charts", ChartProductSlug, version.Number)
			} else if len(charts) == 1 {
				chart = charts[0]
			} else {
				return fmt.Errorf("multiple charts found for %s for version %s, please use the --chard-id parameter", ChartProductSlug, version.Number)
			}
		}

		return Output.RenderChart(chart)
	},
}

var AttachChartCmd = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a chart",
	Long:    "Attaches a Helm Chart to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ChartProductSlug, ChartProductVersion)
		if err != nil {
			return err
		}

		product.SetDeploymentType(models.DeploymentTypeHelm)
		product.PrepForUpdate()

		chartURL, err := url.Parse(ChartURL)
		if err != nil {
			return fmt.Errorf("failed to parse chart URL: %w", err)
		}

		var chart *models.ChartVersion
		if chartURL.Scheme == "" || chartURL.Scheme == "file" {
			chart, err = pkg.LoadChart(ChartURL)
			if err != nil {
				return err
			}

			err = GetUploadCredentials(cmd, args)
			if err != nil {
				return err
			}

			uploader := Marketplace.GetUploader(product.PublisherDetails.OrgId, UploadCredentials)
			_, uploadedChartUrl, err := uploader.UploadProductFile(ChartURL)
			if err != nil {
				return err
			}

			chart.HelmTarUrl = uploadedChartUrl
		} else if chartURL.Scheme == "http" || chartURL.Scheme == "https" {
			chart, err = Marketplace.DownloadChart(chartURL)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported protocol scheme: %s", chartURL.Scheme)
		}

		chart.AppVersion = version.Number
		chart.Readme = ChartReadme

		product.AddChart(chart)
		updatedProduct, err := Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Charts for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderCharts(updatedProduct.GetChartsForVersion(version.Number))
	},
}
