// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output/outputfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("ContainerImage", func() {
	var (
		marketplace *pkgfakes.FakeMarketplaceInterface
		output      *outputfakes.FakeFormat
		product     *models.Product
	)

	BeforeEach(func() {
		marketplace = &pkgfakes.FakeMarketplaceInterface{}
		cmd.Marketplace = marketplace

		marketplace.GetProductWithVersionStub = func(slug string, version string) (*models.Product, *models.Version, error) {
			Expect(slug).To(Equal("my-super-product"))
			Expect(version).To(Equal("1.2.3"))
			return product, &models.Version{Number: "1.2.3"}, nil
		}

		output = &outputfakes.FakeFormat{}
		cmd.Output = output
	})

	Describe("ListContainerImageCmd", func() {
		BeforeEach(func() {
			container := test.CreateFakeContainerImage(
				"myId",
				"0.0.1",
				"latest",
			)

			product = test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")
			test.AddContainerImages(product, "1.2.3", "Machine wash cold with like colors", container)
		})

		It("outputs the container images", func() {
			cmd.ContainerImageProductSlug = "my-super-product"
			cmd.ContainerImageProductVersion = "1.2.3"
			err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderContainerImagesCallCount()).To(Equal(1))
				images := output.RenderContainerImagesArgsForCall(0)
				Expect(images.AppVersion).To(Equal("1.2.3"))
				Expect(images.DockerURLs).To(HaveLen(1))
				Expect(images.DockerURLs[0].ImageTags).To(HaveLen(2))
				Expect(images.DockerURLs[0].ImageTags[0].Tag).To(Equal("0.0.1"))
				Expect(images.DockerURLs[0].ImageTags[0].Type).To(Equal("FIXED"))
				Expect(images.DockerURLs[0].ImageTags[1].Tag).To(Equal("latest"))
				Expect(images.DockerURLs[0].ImageTags[1].Type).To(Equal("FLOATING"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})

	Describe("AttachContainerImageCmd", func() {
		var productID string
		BeforeEach(func() {
			nginx := test.CreateFakeContainerImage("nginx", "latest")
			python := test.CreateFakeContainerImage("python", "1.2.3")

			productID = uuid.New().String()
			product = test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3")
			test.AddContainerImages(product, "1.2.3", "Machine wash cold with like colors", nginx)

			updatedProduct := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(updatedProduct, "1.2.3")
			test.AddContainerImages(updatedProduct, "1.2.3", "Machine wash cold with like colors", nginx, python)
			marketplace.PutProductReturns(updatedProduct, nil)
		})

		It("outputs the new container image", func() {
			cmd.ContainerImageProductSlug = "my-super-product"
			cmd.ContainerImageProductVersion = "1.2.3"
			cmd.ImageRepository = "python"
			cmd.ImageTag = "1.2.3"
			cmd.ImageTagType = cmd.ImageTagTypeFixed
			err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
			})

			By("updating the product with the new container image", func() {
				Expect(marketplace.PutProductCallCount()).To(Equal(1))
				updatedProduct, versionUpdate := marketplace.PutProductArgsForCall(0)
				Expect(versionUpdate).To(BeFalse())

				Expect(updatedProduct.DeploymentTypes).To(ContainElement("DOCKERLINK"))
				Expect(updatedProduct.DockerLinkVersions).To(HaveLen(1))
				dockerLink := updatedProduct.DockerLinkVersions[0]
				Expect(dockerLink.AppVersion).To(Equal("1.2.3"))
				Expect(dockerLink.DockerURLs).To(HaveLen(2))
				dockerUrl := dockerLink.DockerURLs[0]
				Expect(dockerUrl.Url).To(Equal("nginx"))
				Expect(dockerUrl.ImageTags).To(HaveLen(1))
				tag := dockerUrl.ImageTags[0]
				Expect(tag.Tag).To(Equal("latest"))
				Expect(tag.Type).To(Equal("FLOATING"))

				dockerUrl = dockerLink.DockerURLs[1]
				Expect(dockerUrl.Url).To(Equal("python"))
				Expect(dockerUrl.ImageTags).To(HaveLen(1))
				tag = dockerUrl.ImageTags[0]
				Expect(tag.Tag).To(Equal("1.2.3"))
				Expect(tag.Type).To(Equal("FIXED"))
			})

			By("outputting the response", func() {
				Expect(output.RenderContainerImagesCallCount()).To(Equal(1))
				images := output.RenderContainerImagesArgsForCall(0)
				Expect(images.DockerURLs).To(HaveLen(2))
			})
		})

		Context("Adding a new tag to an existing container image", func() {
			BeforeEach(func() {
				nginx := test.CreateFakeContainerImage("nginx", "latest")
				nginxUpdated := test.CreateFakeContainerImage("nginx", "latest", "5.5.5")

				productID = uuid.New().String()
				product = test.CreateFakeProduct(
					productID,
					"My Super Product",
					"my-super-product",
					"PENDING")
				test.AddVerions(product, "1.2.3")
				test.AddContainerImages(product, "1.2.3", "Machine wash cold with like colors", nginx)

				updatedProduct := test.CreateFakeProduct(
					productID,
					"My Super Product",
					"my-super-product",
					"PENDING")
				test.AddVerions(updatedProduct, "1.2.3")
				test.AddContainerImages(updatedProduct, "1.2.3", "Machine wash cold with like colors", nginxUpdated)
				marketplace.PutProductReturns(updatedProduct, nil)
			})

			It("outputs the new container image", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("getting the product details", func() {
					Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				})

				By("updating the product with the new container image", func() {
					Expect(marketplace.PutProductCallCount()).To(Equal(1))
					updatedProduct, versionUpdate := marketplace.PutProductArgsForCall(0)
					Expect(versionUpdate).To(BeFalse())

					Expect(updatedProduct.DeploymentTypes).To(ContainElement("DOCKERLINK"))
					Expect(updatedProduct.DockerLinkVersions).To(HaveLen(1))
					dockerLink := updatedProduct.DockerLinkVersions[0]
					Expect(dockerLink.AppVersion).To(Equal("1.2.3"))
					Expect(dockerLink.DockerURLs).To(HaveLen(1))
					dockerUrl := dockerLink.DockerURLs[0]
					Expect(dockerUrl.Url).To(Equal("nginx"))
					Expect(dockerUrl.ImageTags).To(HaveLen(2))
					tag := dockerUrl.ImageTags[0]
					Expect(tag.Tag).To(Equal("latest"))
					Expect(tag.Type).To(Equal("FLOATING"))
					tag = dockerUrl.ImageTags[1]
					Expect(tag.Tag).To(Equal("5.5.5"))
					Expect(tag.Type).To(Equal("FIXED"))
				})

				By("outputting the response", func() {
					Expect(output.RenderContainerImagesCallCount()).To(Equal(1))
					images := output.RenderContainerImagesArgsForCall(0)
					Expect(images.DockerURLs[0].ImageTags).To(HaveLen(2))
					Expect(images.DockerURLs[0].ImageTags[0].Tag).To(Equal("latest"))
					Expect(images.DockerURLs[0].ImageTags[0].Type).To(Equal("FLOATING"))
					Expect(images.DockerURLs[0].ImageTags[1].Tag).To(Equal("5.5.5"))
					Expect(images.DockerURLs[0].ImageTags[1].Type).To(Equal("FIXED"))
				})
			})
		})

		Context("Container image already exists", func() {
			It("says the image already exists", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "latest"
				cmd.ImageTagType = cmd.ImageTagTypeFloating
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("my-super-product 1.2.3 already has the tag nginx:latest"))
			})
		})

		Context("invalid tag type", func() {
			It("prints the error", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = "fancy"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid image tag type: FANCY. must be either \"FIXED\" or \"FLOATING\""))
			})
		})

		Context("Error putting product", func() {
			BeforeEach(func() {
				marketplace.PutProductReturns(nil, fmt.Errorf("put product failed"))
			})

			It("prints the error", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("put product failed"))
			})
		})
	})
})
