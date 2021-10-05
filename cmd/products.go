// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

var (
	allOrgs        = false
	searchTerm     string
	ProductSlug    string
	ProductVersion string
)

func init() {
	rootCmd.AddCommand(ProductCmd)
	ProductCmd.AddCommand(ListProductsCmd)
	ProductCmd.AddCommand(GetProductCmd)
	ProductCmd.AddCommand(AddProductVersionCmd)
	ProductCmd.AddCommand(ListProductVersionsCmd)

	ListProductsCmd.Flags().StringVar(&searchTerm, "search-text", "", "Filter product list by text")
	ListProductsCmd.Flags().BoolVarP(&allOrgs, "all-orgs", "a", false, "Show products from all organizations")

	GetProductCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = GetProductCmd.MarkFlagRequired("product")

	AddProductVersionCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = AddProductVersionCmd.MarkFlagRequired("product")
	AddProductVersionCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version (required)")
	_ = AddProductVersionCmd.MarkFlagRequired("product-version")

	ListProductVersionsCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListProductVersionsCmd.MarkFlagRequired("product")
}

var ProductCmd = &cobra.Command{
	Use:       "product",
	Aliases:   []string{"products"},
	Short:     "Get information about products",
	Long:      "Get information about products in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListProductsCmd.Use, GetProductCmd.Use},
}

var ListProductsCmd = &cobra.Command{
	Use:   "list",
	Short: "List products",
	Long:  "List and search for products in the VMware Marketplace",
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
	Short: "Show details about a product",
	Long:  "Show details about a product in the VMware Marketplace",
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

var AddProductVersionCmd = &cobra.Command{
	Use:   "add-version",
	Short: "Add a new version",
	Long:  "Adds a new version to the given product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		if product.HasVersion(ProductVersion) {
			return fmt.Errorf("product \"%s\" already has version %s", ProductSlug, ProductVersion)
		}
		product.Versions = append(product.AllVersions, &models.Version{
			Number: ProductVersion,
		})

		product.PrepForUpdate()
		updatedProduct, err := Marketplace.PutProduct(product, true)
		if err != nil {
			return err
		}

		return Output.RenderVersions(updatedProduct)
	},
}

var ListProductVersionsCmd = &cobra.Command{
	Use:   "list-versions",
	Short: "List product versions",
	Long:  "Prints the list of versions for the given product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		return Output.RenderVersions(product)
	},
}
