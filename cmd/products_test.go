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
					"PENDING"),
				test.CreateFakeProduct(
					"",
					"My Other Product",
					"my-other-product",
					"PENDING"),
			}

			marketplace.ListProductsReturns(products, nil)
		})

		It("outputs the list of products", func() {
			err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			By("getting the list of products from the Marketplace", func() {
				Expect(marketplace.ListProductsCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderProductsCallCount()).To(Equal(1))
				products := output.RenderProductsArgsForCall(0)
				Expect(products).To(HaveLen(2))
				Expect(products[0].Slug).To(Equal("my-super-product"))
				Expect(products[1].Slug).To(Equal("my-other-product"))
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
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")
			marketplace.GetProductReturns(product, nil)
		})

		It("outputs the product", func() {
			cmd.ProductSlug = "my-super-product"
			err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product from the Marketplace", func() {
				Expect(marketplace.GetProductCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderProductCallCount()).To(Equal(1))
				product := output.RenderProductArgsForCall(0)
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(product.DisplayName).To(Equal("My Super Product"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductReturns(nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})

	Describe("AddProductVersionCmd", func() {
		var productID string
		BeforeEach(func() {
			productID = uuid.New().String()
			product := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "0.1.2", "1.2.3")

			updatedProduct := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(updatedProduct, "0.1.2", "1.2.3", "9.9.9")

			marketplace.GetProductReturns(product, nil)
			marketplace.PutProductReturns(updatedProduct, nil)
		})

		It("adds the new version", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "9.9.9"
			err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("first, getting the existing product", func() {
				Expect(marketplace.GetProductCallCount()).To(Equal(1))
			})

			By("second, sending the new product", func() {
				Expect(marketplace.PutProductCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderVersionsCallCount()).To(Equal(1))
				product := output.RenderVersionsArgsForCall(0)
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(product.AllVersions).To(HaveLen(3))
			})
		})

		Context("Version already exists", func() {
			It("says that the version already exists", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" already has version 1.2.3"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductReturns(nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})

		Context("Error putting product", func() {
			BeforeEach(func() {
				marketplace.PutProductReturns(nil, fmt.Errorf("put product failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("put product failed"))
			})
		})
	})

	Describe("ListProductVersionsCmd", func() {
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "0.1.2", "1.2.3")
			marketplace.GetProductReturns(product, nil)
		})

		It("outputs the list of versions", func() {
			cmd.ProductSlug = "my-super-product"
			err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product from the Marketplace", func() {
				Expect(marketplace.GetProductCallCount()).To(Equal(1))
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
