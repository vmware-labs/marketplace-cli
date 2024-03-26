// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	ProductSlug           string
	ProductVersion        string
	ListProductsAllOrgs   = false
	ListProductsOrgId     string
	ListProductSearchText string
	SetOSLFile            string
)

func init() {
	rootCmd.AddCommand(ProductCmd)
	ProductCmd.AddCommand(ListProductsCmd)
	ProductCmd.AddCommand(GetProductCmd)
	ProductCmd.AddCommand(ListAssetsCmd)
	ProductCmd.AddCommand(ListProductVersionsCmd)
	ProductCmd.AddCommand(SetCmd)
	ProductCmd.AddCommand(StartProcessCmd)
	ProductCmd.AddCommand(StopProcessCmd)

	ListProductsCmd.Flags().StringVar(&ListProductSearchText, "search-text", "", "Filter product list by text")
	ListProductsCmd.Flags().BoolVarP(&ListProductsAllOrgs, "all-orgs", "a", false, "Show published products from all organizations")
	ListProductsCmd.Flags().StringVar(&ListProductsOrgId, "org-id", "", "Filter product list by organization id")

	GetProductCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = GetProductCmd.MarkFlagRequired("product")
	GetProductCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")

	ListAssetsCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListAssetsCmd.MarkFlagRequired("product")
	ListAssetsCmd.Flags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	ListAssetsCmd.Flags().StringVarP(&AssetType, "type", "t", "", "Filter assets by type (one of "+strings.Join(assetTypesList(), ", ")+")")

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
	Short:     "Manage products",
	Long:      "Manage products in the VMware Marketplace",
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
		filter := &pkg.ListProductFilter{
			Text:    ListProductSearchText,
			AllOrgs: ListProductsAllOrgs,
		}
		if ListProductsOrgId != "" {
			filter.OrgIds = []string{ListProductsOrgId}
		}

		products, err := Marketplace.ListProducts(filter)
		if err != nil {
			return err
		}

		header := "All products"
		if filter.AllOrgs {
			header += " from all organizations"
		} else if len(products) > 0 {
			header += fmt.Sprintf(" from %s", products[0].PublisherDetails.OrgDisplayName)
		}
		if filter.Text != "" {
			header += fmt.Sprintf(" filtered by \"%s\"", filter.Text)
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

var ListAssetsCmd = &cobra.Command{
	Use:     "list-assets",
	Short:   "List attached assets",
	Long:    "Print the list of assets attached to the given product",
	Args:    cobra.NoArgs,
	PreRunE: RunSerially(ValidateAssetTypeFilter, GetRefreshToken),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		var assets []*pkg.Asset
		if AssetType == "" {
			assets = pkg.GetAssets(product, version.Number)
			Output.PrintHeader(fmt.Sprintf("Assets for %s %s:", product.DisplayName, version.Number))
		} else {
			assetType := assetTypeMapping[AssetType]
			assets = pkg.GetAssetsByType(assetType, product, version.Number)
			Output.PrintHeader(fmt.Sprintf("%s assets for %s %s:", assetType, product.DisplayName, version.Number))
		}

		return Output.RenderAssets(assets)
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
	Short:   "Modify product details",
	Long:    "Modify fields in a given product",
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
			uploader, err := Marketplace.GetUploader(product.PublisherDetails.OrgId)
			if err != nil {
				return err
			}
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

var StartProcessCmd = &cobra.Command{
	Use:     "start-mkpl-agent-process",
	Short:   "start background process required for Markeptlace Agent",
	Long:    "start background process performing for monitoring, polling, heart-beat and update subscription status",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		startBackgroundProcess()
		return nil
	},
}

var StopProcessCmd = &cobra.Command{
	Use:     "stop-mkpl-agent-process",
	Short:   "stop background process of Marketplace Agent",
	Long:    "stop background process performing monitoring, polling, heart-beat and update subscription status",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		stopBackgroundProcess()
		return nil
	},
}

func startBackgroundProcess() {
	log.Println("starting the agent process...")
	//dir, _ := os.Getwd()
	cmd := exec.Command("./agent")
	err := cmd.Start()
	if err != nil {
		log.Println("Error starting the process: ", err)
		os.Exit(1)
	}

	log.Println("Process has started successfully...")
}

func stopBackgroundProcess() {
	cmd := exec.Command("pkill", "-f", "agent")
	err := cmd.Start()
	if err != nil {
		log.Println("Error in stopping the process:", err)
		os.Exit(1)
	}

	log.Println("Process stopped successfully")
}
