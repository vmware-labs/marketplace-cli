// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	. "github.com/vmware-labs/marketplace-cli/v2/lib"
	"github.com/vmware-labs/marketplace-cli/v2/models"
)

func init() {
	rootCmd.AddCommand(ProductVersionCmd)
	ProductVersionCmd.AddCommand(ListProductVersionsCmd)
	ProductVersionCmd.AddCommand(GetProductVersionCmd)
	ProductVersionCmd.AddCommand(CreateProductVersionCmd)
	ProductVersionCmd.PersistentFlags().StringVarP(&OutputFormat, "OutputFormat", "f", FormatTable, "Output OutputFormat")

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
		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		err = RenderVersions(OutputFormat, response.Response.Data, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return errors.Wrap(err, "failed to render the product version")
		}

		return nil
	},
}

var GetProductVersionCmd = &cobra.Command{
	Use:  "get",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		product := response.Response.Data
		version := product.GetVersion(ProductVersion)
		if version == nil {
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" does not have a version %s", ProductSlug, ProductVersion)
		}

		err = RenderVersion(OutputFormat, ProductVersion, product, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return errors.Wrap(err, "failed to render the product version")
		}

		return nil
	},
}

var CreateProductVersionCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		product := response.Response.Data

		if product.HasVersion(ProductVersion) {
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" already has version %s", ProductSlug, ProductVersion)
		}
		product.Versions = append(product.Versions, &models.Version{
			Number: ProductVersion,
		})

		response = &GetProductResponse{}
		err = PutProduct(product, true, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		err = RenderVersions(OutputFormat, response.Response.Data, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return errors.Wrap(err, "failed to render the product version")
		}

		return nil
	},
}
