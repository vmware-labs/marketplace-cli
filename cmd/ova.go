// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	. "github.com/vmware-labs/marketplace-cli/v2/lib"
	"github.com/vmware-labs/marketplace-cli/v2/models"
)

func init() {
	rootCmd.AddCommand(OVACmd)
	OVACmd.AddCommand(ListOVACmd)
	OVACmd.AddCommand(CreateOVACmd)
	OVACmd.PersistentFlags().StringVarP(&OutputFormat, "output-format", "f", FormatTable, "Output format")

	OVACmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = OVACmd.MarkPersistentFlagRequired("product")
	OVACmd.PersistentFlags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	_ = OVACmd.MarkPersistentFlagRequired("product-version")

	CreateOVACmd.Flags().StringVar(&OVAFile, "ova-file", "", "OVA file to upload")
}

var OVACmd = &cobra.Command{
	Use:               "ova",
	Aliases:           []string{"ovas"},
	Short:             "ova",
	Long:              "",
	Args:              cobra.OnlyValidArgs,
	ValidArgs:         []string{"get", "list", "create"},
	PersistentPreRunE: GetRefreshToken,
}

var ListOVACmd = &cobra.Command{
	Use:   "list",
	Short: "list OVAs",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			return err
		}

		product := response.Response.Data
		if !product.HasVersion(ProductVersion) {
			return fmt.Errorf("product \"%s\" does not have a version %s", ProductSlug, ProductVersion)
		}

		return RenderOVAs(OutputFormat, ProductVersion, product, cmd.OutOrStdout())
	},
}

var CreateOVACmd = &cobra.Command{
	Use:     "create",
	Short:   "add an OVA to a product",
	Long:    "",
	Args:    cobra.NoArgs,
	PreRunE: GetUploadCredentials,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		product := response.Response.Data

		if !product.HasVersion(ProductVersion) {
			cmd.SilenceUsage = true
			return fmt.Errorf("product \"%s\" does not have a version %s, please add it first", ProductSlug, ProductVersion)
		}

		product.SetDeploymentType(models.DeploymentTypeOVA)

		hashAlgo := internal.HashAlgoSHA1
		uploader := internal.NewS3Uploader("us-east-2", hashAlgo, product.PublisherDetails.OrgId, UploadCredentials)
		fileURL, fileHash, err := uploader.Upload(OVAFile)
		if err != nil {
			return err
		}

		product.ProductDeploymentFiles = []*models.ProductDeploymentFile{{
			Url:        fileURL,
			AppVersion: ProductVersion,
			HashDigest: fileHash,
			HashAlgo:   models.HashAlgoSHA1,
		}}

		var putResponse GetProductResponse
		err = PutProduct(product, false, &putResponse)
		if err != nil {
			return err
		}

		return RenderOVAs(OutputFormat, ProductVersion, putResponse.Response.Data, cmd.OutOrStdout())
	},
}
