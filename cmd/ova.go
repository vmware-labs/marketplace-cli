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
	OVAProductSlug        string
	OVAProductVersion     string
	downloadedOVAFilename string
)

func init() {
	rootCmd.AddCommand(OVACmd)
	OVACmd.AddCommand(ListOVACmd)
	OVACmd.AddCommand(GetOVACmd)
	OVACmd.AddCommand(DownloadOVACmd)
	OVACmd.AddCommand(AttachOVACmd)

	ListOVACmd.Flags().StringVarP(&OVAProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListOVACmd.MarkFlagRequired("product")
	ListOVACmd.Flags().StringVarP(&OVAProductVersion, "product-version", "v", "", "Product version (default to latest version)")

	GetOVACmd.Flags().StringVarP(&OVAProductSlug, "product", "p", "", "Product slug (required)")
	_ = GetOVACmd.MarkFlagRequired("product")
	GetOVACmd.Flags().StringVarP(&OVAProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	GetOVACmd.Flags().StringVar(&ovaFile, "file-id", "", "The file ID of the file to get")

	DownloadOVACmd.Flags().StringVarP(&OVAProductSlug, "product", "p", "", "Product slug (required)")
	_ = DownloadOVACmd.MarkFlagRequired("product")
	DownloadOVACmd.Flags().StringVarP(&OVAProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	DownloadOVACmd.Flags().StringVar(&ovaFile, "file-id", "", "The file ID of the file to download")
	DownloadOVACmd.Flags().StringVarP(&downloadedOVAFilename, "filename", "f", "", "Downloaded file name (default is original file name)")

	AttachOVACmd.Flags().StringVarP(&OVAProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachOVACmd.MarkFlagRequired("product")
	AttachOVACmd.Flags().StringVarP(&OVAProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachOVACmd.Flags().StringVar(&ovaFile, "ova-file", "", "OVA file to upload")
}

var OVACmd = &cobra.Command{
	Use:       "ova",
	Aliases:   []string{"ovas"},
	Short:     "List and manage OVAs attached to a product",
	Long:      "List and manage OVAs attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListOVACmd.Use, GetOVACmd.Use, DownloadOVACmd.Use, AttachOVACmd.Use},
}

var ListOVACmd = &cobra.Command{
	Use:   "list",
	Short: "List product OVAs",
	Long:  "Prints the list of OVAs attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(OVAProductSlug, OVAProductVersion)
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
	Short: "Get details for an OVA",
	Long:  "Prints detailed information about an OVA file attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(OVAProductSlug, OVAProductVersion)
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
				return fmt.Errorf("no files found for %s for version %s", OVAProductSlug, version.Number)
			} else if len(files) == 1 {
				file = files[0]
			} else {
				return fmt.Errorf("multiple files found for %s for version %s, please use the --file-id parameter", OVAProductSlug, version.Number)
			}
		}

		return Output.RenderOVA(file)
	},
}

var DownloadOVACmd = &cobra.Command{
	Use:   "download",
	Short: "Download an OVA",
	Long:  "Downloads an OVA attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(OVAProductSlug, OVAProductVersion)
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
				return fmt.Errorf("no files found for %s for version %s", OVAProductSlug, version.Number)
			} else if len(files) == 1 {
				file = files[0]
			} else {
				return fmt.Errorf("multiple files found for %s for version %s, please use the --file-id parameter", OVAProductSlug, version.Number)
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

var AttachOVACmd = &cobra.Command{
	Use:     "attach",
	Short:   "Upload and attach an OVA",
	Long:    "Uploads and attaches an OVA to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetUploadCredentials,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(OVAProductSlug, OVAProductVersion)
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
