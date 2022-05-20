// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

var (
	ChartID             string
	ChartReadme         string
	ChartProductSlug    string
	ChartProductVersion string
	ChartURL            string
	ChartCreateVersion  bool
)

func init() {
	rootCmd.AddCommand(ChartCmd)
	ChartCmd.AddCommand(ListChartsCmd)
	ChartCmd.AddCommand(GetChartCmd)
	ChartCmd.AddCommand(OldAttachChartCmd)

	ListChartsCmd.Flags().StringVarP(&ChartProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListChartsCmd.MarkFlagRequired("product")
	ListChartsCmd.Flags().StringVarP(&ChartProductVersion, "product-version", "v", "", "Product version (default to latest version)")

	GetChartCmd.Flags().StringVarP(&ChartProductSlug, "product", "p", "", "Product slug (required)")
	_ = GetChartCmd.MarkFlagRequired("product")
	GetChartCmd.Flags().StringVarP(&ChartProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	GetChartCmd.Flags().StringVar(&ChartID, "chart-id", "", "chart ID")

	OldAttachChartCmd.Flags().StringVarP(&ChartProductSlug, "product", "p", "", "Product slug (required)")
	_ = OldAttachChartCmd.MarkFlagRequired("product")
	OldAttachChartCmd.Flags().StringVarP(&ChartProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	OldAttachChartCmd.Flags().StringVarP(&ChartURL, "chart", "c", "", "path to to chart, either local tgz or public URL (required)")
	_ = OldAttachChartCmd.MarkFlagRequired("chart")
	OldAttachChartCmd.Flags().StringVar(&ChartReadme, "readme", "", "readme information")
	_ = OldAttachChartCmd.MarkFlagRequired("readme")
	OldAttachChartCmd.Flags().BoolVar(&ChartCreateVersion, "create-version", false, "create the product version, if it doesn't already exist")
}

var ChartCmd = &cobra.Command{
	Use:       "chart",
	Aliases:   []string{"charts"},
	Short:     "List and manage Helm charts attached to a product",
	Long:      "List and manage Helm charts attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListChartsCmd.Use, GetChartCmd.Use, OldAttachChartCmd.Use},
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

var OldAttachChartCmd = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a chart",
	Long:    "Attaches a Helm Chart to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.PrintErrln("mkpcli chart attach has been deprecated and will be removed in the next major version. Please use mkpcli attach chart instead.")
		AttachProductSlug = ChartProductSlug
		AttachProductVersion = ChartProductVersion
		AttachCreateVersion = ChartCreateVersion
		AttachChartURL = ChartURL
		AttachInstructions = ChartReadme
		return AttachChartCmd.RunE(cmd, args)
	},
}
