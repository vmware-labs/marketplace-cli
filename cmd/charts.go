// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	chartID                 string
	downloadedChartFilename string
)

func init() {
	rootCmd.AddCommand(ChartCmd)
	ChartCmd.AddCommand(ListChartsCmd)
	ChartCmd.AddCommand(GetChartCmd)
	ChartCmd.AddCommand(DownloadChartCmd)
	ChartCmd.AddCommand(CreateChartCmd)

	ChartCmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = ChartCmd.MarkPersistentFlagRequired("product")
	ChartCmd.PersistentFlags().StringVarP(&ProductVersion, "product-version", "v", "latest", "Product version")

	GetChartCmd.Flags().StringVar(&chartID, "chart-id", "", "chart ID")

	DownloadChartCmd.Flags().StringVar(&chartID, "chart-id", "", "chart ID")
	DownloadChartCmd.Flags().StringVarP(&downloadedChartFilename, "filename", "f", "chart.tgz", "output file name")

	CreateChartCmd.Flags().StringVarP(&ChartVersion, "chart-version", "e", "", "chart version")
	_ = CreateChartCmd.MarkFlagRequired("chart-version")
	CreateChartCmd.Flags().StringVarP(&ChartURL, "chart-url", "u", "", "url to chart tgz")
	_ = CreateChartCmd.MarkFlagRequired("chart-url")
	CreateChartCmd.Flags().StringVar(&ChartRepositoryURL, "repository-url", "", "chart public repository url")
	_ = CreateChartCmd.MarkFlagRequired("repository-url")
	CreateChartCmd.Flags().StringVar(&ChartRepositoryName, "repository-name", "", "chart public repository name")
	_ = CreateChartCmd.MarkFlagRequired("repository-name")
	CreateChartCmd.Flags().StringVar(&DeploymentInstructions, "readme", "", "readme information")
	_ = CreateChartCmd.MarkFlagRequired("readme")
}

var ChartCmd = &cobra.Command{
	Use:       "chart",
	Aliases:   []string{"charts"},
	Short:     "List and manage Helm charts attached to a product",
	Long:      "List and manage Helm charts attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"list", "get", "download", "create"},
}

var ListChartsCmd = &cobra.Command{
	Use:   "list",
	Short: "List product charts",
	Long:  "Prints the list of Helm charts attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		charts := product.GetChartsForVersion(version.Number)
		if len(charts) == 0 {
			fmt.Printf("%s %s does not have any charts\n", product.Slug, version.Number)
			return nil
		}

		return Output.RenderCharts(charts)
	},
}

var GetChartCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details for a chart",
	Long:  "Prints detailed information about a Helm chart attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		var chart *models.ChartVersion
		if chartID != "" {
			chart = product.GetChart(chartID)
			if chart == nil {
				return fmt.Errorf("product \"%s\" %s does not have the container image \"%s\"", ProductSlug, version.Number, ImageRepository)
			}
		} else {
			charts := product.GetChartsForVersion(version.Number)
			if len(charts) == 0 {
				return fmt.Errorf("product \"%s\" does not have any charts for version %s", ProductSlug, version.Number)
			} else if len(charts) == 1 {
				chart = charts[0]
			} else {
				return fmt.Errorf("multiple container images found for %s for version %s, please use the --image-repository parameter", ProductSlug, version.Number)
			}
		}

		return Output.RenderChart(chart)
	},
}

var DownloadChartCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a chart",
	Long:  "Downloads a Helm chart attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		var chart *models.ChartVersion
		if chartID != "" {
			chart = product.GetChart(chartID)
			if chart == nil {
				return fmt.Errorf("product \"%s\" %s does not have the container image \"%s\"", ProductSlug, version.Number, ImageRepository)
			}
		} else {
			charts := product.GetChartsForVersion(version.Number)
			if len(charts) == 0 {
				return fmt.Errorf("product \"%s\" does not have any charts for version %s", ProductSlug, version.Number)
			} else if len(charts) == 1 {
				chart = charts[0]
			} else {
				return fmt.Errorf("multiple container images found for %s for version %s, please use the --image-repository parameter", ProductSlug, version.Number)
			}
		}

		if !chart.IsUpdatedInMarketplaceRegistry {
			return fmt.Errorf("%s %s in %s %s is not downloadable: %s", chart.TarUrl, chart.Version, ProductSlug, version.Number, chart.ValidationStatus)
		}

		cmd.Printf("Downloading chart to %s...\n", downloadedChartFilename)
		return Marketplace.Download(product.ProductId, downloadedChartFilename, &pkg.DownloadRequestPayload{
			AppVersion:   chart.AppVersion,
			ChartVersion: chart.Version,
			EulaAccepted: true,
		}, cmd.ErrOrStderr())
	},
}

var CreateChartCmd = &cobra.Command{
	Use:   "create",
	Short: "Attach a chart",
	Long:  "Attaches a Helm Chart to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		product.SetDeploymentType(models.DeploymentTypeHelm)
		product.PrepForUpdate()
		product.AddChart(&models.ChartVersion{
			TarUrl:     ChartURL,
			Version:    ChartVersion,
			AppVersion: version.Number,
			Readme:     DeploymentInstructions,
			Repo: &models.Repo{
				Name: ChartRepositoryName,
				Url:  ChartRepositoryURL,
			},
		})

		updatedProduct, err := Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		return Output.RenderCharts(updatedProduct.GetChartsForVersion(version.Number))
	},
}
