// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	. "github.com/vmware-labs/marketplace-cli/v2/lib"
	"github.com/vmware-labs/marketplace-cli/v2/models"
)

const (
	ImageTagTypeFixed    = "FIXED"
	ImageTagTypeFloating = "FLOATING"
)

func init() {

	rootCmd.AddCommand(ContainerImageCmd)
	ContainerImageCmd.AddCommand(ListContainerImageCmd)
	ContainerImageCmd.AddCommand(GetContainerImageCmd)
	ContainerImageCmd.AddCommand(CreateContainerImageCmd)
	ContainerImageCmd.PersistentFlags().StringVarP(&OutputFormat, "output-format", "f", FormatTable, "Output OutputFormat")

	ContainerImageCmd.PersistentFlags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = ContainerImageCmd.MarkPersistentFlagRequired("product")
	ContainerImageCmd.PersistentFlags().StringVarP(&ProductVersion, "product-version", "v", "", "Product version")
	_ = ContainerImageCmd.MarkPersistentFlagRequired("product-version")

	GetContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")
	_ = GetContainerImageCmd.MarkFlagRequired("image-repository")

	CreateContainerImageCmd.Flags().StringVarP(&ImageRepository, "image-repository", "r", "", "container repository")
	_ = CreateContainerImageCmd.MarkFlagRequired("image-repository")
	CreateContainerImageCmd.Flags().StringVarP(&ImageTag, "image-tag", "t", "", "container repository tag")
	_ = CreateContainerImageCmd.MarkFlagRequired("image-tag")
	CreateContainerImageCmd.Flags().StringVarP(&ImageTagType, "image-tag-type", "y", "", "container repository tag type (fixed or floating)")
	_ = CreateContainerImageCmd.MarkFlagRequired("image-tag-type")
	CreateContainerImageCmd.Flags().StringVarP(&DeploymentInstructions, "deployment-instructions", "i", "", "deployment instructions")
}

var ContainerImageCmd = &cobra.Command{
	Use:               "container-image",
	Aliases:           []string{"container-images"},
	Short:             "container images",
	Long:              "",
	Args:              cobra.OnlyValidArgs,
	ValidArgs:         []string{"get", "list", "create"},
	PersistentPreRunE: GetRefreshToken,
}

var ListContainerImageCmd = &cobra.Command{
	Use:   "list",
	Short: "list container images",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		product := response.Response.Data
		if !product.HasVersion(ProductVersion) {
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" does not have a version %s", ProductSlug, ProductVersion)
		}
		containerImages := product.GetDockerImagesForVersion(ProductVersion)
		if containerImages == nil {
			cmd.Printf("product \"%s\" %s does not have any container images", ProductSlug, ProductVersion)
			return nil
		}

		err = RenderContainerImages(OutputFormat, containerImages, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return errors.Wrapf(err, "failed to render the container images")
		}

		return nil
	},
}

var GetContainerImageCmd = &cobra.Command{
	Use:   "get",
	Short: "get a container image",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		product := response.Response.Data
		if !product.HasVersion(ProductVersion) {
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" does not have a version %s", ProductSlug, ProductVersion)
		}

		containerImages := product.GetDockerImagesForVersion(ProductVersion)
		if containerImages == nil {
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" does not have any container images for version %s", ProductSlug, ProductVersion)
		}

		containerImage := containerImages.GetImage(ImageRepository)
		if containerImage == nil {
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" %s does not have the container image \"%s\"", ProductSlug, ProductVersion, ImageRepository)
		}

		err = RenderContainerImage(OutputFormat, containerImage, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return errors.Wrap(err, "failed to render the container images")
		}

		return nil
	},
}

var CreateContainerImageCmd = &cobra.Command{
	Use:   "create",
	Short: "add a container image to a product version",
	Long:  "",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ImageTagType = strings.ToUpper(ImageTagType)
		if ImageTagType != ImageTagTypeFixed && ImageTagType != ImageTagTypeFloating {
			return errors.Errorf("invalid image tag type: %s. must be either \"%s\" or \"%s\"", ImageTagType, ImageTagTypeFixed, ImageTagTypeFloating)
		}

		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		product := response.Response.Data

		if !product.HasVersion(ProductVersion) {
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" does not have a version %s, please add it first", ProductSlug, ProductVersion)
		}

		product.SetDeploymentType(models.DeploymentTypesDocker)

		containerImages := product.GetDockerImagesForVersion(ProductVersion)
		if containerImages == nil {
			if DeploymentInstructions == "" {
				return errors.Errorf("must specify the deployment instructions for the first container image. Please run again with --deployment-instructions <string>")
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
			cmd.SilenceUsage = true
			return errors.Errorf("product \"%s\" %s already has the container image %s:%s", ProductSlug, ProductVersion, ImageRepository, ImageTag)
		}
		containerImage.ImageTags = append(containerImage.ImageTags, &models.DockerImageTag{
			Tag:  ImageTag,
			Type: ImageTagType,
		})

		response = &GetProductResponse{}
		err = PutProduct(product, false, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		containerImages = response.Response.Data.GetDockerImagesForVersion(ProductVersion)
		err = RenderContainerImages(OutputFormat, containerImages, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return errors.Wrap(err, "failed to render the container images")
		}
		return nil
	},
}
