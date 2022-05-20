// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	AttachProductSlug    string
	AttachProductVersion string
	AttachCreateVersion  bool

	AttachChartURL string

	AttachContainerImage        string
	AttachContainerImageTag     string
	AttachContainerImageTagType string

	AttachVMFile string

	AttachInstructions string
)

func init() {
	rootCmd.AddCommand(AttachCmd)
	AttachCmd.AddCommand(AttachChartCmd)
	AttachCmd.AddCommand(AttachContainerImageCmd)
	AttachCmd.AddCommand(AttachVMCmd)

	AttachChartCmd.Flags().StringVarP(&AttachProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachChartCmd.MarkFlagRequired("product")
	AttachChartCmd.Flags().StringVarP(&AttachProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachChartCmd.Flags().StringVarP(&AttachChartURL, "chart", "c", "", "path to to chart, either local tgz or public URL (required)")
	_ = AttachChartCmd.MarkFlagRequired("chart")
	AttachChartCmd.Flags().StringVar(&AttachInstructions, "instructions", "", "readme information")
	_ = AttachChartCmd.MarkFlagRequired("instructions")
	AttachChartCmd.Flags().BoolVar(&AttachCreateVersion, "create-version", false, "create the product version, if it doesn't already exist")

	AttachContainerImageCmd.Flags().StringVarP(&AttachProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachContainerImageCmd.MarkFlagRequired("product")
	AttachContainerImageCmd.Flags().StringVarP(&AttachProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachContainerImageCmd.Flags().StringVarP(&AttachContainerImage, "image-repository", "r", "", "container repository")
	_ = AttachContainerImageCmd.MarkFlagRequired("image-repository")
	AttachContainerImageCmd.Flags().StringVar(&AttachContainerImageTag, "tag", "", "container repository tag")
	_ = AttachContainerImageCmd.MarkFlagRequired("tag")
	AttachContainerImageCmd.Flags().StringVar(&AttachContainerImageTagType, "tag-type", "", "container repository tag type (fixed or floating)")
	_ = AttachContainerImageCmd.MarkFlagRequired("tag-type")
	AttachContainerImageCmd.Flags().StringVarP(&AttachInstructions, "deployment-instructions", "i", "", "deployment instructions")
	_ = AttachContainerImageCmd.MarkFlagRequired("deployment-instructions")
	AttachContainerImageCmd.Flags().BoolVar(&AttachCreateVersion, "create-version", false, "create the product version, if it doesn't already exist")

	AttachVMCmd.Flags().StringVarP(&AttachProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachVMCmd.MarkFlagRequired("product")
	AttachVMCmd.Flags().StringVarP(&AttachProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachVMCmd.Flags().StringVar(&AttachVMFile, "file", "", "Virtual machine file to upload (required)")
	_ = AttachVMCmd.MarkFlagRequired("file")
	AttachVMCmd.Flags().BoolVar(&AttachCreateVersion, "create-version", false, "create the product version, if it doesn't already exist")
}

var AttachCmd = &cobra.Command{
	Use:       "attach",
	Short:     "Attach assets to a product",
	Long:      "List and manage virtual machine files attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{AttachChartCmd.Use, AttachContainerImageCmd.Use, AttachVMCmd.Use},
}

var AttachChartCmd = &cobra.Command{
	Use:     "chart",
	Short:   "Attach a chart",
	Long:    "Attaches a Helm Chart to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(AttachProductSlug, AttachProductVersion)
		if err != nil {
			if errors.Is(err, &pkg.VersionDoesNotExistError{}) && AttachCreateVersion {
				version = product.NewVersion(AttachProductVersion)
			} else {
				return err
			}
		}

		chartURL, err := url.Parse(AttachChartURL)
		if err != nil {
			return fmt.Errorf("failed to parse chart URL: %w", err)
		}

		var updatedProduct *models.Product
		if chartURL.Scheme == "" || chartURL.Scheme == "file" {
			updatedProduct, err = Marketplace.AttachLocalChart(AttachChartURL, AttachInstructions, product, version)
		} else if chartURL.Scheme == "http" || chartURL.Scheme == "https" {
			updatedProduct, err = Marketplace.AttachPublicChart(chartURL, AttachInstructions, product, version)
		} else {
			return fmt.Errorf("unsupported protocol scheme: %s", chartURL.Scheme)
		}

		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Charts for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderCharts(updatedProduct.GetChartsForVersion(version.Number))
	},
}

func ValidateTagType(cmd *cobra.Command, args []string) error {
	AttachContainerImageTagType = strings.ToUpper(AttachContainerImageTagType)
	if AttachContainerImageTagType != models.ImageTagTypeFixed && AttachContainerImageTagType != models.ImageTagTypeFloating {
		return fmt.Errorf("invalid image tag type: %s. must be either \"%s\" or \"%s\"", AttachContainerImageTagType, models.ImageTagTypeFixed, models.ImageTagTypeFloating)
	}
	return nil
}

var AttachContainerImageCmd = &cobra.Command{
	Use:     "image",
	Short:   "Attach a container image",
	Long:    "Attaches a container image to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: RunSerially(ValidateTagType, GetRefreshToken),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(AttachProductSlug, AttachProductVersion)
		if err != nil {
			if errors.Is(err, &pkg.VersionDoesNotExistError{}) && AttachCreateVersion {
				version = product.NewVersion(AttachProductVersion)
			} else {
				return err
			}
		}

		updatedProduct, err := Marketplace.AttachPublicContainerImage(AttachContainerImage, AttachContainerImageTag, AttachContainerImageTagType, AttachInstructions, product, version)
		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Container images for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderContainerImages(updatedProduct.GetContainerImagesForVersion(version.Number))
	},
}

var AttachVMCmd = &cobra.Command{
	Use:     "vm",
	Short:   "Upload and attach a virtual machine file (ISO or OVA)",
	Long:    "Uploads and attaches a virtual machine file (ISO or OVA) to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(AttachProductSlug, AttachProductVersion)
		if err != nil {
			if errors.Is(err, &pkg.VersionDoesNotExistError{}) && AttachCreateVersion {
				version = product.NewVersion(AttachProductVersion)
			} else {
				return err
			}
		}

		updatedProduct, err := Marketplace.UploadVM(AttachVMFile, product, version)
		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Virtual machine files for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderFiles(updatedProduct.GetFilesForVersion(version.Number))
	},
}
