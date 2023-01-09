// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output/outputfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
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
		cmd.AttachPCAFile = ""
	})

	Describe("AttachChartCmd", func() {
		var (
			testProduct    *models.Product
			updatedProduct *models.Product
		)

		BeforeEach(func() {
			testProduct = test.CreateFakeProduct("", "My Super Product", "my-super-product", models.SolutionTypeChart)
			test.AddVersions(testProduct, "1.1.1")
			marketplace.GetProductWithVersionReturns(testProduct, &models.Version{Number: "1.1.1"}, nil)

			updatedProduct = test.CreateFakeProduct(testProduct.ProductId, "My Super Product", "my-super-product", models.SolutionTypeChart)
			test.AddVersions(updatedProduct, "1.1.1")
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
				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
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
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("attach local chart failed"))
				})
			})

			When("attaching a PCA file", func() {
				var uploader *internalfakes.FakeUploader
				BeforeEach(func() {
					uploader = &internalfakes.FakeUploader{}
					marketplace.GetUploaderReturns(uploader, nil)
					uploader.UploadMediaFileReturns("", "https://example.com/path/to/pca.pdf", nil)
				})
				It("attaches the PCA file", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachChartURL = "/path/to/my-chart"
					cmd.AttachInstructions = "helm install it"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
					Expect(err).ToNot(HaveOccurred())

					By("uploading the PCA file", func() {
						Expect(marketplace.GetUploaderCallCount()).To(Equal(1))
						Expect(uploader.UploadMediaFileCallCount()).To(Equal(1))
						Expect(uploader.UploadMediaFileArgsForCall(0)).To(Equal("/path/to/pca.pdf"))
					})

					By("adding the url to the product", func() {
						_, _, product, _ := marketplace.AttachLocalChartArgsForCall(0)
						Expect(product.PCADetails).ToNot(BeNil())
						Expect(product.PCADetails.URL).To(Equal("https://example.com/path/to/pca.pdf"))
						Expect(product.PCADetails.Version).To(Equal("1.1.1"))
					})
				})

				When("getting the uploader fails", func() {
					BeforeEach(func() {
						marketplace.GetUploaderReturns(nil, errors.New("get uploader failed"))
					})
					It("returns an error", func() {
						cmd.AttachProductSlug = "my-super-product"
						cmd.AttachProductVersion = "1.1.1"
						cmd.AttachChartURL = "/path/to/my-chart"
						cmd.AttachInstructions = "helm install it"
						cmd.AttachPCAFile = "/path/to/pca.pdf"
						err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("get uploader failed"))
					})
				})

				When("uploading the file fails", func() {
					BeforeEach(func() {
						uploader.UploadMediaFileReturns("", "", errors.New("upload media file failed"))
					})
					It("returns an error", func() {
						cmd.AttachProductSlug = "my-super-product"
						cmd.AttachProductVersion = "1.1.1"
						cmd.AttachChartURL = "/path/to/my-chart"
						cmd.AttachInstructions = "helm install it"
						cmd.AttachPCAFile = "/path/to/pca.pdf"
						err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("upload media file failed"))
					})
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
				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
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
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
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
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
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
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
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
						err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
						Expect(err).ToNot(HaveOccurred())

						By("passing a new version to upload vm", func() {
							_, _, _, version := marketplace.AttachPublicChartArgsForCall(0)
							Expect(version.Number).To(Equal("9.9.9"))
							Expect(version.IsNewVersion).To(BeTrue())
						})
					})
				})
			})

			When("product is the wrong solution type", func() {
				BeforeEach(func() {
					testProduct.SolutionType = models.SolutionTypeOthers
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
					cmd.AttachInstructions = "helm install it"
					err := cmd.AttachChartCmd.RunE(cmd.DownloadCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("cannot attach a chart to my-super-product which is of type OTHERS"))
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
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("render failed"))
				})
			})

			When("attaching a PCA file", func() {
				var uploader *internalfakes.FakeUploader
				BeforeEach(func() {
					uploader = &internalfakes.FakeUploader{}
					marketplace.GetUploaderReturns(uploader, nil)
					uploader.UploadMediaFileReturns("", "https://example.com/path/to/pca.pdf", nil)
				})
				It("attaches the PCA file", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
					cmd.AttachInstructions = "helm install it"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
					Expect(err).ToNot(HaveOccurred())

					By("uploading the PCA file", func() {
						Expect(marketplace.GetUploaderCallCount()).To(Equal(1))
						Expect(uploader.UploadMediaFileCallCount()).To(Equal(1))
						Expect(uploader.UploadMediaFileArgsForCall(0)).To(Equal("/path/to/pca.pdf"))
					})

					By("adding the url to the product", func() {
						_, _, product, _ := marketplace.AttachPublicChartArgsForCall(0)
						Expect(product.PCADetails).ToNot(BeNil())
						Expect(product.PCADetails.URL).To(Equal("https://example.com/path/to/pca.pdf"))
						Expect(product.PCADetails.Version).To(Equal("1.1.1"))
					})
				})

				When("getting the uploader fails", func() {
					BeforeEach(func() {
						marketplace.GetUploaderReturns(nil, errors.New("get uploader failed"))
					})
					It("returns an error", func() {
						cmd.AttachProductSlug = "my-super-product"
						cmd.AttachProductVersion = "1.1.1"
						cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
						cmd.AttachInstructions = "helm install it"
						cmd.AttachPCAFile = "/path/to/pca.pdf"
						err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("get uploader failed"))
					})
				})

				When("uploading the file fails", func() {
					BeforeEach(func() {
						uploader.UploadMediaFileReturns("", "", errors.New("upload media file failed"))
					})
					It("returns an error", func() {
						cmd.AttachProductSlug = "my-super-product"
						cmd.AttachProductVersion = "1.1.1"
						cmd.AttachChartURL = "https://example.com/public/my-chart.tgz"
						cmd.AttachInstructions = "helm install it"
						cmd.AttachPCAFile = "/path/to/pca.pdf"
						err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("upload media file failed"))
					})
				})
			})
		})
	})

	Describe("AttachContainerImageCmd", func() {
		var testProduct *models.Product

		BeforeEach(func() {
			testProduct = test.CreateFakeProduct("", "My Super Product", "my-super-product", models.SolutionTypeImage)
			test.AddVersions(testProduct, "1.1.1")
			marketplace.GetProductWithVersionReturns(testProduct, &models.Version{Number: "1.1.1"}, nil)

			updatedProduct := test.CreateFakeProduct(testProduct.ProductId, "My Super Product", "my-super-product", models.SolutionTypeImage)
			test.AddVersions(updatedProduct, "1.1.1")
			nginx := test.CreateFakeContainerImage("nginx", "1.21.6")
			test.AddContainerImages(updatedProduct, "1.1.1", "docker run it", nginx)
			marketplace.AttachLocalContainerImageReturns(updatedProduct, nil)
			marketplace.AttachPublicContainerImageReturns(updatedProduct, nil)

			cmd.AttachContainerImageFile = ""
			cmd.AttachPCAFile = ""
		})

		It("attaches the container image", func() {
			cmd.AttachProductSlug = "my-super-product"
			cmd.AttachProductVersion = "1.1.1"
			cmd.AttachContainerImage = "docker.io/bitnami/nginx"
			cmd.AttachContainerImageTag = "1.21.6"
			cmd.AttachContainerImageTagType = "FIXED"
			cmd.AttachInstructions = "docker run it"
			err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				slug, version := marketplace.GetProductWithVersionArgsForCall(0)
				Expect(slug).To(Equal("my-super-product"))
				Expect(version).To(Equal("1.1.1"))
			})

			By("attaching the container image", func() {
				Expect(marketplace.AttachPublicContainerImageCallCount()).To(Equal(1))
				image, tag, tagType, instructions, product, version := marketplace.AttachPublicContainerImageArgsForCall(0)
				Expect(image).To(Equal("docker.io/bitnami/nginx"))
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

		When("attaching a PCA file", func() {
			var uploader *internalfakes.FakeUploader
			BeforeEach(func() {
				uploader = &internalfakes.FakeUploader{}
				marketplace.GetUploaderReturns(uploader, nil)
				uploader.UploadMediaFileReturns("", "https://example.com/path/to/pca.pdf", nil)
			})
			It("attaches the PCA file", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "docker.io/bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				cmd.AttachPCAFile = "/path/to/pca.pdf"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("uploading the PCA file", func() {
					Expect(marketplace.GetUploaderCallCount()).To(Equal(1))
					Expect(uploader.UploadMediaFileCallCount()).To(Equal(1))
					Expect(uploader.UploadMediaFileArgsForCall(0)).To(Equal("/path/to/pca.pdf"))
				})

				By("adding the url to the product", func() {
					_, _, _, _, product, _ := marketplace.AttachPublicContainerImageArgsForCall(0)
					Expect(product.PCADetails).ToNot(BeNil())
					Expect(product.PCADetails.URL).To(Equal("https://example.com/path/to/pca.pdf"))
					Expect(product.PCADetails.Version).To(Equal("1.1.1"))
				})
			})

			When("getting the uploader fails", func() {
				BeforeEach(func() {
					marketplace.GetUploaderReturns(nil, errors.New("get uploader failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachContainerImage = "docker.io/bitnami/nginx"
					cmd.AttachContainerImageTag = "1.21.6"
					cmd.AttachContainerImageTagType = "FIXED"
					cmd.AttachInstructions = "docker run it"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("get uploader failed"))
				})
			})

			When("uploading the file fails", func() {
				BeforeEach(func() {
					uploader.UploadMediaFileReturns("", "", errors.New("upload media file failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachContainerImage = "docker.io/bitnami/nginx"
					cmd.AttachContainerImageTag = "1.21.6"
					cmd.AttachContainerImageTagType = "FIXED"
					cmd.AttachInstructions = "docker run it"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("upload media file failed"))
				})
			})
		})

		Context("local image file", func() {
			It("uploads and attaches the container image", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "docker.io/bitnami/nginx"
				cmd.AttachContainerImageFile = "/path/tp/image.tar"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("uploadings and attaching the container image", func() {
					Expect(marketplace.AttachLocalContainerImageCallCount()).To(Equal(1))
					imageFile, image, tag, tagType, instructions, product, version := marketplace.AttachLocalContainerImageArgsForCall(0)
					Expect(imageFile).To(Equal("/path/tp/image.tar"))
					Expect(image).To(Equal("docker.io/bitnami/nginx"))
					Expect(tag).To(Equal("1.21.6"))
					Expect(tagType).To(Equal("FIXED"))
					Expect(instructions).To(Equal("docker run it"))
					Expect(product.Slug).To(Equal("my-super-product"))
					Expect(version.Number).To(Equal("1.1.1"))
					Expect(version.IsNewVersion).To(BeFalse())
				})
			})
		})

		When("getting the product fails", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, errors.New("get product with version failed"))
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "docker.io/bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
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
				cmd.AttachContainerImage = "docker.io/bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have version 9.9.9"))
			})

			Context("But, we want to create the version", func() {
				It("attaches the asset, but with a new version", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachContainerImage = "docker.io/bitnami/nginx"
					cmd.AttachContainerImageTag = "1.21.6"
					cmd.AttachContainerImageTagType = "FIXED"
					cmd.AttachInstructions = "docker run it"
					cmd.AttachCreateVersion = true
					err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
					Expect(err).ToNot(HaveOccurred())

					By("passing a new version to upload vm", func() {
						_, _, _, _, _, version := marketplace.AttachPublicContainerImageArgsForCall(0)
						Expect(version.Number).To(Equal("9.9.9"))
						Expect(version.IsNewVersion).To(BeTrue())
					})
				})
			})
		})

		When("product is the wrong solution type", func() {
			BeforeEach(func() {
				testProduct.SolutionType = models.SolutionTypeOthers
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "docker.io/bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("cannot attach an image to my-super-product which is of type OTHERS"))
			})
		})

		When("attaching the container image fails", func() {
			BeforeEach(func() {
				marketplace.AttachPublicContainerImageReturns(nil, errors.New("attach container image failed"))
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "docker.io/bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("attach container image failed"))
			})
		})

		When("rendering fails", func() {
			BeforeEach(func() {
				output.RenderContainerImagesReturns(errors.New("render failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachContainerImage = "docker.io/bitnami/nginx"
				cmd.AttachContainerImageTag = "1.21.6"
				cmd.AttachContainerImageTagType = "FIXED"
				cmd.AttachInstructions = "docker run it"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("render failed"))
			})
		})
	})

	Describe("AttachOtherCmd", func() {
		var testProduct *models.Product

		BeforeEach(func() {
			testProduct = test.CreateFakeProduct("", "My Super Product", "my-super-product", models.SolutionTypeOthers)
			test.AddVersions(testProduct, "1.1.1")
			marketplace.GetProductWithVersionReturns(testProduct, &models.Version{Number: "1.1.1"}, nil)

			updatedProduct := test.CreateFakeProduct(testProduct.ProductId, "My Super Product", "my-super-product", models.SolutionTypeOthers)
			test.AddVersions(updatedProduct, "1.1.1")
			updatedProduct.AddOnFiles = append(updatedProduct.AddOnFiles, test.CreateFakeOtherFile("fake-file", "1.1.1"))
			marketplace.AttachOtherFileReturns(updatedProduct, nil)

			cmd.AttachPCAFile = ""
		})

		It("attaches the asset", func() {
			cmd.AttachProductSlug = "my-super-product"
			cmd.AttachProductVersion = "1.1.1"
			cmd.AttachOtherFile = "path/to/a/file.tgz"
			err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				slug, version := marketplace.GetProductWithVersionArgsForCall(0)
				Expect(slug).To(Equal("my-super-product"))
				Expect(version).To(Equal("1.1.1"))
			})

			By("uploading the other file", func() {
				Expect(marketplace.AttachOtherFileCallCount()).To(Equal(1))
				file, product, version := marketplace.AttachOtherFileArgsForCall(0)
				Expect(file).To(Equal("path/to/a/file.tgz"))
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(version.Number).To(Equal("1.1.1"))
				Expect(version.IsNewVersion).To(BeFalse())
			})

			By("outputting the updated list of other files", func() {
				Expect(output.PrintHeaderCallCount()).To(Equal(1))
				header := output.PrintHeaderArgsForCall(0)
				Expect(header).To(Equal("Other files for My Super Product 1.1.1:"))

				Expect(output.RenderAssetsCallCount()).To(Equal(1))
				files := output.RenderAssetsArgsForCall(0)
				Expect(files).To(HaveLen(1))
				Expect(files[0].Filename).To(Equal("fake-file"))
			})
		})

		When("attaching a PCA file", func() {
			var uploader *internalfakes.FakeUploader
			BeforeEach(func() {
				uploader = &internalfakes.FakeUploader{}
				marketplace.GetUploaderReturns(uploader, nil)
				uploader.UploadMediaFileReturns("", "https://example.com/path/to/pca.pdf", nil)
			})
			It("attaches the PCA file", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachOtherFile = "path/to/a/file.tgz"
				cmd.AttachPCAFile = "/path/to/pca.pdf"
				err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("uploading the PCA file", func() {
					Expect(marketplace.GetUploaderCallCount()).To(Equal(1))
					Expect(uploader.UploadMediaFileCallCount()).To(Equal(1))
					Expect(uploader.UploadMediaFileArgsForCall(0)).To(Equal("/path/to/pca.pdf"))
				})

				By("adding the url to the product", func() {
					_, product, _ := marketplace.AttachOtherFileArgsForCall(0)
					Expect(product.PCADetails).ToNot(BeNil())
					Expect(product.PCADetails.URL).To(Equal("https://example.com/path/to/pca.pdf"))
					Expect(product.PCADetails.Version).To(Equal("1.1.1"))
				})
			})

			When("getting the uploader fails", func() {
				BeforeEach(func() {
					marketplace.GetUploaderReturns(nil, errors.New("get uploader failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachOtherFile = "path/to/a/file.tgz"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("get uploader failed"))
				})
			})

			When("uploading the file fails", func() {
				BeforeEach(func() {
					uploader.UploadMediaFileReturns("", "", errors.New("upload media file failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachOtherFile = "path/to/a/file.tgz"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("upload media file failed"))
				})
			})
		})

		When("getting the product fails", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, errors.New("get product with version failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachOtherFile = "path/to/a/file.tgz"
				err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
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
				cmd.AttachOtherFile = "path/to/a/file.tgz"
				err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have version 9.9.9"))
			})

			Context("But, we want to create the version", func() {
				It("attaches the asset, but with a new version", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachOtherFile = "path/to/a/file.tgz"
					cmd.AttachCreateVersion = true
					err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
					Expect(err).ToNot(HaveOccurred())

					By("passing a new version to attach other file", func() {
						_, _, version := marketplace.AttachOtherFileArgsForCall(0)
						Expect(version.Number).To(Equal("9.9.9"))
						Expect(version.IsNewVersion).To(BeTrue())
					})
				})
			})
		})

		When("product is the wrong solution type", func() {
			BeforeEach(func() {
				testProduct.SolutionType = models.SolutionTypeISO
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachOtherFile = "path/to/a/file.tgz"
				err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("cannot attach an other file to my-super-product which is of type ISO"))
			})
		})

		When("uploading the file fails", func() {
			BeforeEach(func() {
				marketplace.AttachOtherFileReturns(nil, errors.New("attach other file failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachOtherFile = "path/to/a/file.tgz"
				err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("attach other file failed"))
			})
		})

		When("rendering fails", func() {
			BeforeEach(func() {
				output.RenderAssetsReturns(errors.New("render failed"))
			})
			It("Returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachOtherFile = "path/to/a/file.tgz"
				err := cmd.AttachOtherCmd.RunE(cmd.AttachOtherCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("render failed"))
			})
		})
	})

	Describe("AttachVMCmd", func() {
		var testProduct *models.Product

		BeforeEach(func() {
			testProduct = test.CreateFakeProduct("", "My Super Product", "my-super-product", models.SolutionTypeOVA)
			test.AddVersions(testProduct, "1.1.1")
			marketplace.GetProductWithVersionReturns(testProduct, &models.Version{Number: "1.1.1"}, nil)

			updatedProduct := test.CreateFakeProduct(testProduct.ProductId, "My Super Product", "my-super-product", models.SolutionTypeOVA)
			test.AddVersions(updatedProduct, "1.1.1")
			updatedProduct.ProductDeploymentFiles = append(updatedProduct.ProductDeploymentFiles, test.CreateFakeOVA("fake-ova", "1.1.1"))
			marketplace.UploadVMReturns(updatedProduct, nil)

			cmd.AttachPCAFile = ""
		})

		It("attaches the asset", func() {
			cmd.AttachProductSlug = "my-super-product"
			cmd.AttachProductVersion = "1.1.1"
			cmd.AttachVMFile = "path/to/a/file.iso"
			err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
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

		When("attaching a PCA file", func() {
			var uploader *internalfakes.FakeUploader
			BeforeEach(func() {
				uploader = &internalfakes.FakeUploader{}
				marketplace.GetUploaderReturns(uploader, nil)
				uploader.UploadMediaFileReturns("", "https://example.com/path/to/pca.pdf", nil)
			})
			It("attaches the PCA file", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachVMFile = "path/to/a/file.iso"
				cmd.AttachPCAFile = "/path/to/pca.pdf"
				err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("uploading the PCA file", func() {
					Expect(marketplace.GetUploaderCallCount()).To(Equal(1))
					Expect(uploader.UploadMediaFileCallCount()).To(Equal(1))
					Expect(uploader.UploadMediaFileArgsForCall(0)).To(Equal("/path/to/pca.pdf"))
				})

				By("adding the url to the product", func() {
					_, product, _ := marketplace.UploadVMArgsForCall(0)
					Expect(product.PCADetails).ToNot(BeNil())
					Expect(product.PCADetails.URL).To(Equal("https://example.com/path/to/pca.pdf"))
					Expect(product.PCADetails.Version).To(Equal("1.1.1"))
				})
			})

			When("getting the uploader fails", func() {
				BeforeEach(func() {
					marketplace.GetUploaderReturns(nil, errors.New("get uploader failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachVMFile = "path/to/a/file.iso"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("get uploader failed"))
				})
			})

			When("uploading the file fails", func() {
				BeforeEach(func() {
					uploader.UploadMediaFileReturns("", "", errors.New("upload media file failed"))
				})
				It("returns an error", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "1.1.1"
					cmd.AttachVMFile = "path/to/a/file.iso"
					cmd.AttachPCAFile = "/path/to/pca.pdf"
					err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("upload media file failed"))
				})
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
				err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
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
				err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have version 9.9.9"))
			})

			Context("But, we want to create the version", func() {
				It("attaches the asset, but with a new version", func() {
					cmd.AttachProductSlug = "my-super-product"
					cmd.AttachProductVersion = "9.9.9"
					cmd.AttachVMFile = "path/to/a/file.iso"
					cmd.AttachCreateVersion = true
					err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
					Expect(err).ToNot(HaveOccurred())

					By("passing a new version to upload vm", func() {
						_, _, version := marketplace.UploadVMArgsForCall(0)
						Expect(version.Number).To(Equal("9.9.9"))
						Expect(version.IsNewVersion).To(BeTrue())
					})
				})
			})
		})

		When("product is the wrong solution type", func() {
			BeforeEach(func() {
				testProduct.SolutionType = models.SolutionTypeOthers
			})
			It("returns an error", func() {
				cmd.AttachProductSlug = "my-super-product"
				cmd.AttachProductVersion = "1.1.1"
				cmd.AttachVMFile = "path/to/a/file.iso"
				err := cmd.AttachVMCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("cannot attach a vm to my-super-product which is of type OTHERS"))
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
				err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
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
				err := cmd.AttachVMCmd.RunE(cmd.AttachVMCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("render failed"))
			})
		})
	})
})
