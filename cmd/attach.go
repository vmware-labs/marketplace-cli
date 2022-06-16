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
	AttachContainerImageFile    string
	AttachContainerImageTag     string
	AttachContainerImageTagType string

	AttachMetaFile        string
	AttachMetaFileVersion string

	AttachVMFile string

	AttachInstructions string

	AttachPCAFile string
)

func init() {
	rootCmd.AddCommand(AttachCmd)
	AttachCmd.AddCommand(AttachChartCmd)
	AttachCmd.AddCommand(AttachContainerImageCmd)
	AttachCmd.AddCommand(AttachMetaFileCmd)
	AttachCmd.AddCommand(AttachVMCmd)

	AttachChartCmd.Flags().StringVarP(&AttachProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachChartCmd.MarkFlagRequired("product")
	AttachChartCmd.Flags().StringVarP(&AttachProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachChartCmd.Flags().StringVarP(&AttachChartURL, "chart", "c", "", "Path to to chart, either local tgz or public URL (required)")
	_ = AttachChartCmd.MarkFlagRequired("chart")
	AttachChartCmd.Flags().StringVar(&AttachInstructions, "instructions", "", "Chart deployment instructions (required)")
	_ = AttachChartCmd.MarkFlagRequired("instructions")
	AttachChartCmd.Flags().BoolVar(&AttachCreateVersion, "create-version", false, "Create the product version, if it doesn't already exist")
	AttachChartCmd.Flags().StringVar(&AttachPCAFile, "pca-file", "", "Path to a PCA file to upload")

	AttachContainerImageCmd.Flags().StringVarP(&AttachProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachContainerImageCmd.MarkFlagRequired("product")
	AttachContainerImageCmd.Flags().StringVarP(&AttachProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachContainerImageCmd.Flags().StringVarP(&AttachContainerImage, "image-repository", "r", "", "Image repository (e.g. registry/repository/image) (required)")
	_ = AttachContainerImageCmd.MarkFlagRequired("image-repository")
	AttachContainerImageCmd.Flags().StringVarP(&AttachContainerImageFile, "file", "f", "", "Path to a local tar file to upload")
	AttachContainerImageCmd.Flags().StringVar(&AttachContainerImageTag, "tag", "", "Image repository tag (required)")
	_ = AttachContainerImageCmd.MarkFlagRequired("tag")
	AttachContainerImageCmd.Flags().StringVar(&AttachContainerImageTagType, "tag-type", "", "Image repository tag type (fixed or floating) (required)")
	_ = AttachContainerImageCmd.MarkFlagRequired("tag-type")
	AttachContainerImageCmd.Flags().StringVarP(&AttachInstructions, "instructions", "i", "", "Image deployment instructions (required)")
	_ = AttachContainerImageCmd.MarkFlagRequired("instructions")
	AttachContainerImageCmd.Flags().BoolVar(&AttachCreateVersion, "create-version", false, "Create the product version, if it doesn't already exist")
	AttachContainerImageCmd.Flags().StringVar(&AttachPCAFile, "pca-file", "", "Path to a PCA file to upload")

	AttachMetaFileCmd.Flags().StringVarP(&AttachProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachMetaFileCmd.MarkFlagRequired("product")
	AttachMetaFileCmd.Flags().StringVarP(&AttachProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachMetaFileCmd.Flags().StringVar(&AttachMetaFile, "metafile", "", "Meta file to upload (required)")
	_ = AttachMetaFileCmd.MarkFlagRequired("metafile")
	AttachMetaFileCmd.Flags().StringVar(&MetaFileType, "metafile-type", "", "Meta file version (required, one of "+strings.Join(metaFileTypesList(), ", ")+")")
	AttachMetaFileCmd.Flags().StringVar(&AttachMetaFileVersion, "metafile-version", "", "Meta file type (default is the product version)")

	AttachVMCmd.Flags().StringVarP(&AttachProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachVMCmd.MarkFlagRequired("product")
	AttachVMCmd.Flags().StringVarP(&AttachProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachVMCmd.Flags().StringVar(&AttachVMFile, "file", "", "Virtual machine file to upload (required)")
	_ = AttachVMCmd.MarkFlagRequired("file")
	AttachVMCmd.Flags().BoolVar(&AttachCreateVersion, "create-version", false, "Create the product version, if it doesn't already exist")
	AttachVMCmd.Flags().StringVar(&AttachPCAFile, "pca-file", "", "Path to a PCA file to upload")
}

var AttachCmd = &cobra.Command{
	Use:       "attach",
	Short:     "Attach assets to a product",
	Long:      "Attach assets to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{AttachChartCmd.Use, AttachContainerImageCmd.Use, AttachVMCmd.Use},
}

var AttachChartCmd = &cobra.Command{
	Use:     "chart",
	Short:   "Attach a chart",
	Long:    "Attaches a Helm Chart to a product in the VMware Marketplace",
	Example: fmt.Sprintf("%s attach chart -p hyperspace-database-chart1 -v 1.2.3 --chart hyperspace-db-1.2.3.tgz --instructions \"helm install...\"", AppName),
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

		if AttachPCAFile != "" {
			uploader, err := Marketplace.GetUploader(product.PublisherDetails.OrgId)
			if err != nil {
				return err
			}
			_, pcaUrl, err := uploader.UploadMediaFile(AttachPCAFile)
			if err != nil {
				return err
			}

			product.SetPCAFile(version.Number, pcaUrl)
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
	Example: fmt.Sprintf("%s attach image -p hyperspace-database-image1 -v 1.2.3 --image hyperspace-labs/hyperspace-db --tag 1.2.3 --tag-type fixed --instructions \"docker run...\"", AppName),
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

		if AttachPCAFile != "" {
			uploader, err := Marketplace.GetUploader(product.PublisherDetails.OrgId)
			if err != nil {
				return err
			}
			_, pcaUrl, err := uploader.UploadMediaFile(AttachPCAFile)
			if err != nil {
				return err
			}

			product.SetPCAFile(version.Number, pcaUrl)
		}

		var updatedProduct *models.Product
		if AttachContainerImageFile != "" {
			updatedProduct, err = Marketplace.AttachLocalContainerImage(AttachContainerImageFile, AttachContainerImage, AttachContainerImageTag, AttachContainerImageTagType, AttachInstructions, product, version)
		} else {
			updatedProduct, err = Marketplace.AttachPublicContainerImage(AttachContainerImage, AttachContainerImageTag, AttachContainerImageTagType, AttachInstructions, product, version)
		}
		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Container images for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderContainerImages(updatedProduct.GetContainerImagesForVersion(version.Number))
	},
}

var AttachMetaFileCmd = &cobra.Command{
	Use:     "metafile",
	Short:   "Attach a meta file",
	Long:    "Upload and attach a meta file to a product in the VMware Marketplace",
	Example: fmt.Sprintf("%s attach metafile -p hyperspace-database-vm1 -v 1.2.3 --metafile deploy.sh --metafile-type cli", AppName),
	Args:    cobra.NoArgs,
	PreRunE: RunSerially(ValidateMetaFileType, GetRefreshToken),
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

		if AttachMetaFileVersion == "" {
			AttachMetaFileVersion = version.Number
		}
		updatedProduct, err := Marketplace.AttachMetaFile(AttachMetaFile, metaFileTypeMapping[MetaFileType], AttachMetaFileVersion, product, version)
		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Assets for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderAssets(pkg.GetAssets(updatedProduct, version.Number))
	},
}

var AttachVMCmd = &cobra.Command{
	Use:     "vm",
	Short:   "Attach a virtual machine file (ISO or OVA)",
	Long:    "Upload and attach a virtual machine file (ISO or OVA) to a product in the VMware Marketplace",
	Example: fmt.Sprintf("%s attach vm -p hyperspace-database-vm1 -v 1.2.3 --file hyperspace-db-1.2.3.iso", AppName),
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

		if AttachPCAFile != "" {
			uploader, err := Marketplace.GetUploader(product.PublisherDetails.OrgId)
			if err != nil {
				return err
			}
			_, pcaUrl, err := uploader.UploadMediaFile(AttachPCAFile)
			if err != nil {
				return err
			}

			product.SetPCAFile(version.Number, pcaUrl)
		}

		updatedProduct, err := Marketplace.UploadVM(AttachVMFile, product, version)
		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Virtual machine files for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderFiles(updatedProduct.GetFilesForVersion(version.Number))
	},
}
