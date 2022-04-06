// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

const (
	ImageTagTypeFixed    = "FIXED"
	ImageTagTypeFloating = "FLOATING"
)

var (
	ContainerImageProductSlug            string
	ContainerImageProductVersion         string
	ContainerImageDeploymentInstructions string
	ImageRepository                      string
	ImageTag                             string
	ImageTagType                         string
)

func init() {
	rootCmd.AddCommand(ContainerImageCmd)
	ContainerImageCmd.AddCommand(ListContainerImageCmd)
	ContainerImageCmd.AddCommand(AttachContainerImageCmd)

	ListContainerImageCmd.Flags().StringVarP(&ContainerImageProductSlug, "product", "p", "", "Product slug (required)")
	_ = ListContainerImageCmd.MarkFlagRequired("product")
	ListContainerImageCmd.Flags().StringVarP(&ContainerImageProductVersion, "product-version", "v", "", "Product version (default to latest version)")

	AttachContainerImageCmd.Flags().StringVarP(&ContainerImageProductSlug, "product", "p", "", "Product slug (required)")
	_ = AttachContainerImageCmd.MarkFlagRequired("product")
	AttachContainerImageCmd.Flags().StringVarP(&ContainerImageProductVersion, "product-version", "v", "", "Product version (default to latest version)")
	AttachContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")
	_ = AttachContainerImageCmd.MarkFlagRequired("image-repository")
	AttachContainerImageCmd.Flags().StringVar(&ImageTag, "tag", "", "container repository tag")
	_ = AttachContainerImageCmd.MarkFlagRequired("tag")
	AttachContainerImageCmd.Flags().StringVar(&ImageTagType, "tag-type", "", "container repository tag type (fixed or floating)")
	_ = AttachContainerImageCmd.MarkFlagRequired("tag-type")
	AttachContainerImageCmd.Flags().StringVarP(&ContainerImageDeploymentInstructions, "deployment-instructions", "i", "", "deployment instructions")
	_ = AttachContainerImageCmd.MarkFlagRequired("deployment-instructions")
}

var ContainerImageCmd = &cobra.Command{
	Use:       "container-image",
	Aliases:   []string{"container-images"},
	Short:     "List and manage container images attached to a product",
	Long:      "List and manage container images attached to a product in the VMware Marketplace",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListContainerImageCmd.Use, AttachContainerImageCmd.Use},
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

var AttachContainerImageCmd = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a container image",
	Long:    "Attaches a container image to a product in the VMware Marketplace",
	Args:    cobra.NoArgs,
	PreRunE: GetRefreshToken,
	RunE: func(cmd *cobra.Command, args []string) error {
		ImageTagType = strings.ToUpper(ImageTagType)
		if ImageTagType != ImageTagTypeFixed && ImageTagType != ImageTagTypeFloating {
			return fmt.Errorf("invalid image tag type: %s. must be either \"%s\" or \"%s\"", ImageTagType, ImageTagTypeFixed, ImageTagTypeFloating)
		}

		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ContainerImageProductSlug, ContainerImageProductVersion)
		if err != nil {
			return err
		}

		product.SetDeploymentType(models.DeploymentTypesDocker)

		if product.HasContainerImage(version.Number, ImageRepository, ImageTag) {
			return fmt.Errorf("%s %s already has the tag %s:%s", ContainerImageProductSlug, version.Number, ImageRepository, ImageTag)
		}

		product.PrepForUpdate()
		product.DockerLinkVersions = append(product.DockerLinkVersions, &models.DockerVersionList{
			AppVersion: version.Number,
			DockerURLs: []*models.DockerURLDetails{
				{
					Url: ImageRepository,
					ImageTags: []*models.DockerImageTag{
						{
							Tag:  ImageTag,
							Type: ImageTagType,
						},
					},
					DeploymentInstruction: ContainerImageDeploymentInstructions,
					DockerType:            models.DockerTypeRegistry,
				},
			},
		})

		updatedProduct, err := Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		Output.PrintHeader(fmt.Sprintf("Container images for %s %s:", updatedProduct.DisplayName, version.Number))
		return Output.RenderContainerImages(updatedProduct.GetContainerImagesForVersion(version.Number))
	},
}
