// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

var (
	vmFile           string
	VMProductSlug    string
	VMProductVersion string
	VMCreateVersion  bool
)

func init() {
	rootCmd.AddCommand(VMCmd)
	VMCmd.AddCommand(ListVMCmd)
	VMCmd.AddCommand(GetVMCmd)
	VMCmd.AddCommand(OldAttachVMCmd)

	ListVMCmd.Flags().StringVarP(&VMProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListVMCmd.MarkFlagRequired("product")
	ListVMCmd.Flags().StringVarP(&VMProductVersion, "product-version", "v", "", "Product version (default to latest version)")

	GetVMCmd.Flags().StringVarP(&VMProductSlug, "product", "p", "", "Product slug (required)")
	_ = GetVMCmd.MarkFlagRequired("product")
	GetVMCmd.Flags().StringVarP(&VMProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	GetVMCmd.Flags().StringVar(&vmFile, "file-id", "", "The file ID of the file to get")

	OldAttachVMCmd.Flags().StringVarP(&VMProductSlug, "product", "p", "", "Product slug (required)")
	_ = OldAttachVMCmd.MarkFlagRequired("product")
	OldAttachVMCmd.Flags().StringVarP(&VMProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	OldAttachVMCmd.Flags().StringVar(&vmFile, "file", "", "Virtual machine file to upload (required)")
	_ = OldAttachVMCmd.MarkFlagRequired("file")
	OldAttachVMCmd.Flags().BoolVar(&VMCreateVersion, "create-version", false, "create the product version, if it doesn't already exist")
}

var VMCmd = &cobra.Command{
	Use:       "vm",
	Aliases:   []string{"vms"},
	Short:     "List and manage virtual machines attached to a product",
	Long:      "List and manage virtual machine files attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListVMCmd.Use, GetVMCmd.Use, AttachVMCmd.Use},
}

var ListVMCmd = &cobra.Command{
	Use:     "list",
	Short:   "List product virtual machines",
	Long:    "Prints the list of virtual machine files attached to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(VMProductSlug, VMProductVersion)
		if err != nil {
			return err
		}

		files := product.GetFilesForVersion(version.Number)

		Output.PrintHeader(fmt.Sprintf("Virtual machine files for %s %s:", product.DisplayName, version.Number))
		return Output.RenderFiles(files)
	},
}

var GetVMCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get details for a virtual machine",
	Long:    "Prints detailed information about a virtual machine file attached to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		product, version, err := Marketplace.GetProductWithVersion(VMProductSlug, VMProductVersion)
		if err != nil {
			return err
		}

		var file *models.ProductDeploymentFile
		if vmFile != "" {
			file = product.GetFile(vmFile)
			if file == nil {
				return fmt.Errorf("no file found with ID %s", vmFile)
			}
		} else {
			files := product.GetFilesForVersion(version.Number)
			if len(files) == 0 {
				return fmt.Errorf("no files found for %s for version %s", VMProductSlug, version.Number)
			} else if len(files) == 1 {
				file = files[0]
			} else {
				return fmt.Errorf("multiple files found for %s for version %s, please use the --file-id parameter", VMProductSlug, version.Number)
			}
		}

		return Output.RenderFile(file)
	},
}

var OldAttachVMCmd = &cobra.Command{
	Use:     "attach",
	Short:   "Upload and attach a virtual machine file (ISO or OVA)",
	Long:    "Uploads and attaches a virtual machine file (ISO or OVA) to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.PrintErrln("mkpcli vm attach has been deprecated and will be removed in the next major version. Please use mkpcli attach vm instead.")
		AttachProductSlug = VMProductSlug
		AttachProductVersion = VMProductVersion
		AttachCreateVersion = VMCreateVersion
		AttachVMFile = vmFile
		return AttachVMCmd.RunE(cmd, args)
	},
}
