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

var (
	ContainerImageProductSlug            string
	ContainerImageProductVersion         string
	ContainerImageDeploymentInstructions string
	downloadedContainerImageFilename     string
	ImageRepository                      string
	ImageTag                             string
	ImageTagType                         string
)

func init() {
	rootCmd.AddCommand(ContainerImageCmd)
	ContainerImageCmd.AddCommand(ListContainerImageCmd)
	ContainerImageCmd.AddCommand(GetContainerImageCmd)
	ContainerImageCmd.AddCommand(DownloadContainerImageCmd)
	ContainerImageCmd.AddCommand(AttachContainerImageCmd)

	ContainerImageCmd.PersistentFlags().StringVarP(&ContainerImageProductSlug, "product", "p", "", "Product slug")
	_ = ContainerImageCmd.MarkPersistentFlagRequired("product")
	ContainerImageCmd.PersistentFlags().StringVarP(&ContainerImageProductVersion, "product-version", "v", "latest", "Product version")

	GetContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")

	DownloadContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")
	DownloadContainerImageCmd.Flags().StringVar(&ImageTag, "tag", "", "container image tag")
	DownloadContainerImageCmd.Flags().StringVarP(&downloadedContainerImageFilename, "filename", "f", "image.tar", "output file name")

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
	ValidArgs: []string{ListContainerImageCmd.Use, GetContainerImageCmd.Use, DownloadContainerImageCmd.Use, AttachContainerImageCmd.Use},
}

var ListContainerImageCmd = &cobra.Command{
	Use:   "list",
	Short: "List product container images",
	Long:  "Prints the list of container images attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ContainerImageProductSlug, ContainerImageProductVersion)
		if err != nil {
			return err
		}

		images := product.GetContainerImagesForVersion(version.Number)
		if images == nil || len(images.DockerURLs) == 0 {
			cmd.Printf("%s %s does not have any container images\n", product.Slug, version.Number)
			return nil
		}

		return Output.RenderContainerImages(images)
	},
}

var GetContainerImageCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details for a container image",
	Long:  "Prints detailed information about a container image attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ContainerImageProductSlug, ContainerImageProductVersion)
		if err != nil {
			return err
		}

		containerImages := product.GetContainerImagesForVersion(version.Number)
		if containerImages == nil {
			return fmt.Errorf("%s %s does not have any container images\n", product.Slug, version.Number)
		}

		var containerImage *models.DockerURLDetails
		if ImageRepository != "" {
			containerImage = containerImages.GetImage(ImageRepository)
			if containerImage == nil {
				return fmt.Errorf("%s %s does not have the container image \"%s\"", ContainerImageProductSlug, version.Number, ImageRepository)
			}
		} else {
			if len(containerImages.DockerURLs) == 0 {
				return fmt.Errorf("%s %s does not have any container images\n", product.Slug, version.Number)
			} else if len(containerImages.DockerURLs) == 1 {
				containerImage = containerImages.DockerURLs[0]
			} else {
				return fmt.Errorf("multiple container images found for %s %s, please use the --image-repository parameter", ContainerImageProductSlug, version.Number)
			}
		}

		return Output.RenderContainerImage(containerImage)
	},
}

var DownloadContainerImageCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a container image",
	Long:  "Downloads a container image attached to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		product, version, err := Marketplace.GetProductWithVersion(ContainerImageProductSlug, ContainerImageProductVersion)
		if err != nil {
			return err
		}

		containerImages := product.GetContainerImagesForVersion(version.Number)
		if containerImages == nil {
			return fmt.Errorf("%s %s does not have any container images\n", product.Slug, version.Number)
		}

		var containerImage *models.DockerURLDetails
		if ImageRepository != "" {
			containerImage = containerImages.GetImage(ImageRepository)
			if containerImage == nil {
				return fmt.Errorf("%s %s does not have the container image \"%s\"", ContainerImageProductSlug, version.Number, ImageRepository)
			}
		} else {
			if len(containerImages.DockerURLs) == 0 {
				return fmt.Errorf("no container images found for %s %s", ContainerImageProductSlug, version.Number)
			} else if len(containerImages.DockerURLs) == 1 {
				containerImage = containerImages.DockerURLs[0]
			} else {
				return fmt.Errorf("multiple container images found for %s %s, please use the --image-repository parameter", ContainerImageProductSlug, version.Number)
			}
		}

		var imageTag *models.DockerImageTag
		if ImageTag != "" {
			imageTag = containerImage.GetTag(ImageTag)
			if imageTag == nil {
				return fmt.Errorf("%s %s does not have the an image tag \"%s\" for %s", ContainerImageProductSlug, version.Number, ImageTag, containerImage.Url)
			}
		} else {
			if len(containerImage.ImageTags) == 0 {
				return fmt.Errorf("no tags images found for %s in %s %s", containerImage.Url, ContainerImageProductSlug, version.Number)
			} else if len(containerImages.DockerURLs) == 1 {
				imageTag = containerImage.ImageTags[0]
			} else {
				return fmt.Errorf("multiple tags found for %s in %s %s, please use the --tag parameter", containerImage.Url, ContainerImageProductSlug, version.Number)
			}
		}

		if !imageTag.IsUpdatedInMarketplaceRegistry {
			return fmt.Errorf("%s with tag %s in %s %s is not downloadable: %s", containerImage.Url, imageTag.Tag, ContainerImageProductSlug, version.Number, imageTag.ProcessingError)
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

var AttachContainerImageCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a container image",
	Long:  "Attaches a container image to a product in the VMware Marketplace",
	Args:  cobra.NoArgs,
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

		containerImages := product.GetContainerImagesForVersion(version.Number)
		if containerImages == nil {
			if ContainerImageDeploymentInstructions == "" {
				cmd.SilenceUsage = false
				return fmt.Errorf("must specify the deployment instructions for the first container image. Please run again with --deployment-instructions <string>")
			}

			containerImages = &models.DockerVersionList{
				AppVersion: version.Number,
				DockerURLs: []*models.DockerURLDetails{},
			}
		}

		containerImage := containerImages.GetImage(ImageRepository)
		if containerImage == nil {
			containerImage = &models.DockerURLDetails{
				Url:                   ImageRepository,
				ImageTags:             []*models.DockerImageTag{},
				DeploymentInstruction: ContainerImageDeploymentInstructions,
			}
			containerImages.DockerURLs = append(containerImages.DockerURLs, containerImage)
		}

		if containerImage.HasTag(ImageTag) {
			return fmt.Errorf("%s %s already has the tag %s:%s", ContainerImageProductSlug, version.Number, ImageRepository, ImageTag)
		}
		containerImage.ImageTags = append(containerImage.ImageTags, &models.DockerImageTag{
			Tag:  ImageTag,
			Type: ImageTagType,
		})

		product.PrepForUpdate()
		product.DockerLinkVersions = append(product.DockerLinkVersions, containerImages)

		updatedProduct, err := Marketplace.PutProduct(product, false)
		if err != nil {
			return err
		}

		return Output.RenderContainerImages(updatedProduct.GetContainerImagesForVersion(version.Number))
	},
}
