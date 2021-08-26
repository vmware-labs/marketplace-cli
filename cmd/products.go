// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	allOrgs    = false
	searchTerm string
)

func init() {
	rootCmd.AddCommand(ProductCmd)
	ProductCmd.AddCommand(ListProductsCmd)
	ProductCmd.AddCommand(GetProductCmd)
	ProductCmd.PersistentFlags().StringVarP(&OutputFormat, "output-format", "f", FormatTable, "Output format")

	ListProductsCmd.Flags().StringVar(&searchTerm, "search-text", "", "Filter by text")
	ListProductsCmd.Flags().BoolVarP(&allOrgs, "all-orgs", "a", false, "Show products from all organizations")

	GetProductCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = GetProductCmd.MarkFlagRequired("product")
}

var ProductCmd = &cobra.Command{
	Use:               "product",
	Aliases:           []string{"products"},
	Short:             "stuff related to products",
	Long:              "",
	Args:              cobra.OnlyValidArgs,
	ValidArgs:         []string{"get", "list"},
	PersistentPreRunE: GetRefreshToken,
}

var ListProductsCmd = &cobra.Command{
	Use:  "list",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		products, err := Marketplace.ListProducts(allOrgs, searchTerm)
		if err != nil {
			return err
		}

		err = RenderProductList(OutputFormat, products, cmd.OutOrStdout())
		if err != nil {
			return fmt.Errorf("failed to render the list of products: %w", err)
		}

		return nil
	},
}

var GetProductCmd = &cobra.Command{
	Use:  "get [product slug]",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		err = RenderProduct(OutputFormat, product, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to render the product: %w", err)
		}
		return nil
	},
}
