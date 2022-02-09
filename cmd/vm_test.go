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

var _ = Describe("Virtual Machines", func() {
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

	Describe("ListVMCmd", func() {
		BeforeEach(func() {
			product = test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")
			product.ProductDeploymentFiles = append(product.ProductDeploymentFiles, test.CreateFakeOVA("fake-ova", "1.2.3"))
		})

		It("outputs the vm files", func() {
			cmd.VMProductSlug = "my-super-product"
			cmd.VMProductVersion = "1.2.3"
			err := cmd.ListVMCmd.RunE(cmd.ListVMCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderFilesCallCount()).To(Equal(1))
				files := output.RenderFilesArgsForCall(0)
				Expect(files).To(HaveLen(1))
				Expect(files[0].AppVersion).To(Equal("1.2.3"))
			})
		})

		Context("Error fetching products", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.VMProductSlug = "my-super-product"
				cmd.VMProductVersion = "1.2.3"
				err := cmd.ListVMCmd.RunE(cmd.ListVMCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})
})
