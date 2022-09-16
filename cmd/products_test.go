// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output/outputfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Products", func() {
	var (
		marketplace *pkgfakes.FakeMarketplaceInterface
		output      *outputfakes.FakeFormat
	)

	BeforeEach(func() {
		marketplace = &pkgfakes.FakeMarketplaceInterface{}
		cmd.Marketplace = marketplace

		output = &outputfakes.FakeFormat{}
		cmd.Output = output
	})

	Describe("ListProductsCmd", func() {
		BeforeEach(func() {
			products := []*models.Product{
				test.CreateFakeProduct(
					"",
					"My Super Product",
					"my-super-product",
					models.SolutionTypeOthers),
				test.CreateFakeProduct(
					"",
					"My Other Product",
					"my-other-product",
					models.SolutionTypeOthers),
			}

			cmd.ListProductSearchText = ""
			cmd.ListProductsAllOrgs = false
			marketplace.ListProductsReturns(products, nil)
		})

		It("outputs the list of products", func() {
			err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			By("getting the list of products from the Marketplace", func() {
				Expect(marketplace.ListProductsCallCount()).To(Equal(1))
				filter := marketplace.ListProductsArgsForCall(0)
				Expect(filter.AllOrgs).To(BeFalse())
				Expect(filter.Text).To(Equal(""))
			})

			By("outputting the response", func() {
				Expect(output.PrintHeaderCallCount()).To(Equal(1))
				Expect(output.PrintHeaderArgsForCall(0)).To(Equal("All products from my-org"))

				Expect(output.RenderProductsCallCount()).To(Equal(1))
				products := output.RenderProductsArgsForCall(0)
				Expect(products).To(HaveLen(2))
				Expect(products[0].Slug).To(Equal("my-super-product"))
				Expect(products[1].Slug).To(Equal("my-other-product"))
			})
		})

		Context("Using all orgs and a search field", func() {
			It("sends the appropriate filter", func() {
				cmd.ListProductSearchText = "tanzu"
				cmd.ListProductsAllOrgs = true
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
				Expect(err).ToNot(HaveOccurred())

				By("using the right filter", func() {
					Expect(marketplace.ListProductsCallCount()).To(Equal(1))
					filter := marketplace.ListProductsArgsForCall(0)
					Expect(filter.AllOrgs).To(BeTrue())
					Expect(filter.Text).To(Equal("tanzu"))
				})

				By("outputting a specific header", func() {
					Expect(output.PrintHeaderCallCount()).To(Equal(1))
					Expect(output.PrintHeaderArgsForCall(0)).To(Equal("All products from all organizations filtered by \"tanzu\""))
				})
			})
		})

		Context("Error getting the product list", func() {
			BeforeEach(func() {
				marketplace.ListProductsReturns([]*models.Product{}, fmt.Errorf("gettings products failed"))
			})

			It("prints the error", func() {
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("gettings products failed"))
			})
		})
	})

	Describe("GetProductCmd", func() {
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				models.SolutionTypeOthers)
			test.AddVersions(product, "1.2.3", "2.3.4")
			marketplace.GetProductReturns(product, nil)
			marketplace.GetProductWithVersionStub = func(slug string, version string) (*models.Product, *models.Version, error) {
				Expect(slug).To(Equal("my-super-product"))
				Expect(version).To(Equal("1.2.3"))
				return product, &models.Version{Number: "1.2.3"}, nil
			}
		})

		It("outputs the product", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = ""
			err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product from the Marketplace", func() {
				Expect(marketplace.GetProductCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderProductCallCount()).To(Equal(1))
				product, version := output.RenderProductArgsForCall(0)
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(product.DisplayName).To(Equal("My Super Product"))
				Expect(version.Number).To(Equal("2.3.4"))
			})
		})

		Context("Version number given", func() {
			It("outputs the product", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("getting the product from the Marketplace", func() {
					Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				})

				By("outputting the response", func() {
					Expect(output.RenderProductCallCount()).To(Equal(1))
					product, version := output.RenderProductArgsForCall(0)
					Expect(product.Slug).To(Equal("my-super-product"))
					Expect(product.DisplayName).To(Equal("My Super Product"))
					Expect(version.Number).To(Equal("1.2.3"))
				})
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductReturns(nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = ""
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})

	Describe("ListAssetsCmd", func() {
		var product *models.Product
		BeforeEach(func() {
			product = test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeOVA)
			version := &models.Version{Number: "1"}
			test.AddVersions(product, "1")
			vm := test.CreateFakeOVA("hyperspace-database.ova", "1")
			product.ProductDeploymentFiles = append(product.ProductDeploymentFiles, vm)

			metafile := test.CreateFakeMetaFile("deploy.sh", "0.0.1", "1")
			product.MetaFiles = append(product.MetaFiles, metafile)

			marketplace.GetProductWithVersionReturns(product, version, nil)
		})

		It("outputs the list of assets", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "1"
			cmd.AssetType = ""
			err := cmd.ListAssetsCmd.RunE(cmd.ListAssetsCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product from the Marketplace", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				slug, version := marketplace.GetProductWithVersionArgsForCall(0)
				Expect(slug).To(Equal("my-super-product"))
				Expect(version).To(Equal("1"))
			})

			By("outputting the response", func() {
				Expect(output.RenderAssetsCallCount()).To(Equal(1))
				assets := output.RenderAssetsArgsForCall(0)
				Expect(assets).To(HaveLen(2))
			})
		})

		Context("an asset type filter is used", func() {
			It("outputs the filtered list of assets", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1"
				cmd.AssetType = "vm"
				err := cmd.ListAssetsCmd.RunE(cmd.ListAssetsCmd, []string{})
				Expect(err).ToNot(HaveOccurred())

				By("outputting the filtered response", func() {
					Expect(output.RenderAssetsCallCount()).To(Equal(1))
					assets := output.RenderAssetsArgsForCall(0)
					Expect(assets).To(HaveLen(1))
				})
			})
		})
	})

	Describe("ListProductVersionsCmd", func() {
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				models.SolutionTypeOVA)
			test.AddVersions(product, "0.1.2", "1.2.3")
			marketplace.GetProductReturns(product, nil)
		})

		It("outputs the list of versions", func() {
			cmd.ProductSlug = "my-super-product"
			err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product from the Marketplace", func() {
				Expect(marketplace.GetProductCallCount()).To(Equal(1))
				Expect(marketplace.GetProductArgsForCall(0)).To(Equal("my-super-product"))
			})

			By("outputting the response", func() {
				Expect(output.RenderVersionsCallCount()).To(Equal(1))
				product := output.RenderVersionsArgsForCall(0)
				Expect(product.Slug).To(Equal("my-super-product"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductReturns(nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})
})
