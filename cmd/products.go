// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
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

	ListProductsCmd.Flags().StringVar(&searchTerm, "search-text", "", "Filter by text")
	ListProductsCmd.Flags().BoolVarP(&allOrgs, "all-orgs", "a", false, "Show products from all organizations")

	GetProductCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = GetProductCmd.MarkFlagRequired("product")
}

var ProductCmd = &cobra.Command{
	Use:       "product",
	Aliases:   []string{"products"},
	Short:     "Commands related to products in general",
	Long:      "Interact with products in the Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"get", "list"},
}

var ListProductsCmd = &cobra.Command{
	Use:   "list",
	Short: "List products",
	Long:  "Lists products in the Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		products, err := Marketplace.ListProducts(allOrgs, searchTerm)
		if err != nil {
			return err
		}

		return Output.RenderProducts(products)
	},
}

var GetProductCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a product",
	Long:  "Gets details for a single product in the Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		return Output.RenderProduct(product)
	},
}
