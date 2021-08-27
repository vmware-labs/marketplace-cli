// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func init() {
	rootCmd.AddCommand(ChartCmd)
	ChartCmd.AddCommand(ListChartsCmd)
	ChartCmd.AddCommand(CreateChartCmd)
	ChartCmd.PersistentFlags().StringVarP(&OutputFormat, "output-format", "f", FormatTable, "Output format")

	ChartCmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = ChartCmd.MarkPersistentFlagRequired("product")
	ChartCmd.PersistentFlags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	_ = ChartCmd.MarkPersistentFlagRequired("product-version")

	CreateChartCmd.Flags().StringVarP(&ChartName, "chart-name", "n", "", "chart name")
	_ = CreateChartCmd.MarkFlagRequired("chart-name")
	CreateChartCmd.Flags().StringVarP(&ChartVersion, "chart-version", "e", "", "chart version")
	_ = CreateChartCmd.MarkFlagRequired("chart-version")
	CreateChartCmd.Flags().StringVarP(&ChartURL, "chart-url", "u", "", "url to chart tgz")
	_ = CreateChartCmd.MarkFlagRequired("chart-url")
	CreateChartCmd.Flags().StringVar(&ChartRepositoryName, "repository-url", "", "chart public repository url")
	_ = CreateChartCmd.MarkFlagRequired("repository-url")
	CreateChartCmd.Flags().StringVar(&ChartRepositoryURL, "repository-name", "", "chart public repository name")
	_ = CreateChartCmd.MarkFlagRequired("repository-name")
	CreateChartCmd.Flags().StringVarP(&DeploymentInstructions, "deployment-instructions", "i", "", "deployment instructions")
}

var ChartCmd = &cobra.Command{
	Use:               "chart",
	Aliases:           []string{"charts"},
	Short:             "charts",
	Long:              "",
	Args:              cobra.OnlyValidArgs,
	ValidArgs:         []string{"get", "list", "create"},
	PersistentPreRunE: GetRefreshToken,
}

var ListChartsCmd = &cobra.Command{
	Use:   "list",
	Short: "list charts",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		if !product.HasVersion(ProductVersion) {
			return fmt.Errorf("product \"%s\" does not have a version %s", ProductSlug, ProductVersion)
		}
		charts := product.GetChartsForVersion(ProductVersion)
		if len(charts) == 0 {
			cmd.Printf("product \"%s\" %s does not have any charts\n", ProductSlug, ProductVersion)
			return nil
		}

		err = RenderCharts(OutputFormat, charts, cmd.OutOrStdout())
		if err != nil {
			return fmt.Errorf("failed to render the charts: %w", err)
		}

		return nil
	},
}

var CreateChartCmd = &cobra.Command{
	Use:   "create",
	Short: "add a chart to a product",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		if !product.HasVersion(ProductVersion) {
			return fmt.Errorf("product \"%s\" does not have a version %s, please add it first", ProductSlug, ProductVersion)
		}

		product.SetDeploymentType(models.DeploymentTypeHelm)

		charts := product.GetChartsForVersion(ProductVersion)
		product.ChartVersions = append(charts, &models.ChartVersion{
			Id:             ChartName,
			TarUrl:         ChartURL,
			Version:        ChartVersion,
			AppVersion:     ProductVersion,
			InstallOptions: DeploymentInstructions,
			Repo: &models.Repo{
				Name: ChartRepositoryName,
				Url:  ChartRepositoryURL,
			},
		})

		updatedProduct, err := Marketplace.PutProduct(product, false)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		err = RenderCharts(OutputFormat, updatedProduct.ChartVersions, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to render the charts: %w", err)
		}
		return nil
	},
}
