// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

func init() {
	rootCmd.AddCommand(ProductVersionCmd)
	ProductVersionCmd.AddCommand(ListProductVersionsCmd)
	ProductVersionCmd.AddCommand(CreateProductVersionCmd)

	ProductVersionCmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = ProductVersionCmd.MarkPersistentFlagRequired("product")

	CreateProductVersionCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	_ = ProductVersionCmd.MarkPersistentFlagRequired("product-version")
}

var ProductVersionCmd = &cobra.Command{
	Use:       "product-version",
	Aliases:   []string{"product-versions"},
	Short:     "List and manage versions of a product",
	Long:      "List and manage versions of a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"list", "add"},
}

var ListProductVersionsCmd = &cobra.Command{
	Use:   "list",
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

var CreateProductVersionCmd = &cobra.Command{
	Use:   "add",
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
