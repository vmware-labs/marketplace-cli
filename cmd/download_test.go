// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("DownloadCmd", func() {
	var (
		marketplace *pkgfakes.FakeMarketplaceInterface
		product     *models.Product
		productId   string
	)

	BeforeEach(func() {
		productId = uuid.New().String()
		product = test.CreateFakeProduct(productId, "My Super Product", "my-super-product", "PENDING")

		// Version 0.0.0 has no assets, 1.1.1 has one asset, 3.3.3 has multiple assets
		test.AddVerions(product, "0.0.0", "1.1.1", "3.3.3")

		product.AddFile(test.CreateFakeOVA("my-db.ova", "1.1.1"))

		product.AddFile(test.CreateFakeOVA("aaa.txt", "3.3.3"))
		product.AddFile(test.CreateFakeOVA("bbb.txt", "3.3.3"))
		product.AddFile(test.CreateFakeOVA("ccc.txt", "3.3.3"))

		marketplace = &pkgfakes.FakeMarketplaceInterface{}
		cmd.Marketplace = marketplace

		marketplace.GetProductWithVersionStub = func(slug string, version string) (*models.Product, *models.Version, error) {
			Expect(slug).To(Equal("my-super-product"))
			return product, &models.Version{Number: version}, nil
		}

		cmd.DownloadFilename = ""
		cmd.DownloadFilter = ""
	})

	It("downloads the asset", func() {
		cmd.DownloadProductSlug = "my-super-product"
		cmd.DownloadProductVersion = "1.1.1"
		err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
		Expect(err).ToNot(HaveOccurred())

		By("getting the product details", func() {
			Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
		})

		By("downloading the asset", func() {
			Expect(marketplace.DownloadCallCount()).To(Equal(1))
			downloadedProductId, filename, assetPayload := marketplace.DownloadArgsForCall(0)
			Expect(downloadedProductId).To(Equal(productId))
			Expect(filename).To(Equal("my-db.ova"))
			Expect(assetPayload.ProductId).To(Equal(productId))
			Expect(assetPayload.AppVersion).To(Equal("1.1.1"))
			Expect(assetPayload.EulaAccepted).To(BeTrue())
			Expect(assetPayload.DeploymentFileId).ToNot(BeEmpty())
		})
	})

	Context("Filename parameter used", func() {
		It("downloads the asset using the given filename", func() {
			cmd.DownloadProductSlug = "my-super-product"
			cmd.DownloadProductVersion = "1.1.1"
			cmd.DownloadFilename = "overridden-filename.ova"
			err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			Expect(marketplace.DownloadCallCount()).To(Equal(1))
			_, filename, _ := marketplace.DownloadArgsForCall(0)
			Expect(filename).To(Equal("overridden-filename.ova"))
		})
	})

	Context("Failed to get the product", func() {
		BeforeEach(func() {
			marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
		})
		It("returns an error", func() {
			err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("get product failed"))
		})
	})

	Context("Failed to download the asset", func() {
		BeforeEach(func() {
			marketplace.DownloadReturns(fmt.Errorf("download failed"))
		})
		It("returns an error", func() {
			err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("download failed"))
		})
	})

	Context("No assets attached to the product", func() {
		It("returns an error", func() {
			cmd.DownloadProductSlug = "my-super-product"
			cmd.DownloadProductVersion = "0.0.0"
			err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("product my-super-product 0.0.0 does not have any downloadable assets"))
		})
	})

	Context("Multiple assets attached to the product", func() {
		It("returns an error", func() {
			cmd.DownloadProductSlug = "my-super-product"
			cmd.DownloadProductVersion = "3.3.3"
			err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("product my-super-product 3.3.3 has multiple downloadable assets, please use the --filter parameter"))
		})
	})

	Context("Using a filter", func() {
		It("downloads the asset matching the filter", func() {
			cmd.DownloadProductSlug = "my-super-product"
			cmd.DownloadProductVersion = "3.3.3"
			cmd.DownloadFilter = "bbb"
			err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("downloading the chosen asset", func() {
				Expect(marketplace.DownloadCallCount()).To(Equal(1))
				downloadedProductId, filename, assetPayload := marketplace.DownloadArgsForCall(0)
				Expect(downloadedProductId).To(Equal(productId))
				Expect(filename).To(Equal("bbb.txt"))
				Expect(assetPayload.ProductId).To(Equal(productId))
				Expect(assetPayload.AppVersion).To(Equal("3.3.3"))
				Expect(assetPayload.EulaAccepted).To(BeTrue())
				Expect(assetPayload.DeploymentFileId).ToNot(BeEmpty())
			})
		})

		Context("No assets matching the filter", func() {
			It("returns an error", func() {
				cmd.DownloadProductSlug = "my-super-product"
				cmd.DownloadProductVersion = "3.3.3"
				cmd.DownloadFilter = "does not match"
				err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product my-super-product 3.3.3 does not have any downloadable assets that match the filter \"does not match\", please adjust the --filter parameter"))
			})
		})

		Context("Multiple assets matching the filter", func() {
			It("returns an error", func() {
				cmd.DownloadProductSlug = "my-super-product"
				cmd.DownloadProductVersion = "3.3.3"
				cmd.DownloadFilter = "txt"
				err := cmd.DownloadCmd.RunE(cmd.DownloadCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product my-super-product 3.3.3 has multiple downloadable assets that match the filter \"txt\", please adjust the --filter parameter"))
			})
		})
	})
})
