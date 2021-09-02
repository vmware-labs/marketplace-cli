// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

const (
	ImageTagTypeFixed    = "FIXED"
	ImageTagTypeFloating = "FLOATING"
)

var downloadedContainerImageFilename string

func init() {
	rootCmd.AddCommand(ContainerImageCmd)
	ContainerImageCmd.AddCommand(ListContainerImageCmd)
	ContainerImageCmd.AddCommand(GetContainerImageCmd)
	ContainerImageCmd.AddCommand(DownloadContainerImageCmd)
	ContainerImageCmd.AddCommand(CreateContainerImageCmd)

	ContainerImageCmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = ContainerImageCmd.MarkPersistentFlagRequired("product")
	ContainerImageCmd.PersistentFlags().StringVarP(&ProductVersion, "product-version", "v", "latest", "Product version")
	_ = ContainerImageCmd.MarkPersistentFlagRequired("product-version")

	GetContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")

	DownloadContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")
	DownloadContainerImageCmd.Flags().StringVarP(&downloadedContainerImageFilename, "filename", "f", "image.tar", "output file name")

	CreateContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")
	_ = CreateContainerImageCmd.MarkFlagRequired("image-repository")
	CreateContainerImageCmd.Flags().StringVarP(&ImageTag, "image-tag", "t", "", "container repository tag")
	_ = CreateContainerImageCmd.MarkFlagRequired("image-tag")
	CreateContainerImageCmd.Flags().StringVarP(&ImageTagType, "image-tag-type", "y", "", "container repository tag type (fixed or floating)")
	_ = CreateContainerImageCmd.MarkFlagRequired("image-tag-type")
	CreateContainerImageCmd.Flags().StringVarP(&DeploymentInstructions, "deployment-instructions", "i", "", "deployment instructions")
}

var ContainerImageCmd = &cobra.Command{
	Use:       "container-image",
	Aliases:   []string{"container-images"},
	Short:     "Container image related commands",
	Long:      "Interact with container images attached to a Marketplace product",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"list", "get", "download", "create"},
}

var ListContainerImageCmd = &cobra.Command{
	Use:   "list",
	Short: "List container images",
	Long:  "List the container images attached to a product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		return Output.RenderContainerImages(product, ProductVersion)
	},
}

var GetContainerImageCmd = &cobra.Command{
	Use:   "get",
	Short: "Get container image details",
	Long:  "Get details for a container image attached to a product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		containerImages := product.GetContainerImagesForVersion(ProductVersion)
		if containerImages == nil {
			return fmt.Errorf("%s %s does not have any container images\n", product.Slug, version)
		}

		var containerImage *models.DockerURLDetails
		if ImageRepository != "" {
			containerImage = containerImages.GetImage(ImageRepository)
			if containerImage == nil {
				return fmt.Errorf("%s %s does not have the container image \"%s\"", ProductSlug, ProductVersion, ImageRepository)
			}
		} else {
			if len(containerImages.DockerURLs) == 0 {
				return fmt.Errorf("%s %s does not have any container images\n", product.Slug, version)
			} else if len(containerImages.DockerURLs) == 1 {
				containerImage = containerImages.DockerURLs[0]
			} else {
				return fmt.Errorf("multiple container images found for %s %s, please use the --image-repository parameter", ProductSlug, ProductVersion)
			}
		}

		return Output.RenderContainerImage(product, ProductVersion, containerImage)
	},
}

var DownloadContainerImageCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a container image",
	Long:  "Download a container image attached to a product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		containerImages := product.GetContainerImagesForVersion(ProductVersion)
		if containerImages == nil {
			return fmt.Errorf("%s %s does not have any container images\n", product.Slug, version)
		}

		var containerImage *models.DockerURLDetails
		if ImageRepository != "" {
			containerImage = containerImages.GetImage(ImageRepository)
			if containerImage == nil {
				return fmt.Errorf("%s %s does not have the container image \"%s\"", ProductSlug, ProductVersion, ImageRepository)
			}
		} else {
			if len(containerImages.DockerURLs) == 0 {
				return fmt.Errorf("no container images found for %s %s", ProductSlug, ProductVersion)
			} else if len(containerImages.DockerURLs) == 1 {
				containerImage = containerImages.DockerURLs[0]
			} else {
				return fmt.Errorf("multiple container images found for %s %s, please use the --image-repository parameter", ProductSlug, ProductVersion)
			}
		}

		cmd.Printf("Downloading image to %s...\n", downloadedContainerImageFilename)
		return Marketplace.Download(product.ProductId, downloadedContainerImageFilename, &pkg.DownloadRequestPayload{
			DockerlinkVersionID: containerImages.ID,
			DockerUrlId:         containerImage.ID,
			ImageTagId:          containerImage.ImageTags[0].ID,
			AppVersion:          containerImages.AppVersion,
			EulaAccepted:        true,
		}, cmd.ErrOrStderr())
	},
}

var CreateContainerImageCmd = &cobra.Command{
	Use:   "create",
	Short: "Attach a container image",
	Long:  "Attach a container image to a product",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ImageTagType = strings.ToUpper(ImageTagType)
		if ImageTagType != ImageTagTypeFixed && ImageTagType != ImageTagTypeFloating {
			return fmt.Errorf("invalid image tag type: %s. must be either \"%s\" or \"%s\"", ImageTagType, ImageTagTypeFixed, ImageTagTypeFloating)
		}

		cmd.SilenceUsage = true
		product, err := Marketplace.GetProductWithVersion(ProductSlug, ProductVersion)
		if err != nil {
			return err
		}

		product.SetDeploymentType(models.DeploymentTypesDocker)

		containerImages := product.GetContainerImagesForVersion(ProductVersion)
		if containerImages == nil {
			if DeploymentInstructions == "" {
				cmd.SilenceUsage = false
				return fmt.Errorf("must specify the deployment instructions for the first container image. Please run again with --deployment-instructions <string>")
			}

			containerImages = &models.DockerVersionList{
				AppVersion: ProductVersion,
				DockerURLs: []*models.DockerURLDetails{},
				//DeploymentInstruction: DeploymentInstructions,
			}
			product.DockerLinkVersions = append(product.DockerLinkVersions, containerImages)
		}

		containerImage := containerImages.GetImage(ImageRepository)
		if containerImage == nil {
			containerImage = &models.DockerURLDetails{
				Url:                   ImageRepository,
				ImageTags:             []*models.DockerImageTag{},
				DeploymentInstruction: DeploymentInstructions,
			}
			containerImages.DockerURLs = append(containerImages.DockerURLs, containerImage)
		}

		if containerImage.HasTag(ImageTag) {
			return fmt.Errorf("%s %s already has the container image %s:%s", ProductSlug, ProductVersion, ImageRepository, ImageTag)
		}
		containerImage.ImageTags = append(containerImage.ImageTags, &models.DockerImageTag{
			Tag:  ImageTag,
			Type: ImageTagType,
		})

		updatedProduct, err := Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		return Output.RenderContainerImages(updatedProduct, ProductVersion)
	},
}
