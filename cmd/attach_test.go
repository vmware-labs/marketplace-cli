// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output/outputfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("ValidateTagType", func() {
	It("validates the tag type parameter", func() {
		By("accepting fixed", func() {
			cmd.AttachContainerImageTagType = "fixed"
			Expect(cmd.ValidateTagType(nil, nil)).To(Succeed())
		})

		By("accepting FIXED", func() {
			cmd.AttachContainerImageTagType = "FIXED"
			Expect(cmd.ValidateTagType(nil, nil)).To(Succeed())
		})

		By("accepting floating", func() {
			cmd.AttachContainerImageTagType = "floating"
			Expect(cmd.ValidateTagType(nil, nil)).To(Succeed())
		})

		By("accepting Floating", func() {
			cmd.AttachContainerImageTagType = "Floating"
			Expect(cmd.ValidateTagType(nil, nil)).To(Succeed())
		})

		By("rejecting anything else", func() {
			cmd.AttachContainerImageTagType = "imaginary"
			Expect(cmd.ValidateTagType(nil, nil)).ToNot(Succeed())
		})
	})
})

var _ = Describe("AttachCmd", func() {
	var (
		marketplace *pkgfakes.FakeMarketplaceInterface
		output      *outputfakes.FakeFormat
	)

	BeforeEach(func() {
		marketplace = &pkgfakes.FakeMarketplaceInterface{}
		output = &outputfakes.FakeFormat{}
		cmd.Marketplace = marketplace
		cmd.Output = output
		cmd.AttachCreateVersion = false
	})

	Describe("AttachChartCmd", func() {
		var (
			testProduct    *models.Product
			updatedProduct *models.Product
		)

		BeforeEach(func() {
			testProduct = test.CreateFakeProduct("", "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(testProduct, "1.1.1")
			marketplace.GetProductWithVersionReturns(testProduct, &models.Version{Number: "1.1.1"}, nil)

			updatedProduct = test.CreateFakeProduct(testProduct.ProductId, "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(updatedProduct, "1.1.1")
		})

		Context("chart is a local file", func() {
			BeforeEach(func() {
				updatedProduct.ChartVersions = append(updatedProduct.ChartVersions, &models.ChartVersion{
					AppVersion: "1.1.1",
					HelmTarUrl: "https://example.com/uploaded-chart.tgz",
					Readme:     "helm install it",
				})
				marketplace.AttachLocalChartReturns(updatedProduct, nil)
			})
			It("attaches the chart", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachChartURL = "/path/to/my-chart"
				cmd.AttachInstructions = "helm install it"
				err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("getting the product details", func() {
					Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
					slug, version := marketplace.GetProductWithVersionArgsForCall(0)
					Expect(slug).To(Equal("my-super-product"))
					Expect(version).To(Equal("1.1.1"))
				})

				By("attaching the local chart", func() {
					Expect(marketplace.AttachLocalChartCallCount()).To(Equal(1))
					chartUrl, instructions, product, version := marketplace.AttachLocalChartArgsForCall(0)
					Expect(chartUrl).To(Equal("/path/to/my-chart"))
					Expect(instructions).To(Equal("helm install it"))
					Expect(product.Slug).To(Equal("my-super-product"))
					Expect(version.Number).To(Equal("1.1.1"))
					Expect(version.IsNewVersion).To(BeFalse())
				})

				By("outputting the updated list of charts", func() {
					Expect(output.PrintHeaderCallCount()).To(Equal(1))
					header := output.PrintHeaderArgsForCall(0)
					Expect(header).To(Equal("Charts for My Super Product 1.1.1:"))

					Expect(output.RenderChartsCallCount()).To(Equal(1))
					charts := output.RenderChartsArgsForCall(0)
					Expect(charts).To(HaveLen(1))
					Expect(charts[0].AppVersion).To(Equal("1.1.1"))
					Expect(charts[0].HelmTarUrl).To(Equal("https://example.com/uploaded-chart.tgz"))
					Expect(charts[0].Readme).To(Equal("helm install it"))
				})
			})

			When("attaching the chart fails", func() {
				BeforeEach(func() {
					marketplace.AttachLocalChartReturns(nil, errors.New("attach local chart failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachChartURL = "/path/to/my-chart"
					cmd.AttachInstructions = "helm install it"
					err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("attach local chart failed"))
				})
			})
		})

		Context("chart is available on the network", func() {
			BeforeEach(func() {
				updatedProduct.ChartVersions = append(updatedProduct.ChartVersions, &models.ChartVersion{
					AppVersion: "1.1.1",
					HelmTarUrl: "https://example.com/public/my-chart.tgz",
					Readme:     "helm install it",
				})
				marketplace.AttachPublicChartReturns(updatedProduct, nil)
			})
			It("attaches the chart", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
				cmd.AttachInstructions = "helm install it"
				err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("getting the product details", func() {
					Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
					slug, version := marketplace.GetProductWithVersionArgsForCall(0)
					Expect(slug).To(Equal("my-super-product"))
					Expect(version).To(Equal("1.1.1"))
				})

				By("attaching the local chart", func() {
					Expect(marketplace.AttachPublicChartCallCount()).To(Equal(1))
					chartUrl, instructions, product, version := marketplace.AttachPublicChartArgsForCall(0)
					Expect(chartUrl.String()).To(Equal("https://example.com/public/my-chart.tgz"))
					Expect(instructions).To(Equal("helm install it"))
					Expect(product.Slug).To(Equal("my-super-product"))
					Expect(version.Number).To(Equal("1.1.1"))
					Expect(version.IsNewVersion).To(BeFalse())
				})

				By("outputting the updated list of charts", func() {
					Expect(output.PrintHeaderCallCount()).To(Equal(1))
					header := output.PrintHeaderArgsForCall(0)
					Expect(header).To(Equal("Charts for My Super Product 1.1.1:"))

					Expect(output.RenderChartsCallCount()).To(Equal(1))
					charts := output.RenderChartsArgsForCall(0)
					Expect(charts).To(HaveLen(1))
					Expect(charts[0].AppVersion).To(Equal("1.1.1"))
					Expect(charts[0].HelmTarUrl).To(Equal("https://example.com/public/my-chart.tgz"))
					Expect(charts[0].Readme).To(Equal("helm install it"))
				})
			})

			When("attaching the chart fails", func() {
				BeforeEach(func() {
					marketplace.AttachPublicChartReturns(nil, errors.New("attach public chart failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
					cmd.AttachInstructions = "helm install it"
					err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("attach public chart failed"))
				})
			})

			When("getting the product fails", func() {
				BeforeEach(func() {
					marketplace.GetProductWithVersionReturns(nil, nil, errors.New("get product with version failed"))
				})

				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
					cmd.AttachInstructions = "helm install it"
					err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("get product with version failed"))
				})
			})

			When("the version does not exist", func() {
				BeforeEach(func() {
					marketplace.GetProductWithVersionReturns(testProduct, nil, &pkg.VersionDoesNotExistError{
						Product: testProduct.Slug,
						Version: "9.9.9",
					})
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
					cmd.AttachInstructions = "helm install it"
					err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("product \"my-super-product\" does not have version 9.9.9"))
				})

				Context("But, we want to create the version", func() {
					It("attaches the asset, but with a new version", func() {
						cmd.AttachProductSlug = "my-super-product"
						cmd.AttachProductVersion = "9.9.9"
						cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
						cmd.AttachInstructions = "helm install it"
						cmd.AttachCreateVersion = true
						err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
						Expect(err).ToNot(HaveOccurred())

						By("passing a new version to upload vm", func() {
							_, _, _, version := marketplace.AttachPublicChartArgsForCall(0)
							Expect(version.Number).To(Equal("9.9.9"))
							Expect(version.IsNewVersion).To(BeTrue())
						})
					})
				})
			})

			When("rendering fails", func() {
				BeforeEach(func() {
					output.RenderChartsReturns(errors.New("render failed"))
				})
				It("Returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
					cmd.AttachInstructions = "helm install it"
					err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("render failed"))
				})
			})
		})
	})

	Describe("AttachContainerImageCmd", func() {
		var testProduct *models.Product

		BeforeEach(func() {
			testProduct = test.CreateFakeProduct("", "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(testProduct, "1.1.1")
			marketplace.GetProductWithVersionReturns(testProduct, &models.Version{Number: "1.1.1"}, nil)

			updatedProduct := test.CreateFakeProduct(testProduct.ProductId, "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(updatedProduct, "1.1.1")
			nginx := test.CreateFakeContainerImage("nginx", "1.21.6")
			test.AddContainerImages(updatedProduct, "1.1.1", "docker run it", nginx)
			marketplace.AttachPublicContainerImageReturns(updatedProduct, nil)
		})

		It("attaches the container image", func() {
			cmd.AttachProductSlug = "my-super-product"
			cmd.AttachProductVersion = "1.1.1"
			cmd.AttachContainerImage = "bitnami/nginx"
			cmd.AttachContainerImageTag = "1.21.6"
			cmd.AttachContainerImageTagType = "FIXED"
			cmd.AttachInstructions = "docker run it"
			err := cmd.AttachContainerImageCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				slug, version := marketplace.GetProductWithVersionArgsForCall(0)
				Expect(slug).To(Equal("my-super-product"))
				Expect(version).To(Equal("1.1.1"))
			})

			By("uploading the vm", func() {
				Expect(marketplace.AttachPublicContainerImageCallCount()).To(Equal(1))
				image, tag, tagType, instructions, product, version := marketplace.AttachPublicContainerImageArgsForCall(0)
				Expect(image).To(Equal("bitnami/nginx"))
				Expect(tag).To(Equal("1.21.6"))
				Expect(tagType).To(Equal("FIXED"))
				Expect(instructions).To(Equal("docker run it"))
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(version.Number).To(Equal("1.1.1"))
				Expect(version.IsNewVersion).To(BeFalse())
			})

			By("outputting the updated list of vms", func() {
				Expect(output.PrintHeaderCallCount()).To(Equal(1))
				header := output.PrintHeaderArgsForCall(0)
				Expect(header).To(Equal("Container images for My Super Product 1.1.1:"))

				Expect(output.RenderContainerImagesCallCount()).To(Equal(1))
				images := output.RenderContainerImagesArgsForCall(0)
				Expect(images).To(HaveLen(1))
				Expect(images[0].AppVersion).To(Equal("1.1.1"))
				Expect(images[0].DockerURLs[0].Url).To(Equal("nginx"))
				Expect(images[0].DockerURLs[0].ImageTags[0].Tag).To(Equal("1.21.6"))
				Expect(images[0].DockerURLs[0].ImageTags[0].Type).To(Equal("FIXED"))
			})
		})

		When("getting the product fails", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, errors.New("get product with version failed"))
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product with version failed"))
			})
		})

		When("the version does not exist", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(testProduct, nil, &pkg.VersionDoesNotExistError{
					Product: testProduct.Slug,
					Version: "9.9.9",
				})
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "9.9.9"
				cmd.AttachContainerImage = "bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have version 9.9.9"))
			})

			Context("But, we want to create the version", func() {
				It("attaches the asset, but with a new version", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachContainerImage = "bitnami/nginx"
					cmd.AttachContainerImageTag = "1.21.6"
					cmd.AttachContainerImageTagType = "FIXED"
					cmd.AttachInstructions = "docker run it"
					cmd.AttachCreateVersion = true
					err := cmd.AttachContainerImageCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).ToNot(HaveOccurred())

					By("passing a new version to upload vm", func() {
						_, _, _, _, _, version := marketplace.AttachPublicContainerImageArgsForCall(0)
						Expect(version.Number).To(Equal("9.9.9"))
						Expect(version.IsNewVersion).To(BeTrue())
					})
				})
			})
		})

		When("attaching the container image fails", func() {
			BeforeEach(func() {
				marketplace.AttachPublicContainerImageReturns(nil, errors.New("upload vm failed"))
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("upload vm failed"))
			})
		})

		When("rendering fails", func() {
			BeforeEach(func() {
				output.RenderContainerImagesReturns(errors.New("render failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("render failed"))
			})
		})
	})

	Describe("AttachVMCmd", func() {
		var testProduct *models.Product

		BeforeEach(func() {
			testProduct = test.CreateFakeProduct("", "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(testProduct, "1.1.1")
			marketplace.GetProductWithVersionReturns(testProduct, &models.Version{Number: "1.1.1"}, nil)

			updatedProduct := test.CreateFakeProduct(testProduct.ProductId, "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(updatedProduct, "1.1.1")
			updatedProduct.ProductDeploymentFiles = append(updatedProduct.ProductDeploymentFiles, test.CreateFakeOVA("fake-ova", "1.1.1"))
			marketplace.UploadVMReturns(updatedProduct, nil)
		})

		It("attaches the asset", func() {
			cmd.AttachProductSlug = "my-super-product"
			cmd.AttachProductVersion = "1.1.1"
			cmd.AttachVMFile = "path/to/a/file.iso"
			err := cmd.AttachVMCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				slug, version := marketplace.GetProductWithVersionArgsForCall(0)
				Expect(slug).To(Equal("my-super-product"))
				Expect(version).To(Equal("1.1.1"))
			})

			By("uploading the vm", func() {
				Expect(marketplace.UploadVMCallCount()).To(Equal(1))
				vmFile, product, version := marketplace.UploadVMArgsForCall(0)
				Expect(vmFile).To(Equal("path/to/a/file.iso"))
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(version.Number).To(Equal("1.1.1"))
				Expect(version.IsNewVersion).To(BeFalse())
			})

			By("outputting the updated list of vms", func() {
				Expect(output.PrintHeaderCallCount()).To(Equal(1))
				header := output.PrintHeaderArgsForCall(0)
				Expect(header).To(Equal("Virtual machine files for My Super Product 1.1.1:"))

				Expect(output.RenderFilesCallCount()).To(Equal(1))
				files := output.RenderFilesArgsForCall(0)
				Expect(files).To(HaveLen(1))
				Expect(files[0].Name).To(Equal("fake-ova"))
			})
		})

		When("getting the product fails", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, errors.New("get product with version failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachVMFile = "path/to/a/file.iso"
				err := cmd.AttachVMCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product with version failed"))
			})
		})

		When("the version does not exist", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(testProduct, nil, &pkg.VersionDoesNotExistError{
					Product: testProduct.Slug,
					Version: "9.9.9",
				})
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "9.9.9"
				cmd.AttachVMFile = "path/to/a/file.iso"
				err := cmd.AttachVMCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have version 9.9.9"))
			})

			Context("But, we want to create the version", func() {
				It("attaches the asset, but with a new version", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachVMFile = "path/to/a/file.iso"
					cmd.AttachCreateVersion = true
					err := cmd.AttachVMCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).ToNot(HaveOccurred())

					By("passing a new version to upload vm", func() {
						_, _, version := marketplace.UploadVMArgsForCall(0)
						Expect(version.Number).To(Equal("9.9.9"))
						Expect(version.IsNewVersion).To(BeTrue())
					})
				})
			})
		})

		When("uploading the VM fails", func() {
			BeforeEach(func() {
				marketplace.UploadVMReturns(nil, errors.New("upload vm failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachVMFile = "path/to/a/file.iso"
				err := cmd.AttachVMCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("upload vm failed"))
			})
		})

		When("rendering fails", func() {
			BeforeEach(func() {
				output.RenderFilesReturns(errors.New("render failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachVMFile = "path/to/a/file.iso"
				err := cmd.AttachVMCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("render failed"))
			})
		})
	})
})
