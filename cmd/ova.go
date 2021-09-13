// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	ovaFile               string
	downloadedOVAFilename string
)

func init() {
	rootCmd.AddCommand(OVACmd)
	OVACmd.AddCommand(ListOVACmd)
	OVACmd.AddCommand(GetOVACmd)
	OVACmd.AddCommand(DownloadOVACmd)
	OVACmd.AddCommand(CreateOVACmd)

	OVACmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = OVACmd.MarkPersistentFlagRequired("product")
	OVACmd.PersistentFlags().StringVarP(&ProductVersion, "product-version", "v", "latest", "Product version")

	GetOVACmd.Flags().StringVar(&ovaFile, "file-id", "", "The file ID of the file to get")

	DownloadOVACmd.Flags().StringVar(&ovaFile, "file-id", "", "The file ID of the file to download")
	DownloadOVACmd.Flags().StringVarP(&downloadedOVAFilename, "filename", "f", "", "Downloaded file name (default is original file name)")

	CreateOVACmd.Flags().StringVar(&ovaFile, "ova-file", "", "OVA file to upload")
}

var OVACmd = &cobra.Command{
	Use:       "ova",
	Aliases:   []string{"ovas"},
	Short:     "OVA related commands",
	Long:      "Interact with OVAs attached to a Marketplace product",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"list", "get", "download", "create"},
}

var ListOVACmd = &cobra.Command{
	Use:   "list",
	Short: "List OVAs",
	Long:  "List the OVAs attached to a product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		ovas := product.GetFilesForVersion(version.Number)
		if len(ovas) == 0 {
			cmd.Printf("product \"%s\" %s does not have any OVAs\n", product.Slug, version.Number)
			return nil
		}

		return Output.RenderOVAs(ovas)
	},
}

var GetOVACmd = &cobra.Command{
	Use:   "get",
	Short: "Get OVA details",
	Long:  "Get details for an OVA file attached to a product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		var file *models.ProductDeploymentFile
		if ovaFile != "" {
			file = product.GetFile(ovaFile)
			if file == nil {
				return fmt.Errorf("no file found with ID %s", ovaFile)
			}
		} else {
			files := product.GetFilesForVersion(version.Number)
			if len(files) == 0 {
				return fmt.Errorf("no files found for %s for version %s", ProductSlug, version.Number)
			} else if len(files) == 1 {
				file = files[0]
			} else {
				return fmt.Errorf("multiple files found for %s for version %s, please use the --file-id parameter", ProductSlug, version.Number)
			}
		}

		return Output.RenderOVA(file)
	},
}

var DownloadOVACmd = &cobra.Command{
	Use:   "download",
	Short: "Download an OVA",
	Long:  "Download an OVA attached to a product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		var file *models.ProductDeploymentFile
		if ovaFile != "" {
			file = product.GetFile(ovaFile)
			if file == nil {
				return fmt.Errorf("no file found with ID %s", ovaFile)
			}
		} else {
			files := product.GetFilesForVersion(version.Number)
			if len(files) == 0 {
				return fmt.Errorf("no files found for %s for version %s", ProductSlug, version.Number)
			} else if len(files) == 1 {
				file = files[0]
			} else {
				return fmt.Errorf("multiple files found for %s for version %s, please use the --file-id parameter", ProductSlug, version.Number)
			}
		}

		if downloadedOVAFilename == "" {
			downloadedOVAFilename = file.Name
		}
		cmd.Printf("Downloading file to %s...\n", downloadedOVAFilename)
		return Marketplace.Download(product.ProductId, downloadedOVAFilename, &pkg.DownloadRequestPayload{
			DeploymentFileId: file.FileID,
			AppVersion:       version.Number,
			EulaAccepted:     true,
		}, cmd.ErrOrStderr())
	},
}

var CreateOVACmd = &cobra.Command{
	Use:     "create",
	Short:   "Upload an attach an OVA",
	Long:    "Upload an attach an OVA to a product",
	Args:    cobra.NoArgs,
	PreRunE: GetUploadCredentials,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		uploader := internal.NewS3Uploader(Marketplace.StorageRegion, internal.HashAlgoSHA1, product.PublisherDetails.OrgId, UploadCredentials)
		file, err := uploader.Upload(Marketplace.StorageBucket, ovaFile)
		if err != nil {
			return err
		}

		file.AppVersion = version.Number
		product.PrepForUpdate()
		product.AddFile(file)

		updatedProduct, err := Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		return Output.RenderOVAs(updatedProduct.GetFilesForVersion(version.Number))
	},
}
