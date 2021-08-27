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
	ProductVersionCmd.AddCommand(GetProductVersionCmd)
	ProductVersionCmd.AddCommand(CreateProductVersionCmd)
	ProductVersionCmd.PersistentFlags().StringVarP(&OutputFormat, "output-format", "f", FormatTable, "Output format")

	ProductVersionCmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = ProductVersionCmd.MarkPersistentFlagRequired("product")

	GetProductVersionCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	_ = GetProductVersionCmd.MarkFlagRequired("product-version")

	CreateProductVersionCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	_ = CreateProductVersionCmd.MarkFlagRequired("product-version")
}

var ProductVersionCmd = &cobra.Command{
	Use:               "product-version",
	Aliases:           []string{"product-versions"},
	Short:             "product versions",
	Long:              "",
	Args:              cobra.OnlyValidArgs,
	ValidArgs:         []string{"get", "list", "create"},
	PersistentPreRunE: GetRefreshToken,
}

var ListProductVersionsCmd = &cobra.Command{
	Use:  "list",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		err = RenderVersions(OutputFormat, product, cmd.OutOrStdout())
		if err != nil {
			return fmt.Errorf("failed to render the product version: %w", err)
		}

		return nil
	},
}

var GetProductVersionCmd = &cobra.Command{
	Use:  "get",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		version := product.GetVersion(ProductVersion)
		if version == nil {
			return fmt.Errorf("product \"%s\" does not have a version %s", ProductSlug, ProductVersion)
		}

		err = RenderVersion(OutputFormat, ProductVersion, product, cmd.OutOrStdout())
		if err != nil {
			return fmt.Errorf("failed to render the product version: %w", err)
		}

		return nil
	},
}

var CreateProductVersionCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.NoArgs,
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

		updatedProduct, err := Marketplace.PutProduct(product, true)
		if err != nil {
			return err
		}

		err = RenderVersions(OutputFormat, updatedProduct, cmd.OutOrStdout())
		if err != nil {
			return fmt.Errorf("failed to render the product version: %w", err)
		}

		return nil
	},
}
