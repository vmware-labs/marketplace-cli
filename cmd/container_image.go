// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ContainerImageProductSlug            string
	ContainerImageProductVersion         string
	ContainerImageDeploymentInstructions string
	ContainerImageCreateVersion          bool
	ImageRepository                      string
	ImageTag                             string
)

func init() {
	rootCmd.AddCommand(ContainerImageCmd)
	ContainerImageCmd.AddCommand(ListContainerImageCmd)
	ContainerImageCmd.AddCommand(OldAttachContainerImageCmd)

	ListContainerImageCmd.Flags().StringVarP(&ContainerImageProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListContainerImageCmd.MarkFlagRequired("product")
	ListContainerImageCmd.Flags().StringVarP(&ContainerImageProductVersion, "product-version", "v", "", "Product version (default to latest version)")

	OldAttachContainerImageCmd.Flags().StringVarP(&ContainerImageProductSlug, "product", "p", "", "Product slug (required)")
	_ = OldAttachContainerImageCmd.MarkFlagRequired("product")
	OldAttachContainerImageCmd.Flags().StringVarP(&ContainerImageProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	OldAttachContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")
	_ = OldAttachContainerImageCmd.MarkFlagRequired("image-repository")
	OldAttachContainerImageCmd.Flags().StringVar(&ImageTag, "tag", "", "container repository tag")
	_ = OldAttachContainerImageCmd.MarkFlagRequired("tag")
	OldAttachContainerImageCmd.Flags().StringVar(&AttachContainerImageTagType, "tag-type", "", "container repository tag type (fixed or floating)")
	_ = OldAttachContainerImageCmd.MarkFlagRequired("tag-type")
	OldAttachContainerImageCmd.Flags().StringVarP(&ContainerImageDeploymentInstructions, "deployment-instructions", "i", "", "deployment instructions")
	_ = OldAttachContainerImageCmd.MarkFlagRequired("deployment-instructions")
	OldAttachContainerImageCmd.Flags().BoolVar(&ContainerImageCreateVersion, "create-version", false, "create the product version, if it doesn't already exist")
}

var ContainerImageCmd = &cobra.Command{
	Use:       "container-image",
	Aliases:   []string{"container-images"},
	Short:     "List and manage container images attached to a product",
	Long:      "List and manage container images attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListContainerImageCmd.Use, OldAttachContainerImageCmd.Use},
}

var ListContainerImageCmd = &cobra.Command{
	Use:     "list",
	Short:   "List product container images",
	Long:    "Prints the list of container images attached to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ContainerImageProductSlug, ContainerImageProductVersion)
		if err != nil {
			return err
		}

		images := product.GetContainerImagesForVersion(version.Number)

		Output.PrintHeader(fmt.Sprintf("Container images for %s %s:", product.DisplayName, version.Number))
		return Output.RenderContainerImages(images)
	},
}

var OldAttachContainerImageCmd = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a container image",
	Long:    "Attaches a container image to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: RunSerially(GetRefreshToken, ValidateTagType),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.PrintErrln("mkpcli container-image attach has been deprecated and will be removed in the next major version. Please use mkpcli attach image instead.")
		AttachProductSlug = ContainerImageProductSlug
		AttachProductVersion = ContainerImageProductVersion
		AttachCreateVersion = ContainerImageCreateVersion
		AttachContainerImage = ImageRepository
		AttachContainerImageTag = ImageTag
		return AttachContainerImageCmd.RunE(cmd, args)
	},
}
