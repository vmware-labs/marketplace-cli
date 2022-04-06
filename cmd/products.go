// Copyright 2022 VMware, Inc.
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
	SetOSLFile     string
)

func init() {
	rootCmd.AddCommand(ProductCmd)
	ProductCmd.AddCommand(ListProductsCmd)
	ProductCmd.AddCommand(GetProductCmd)
	ProductCmd.AddCommand(AddProductVersionCmd)
	ProductCmd.AddCommand(ListProductVersionsCmd)
	ProductCmd.AddCommand(SetCmd)

	ListProductsCmd.Flags().StringVar(&searchTerm, "search-text", "", "Filter product list by text")
	ListProductsCmd.Flags().BoolVarP(&allOrgs, "all-orgs", "a", false, "Show published products from all organizations")

	GetProductCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = GetProductCmd.MarkFlagRequired("product")
	GetProductCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")

	AddProductVersionCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = AddProductVersionCmd.MarkFlagRequired("product")
	AddProductVersionCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version (required)")
	_ = AddProductVersionCmd.MarkFlagRequired("product-version")

	ListProductVersionsCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListProductVersionsCmd.MarkFlagRequired("product")

	SetCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = SetCmd.MarkFlagRequired("product")
	SetCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version (required)")
	_ = SetCmd.MarkFlagRequired("product-version")
	SetCmd.Flags().StringVar(&SetOSLFile, "osl-file", "", "File with OSL disclosures")
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
	Long: "List and search for products in the VMware Marketplace\n" +
		"Default without --all-orgs is to list all products (including unpublished) from your organization",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		products, err := Marketplace.ListProducts(allOrgs, searchTerm)
		if err != nil {
			return err
		}

		header := "All products"
		if allOrgs {
			header += " from all organizations"
		} else if len(products) > 0 {
			header += fmt.Sprintf(" from %s", products[0].PublisherDetails.OrgDisplayName)
		}
		if searchTerm != "" {
			header += fmt.Sprintf(" filtered by \"%s\"", searchTerm)
		}

		Output.PrintHeader(header)
		return Output.RenderProducts(products)
	},
}

var GetProductCmd = &cobra.Command{
	Use:     "get",
	Short:   "Show details about a product",
	Long:    "Show details about a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if ProductVersion == "" {
			product, err := Marketplace.GetProduct(ProductSlug)
			if err != nil {
				return err
			}
			return Output.RenderProduct(product, product.GetLatestVersion())
		} else {
			product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
			if err != nil {
				return err
			}
			return Output.RenderProduct(product, version)
		}
	},
}

var AddProductVersionCmd = &cobra.Command{
	Use:     "add-version",
	Short:   "Add a new version",
	Long:    "Adds a new version to the given product",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
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

		Output.PrintHeader(fmt.Sprintf("Versions for %s:", updatedProduct.DisplayName))
		return Output.RenderVersions(updatedProduct)
	},
}

var ListProductVersionsCmd = &cobra.Command{
	Use:     "list-versions",
	Short:   "List product versions",
	Long:    "Prints the list of versions for the given product",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			return err
		}

		models.Sort(product.AllVersions)
		Output.PrintHeader(fmt.Sprintf("Versions for %s:", product.DisplayName))
		return Output.RenderVersions(product)
	},
}

var SetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Modify the product",
	Long:    "Modify the product",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		if SetOSLFile == "" {
			return fmt.Errorf("nothing specified to set")
		}
		cmd.SilenceUsage = true

		product, _, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}
		product.PrepForUpdate()

		if SetOSLFile != "" {
			err = GetUploadCredentials(cmd, args)
			if err != nil {
				return err
			}

			uploader := Marketplace.GetUploader(product.PublisherDetails.OrgId, UploadCredentials)
			_, oslUrl, err := uploader.UploadMediaFile(SetOSLFile)
			if err != nil {
				return err
			}

			product.OpenSourceDisclosure.LicenseDisclosureURL = oslUrl
		}

		_, err = Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		// TODO: what do we render?

		return nil
	},
}
