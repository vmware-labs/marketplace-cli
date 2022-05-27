// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
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
			cmd.ListContainerImageCmd.SetErr(gbytes.NewBuffer())
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
				Expect(images).To(HaveLen(1))
				Expect(images[0].AppVersion).To(Equal("1.2.3"))
				Expect(images[0].DockerURLs).To(HaveLen(1))
				Expect(images[0].DockerURLs[0].ImageTags).To(HaveLen(2))
				Expect(images[0].DockerURLs[0].ImageTags[0].Tag).To(Equal("0.0.1"))
				Expect(images[0].DockerURLs[0].ImageTags[0].Type).To(Equal("FIXED"))
				Expect(images[0].DockerURLs[0].ImageTags[1].Tag).To(Equal("latest"))
				Expect(images[0].DockerURLs[0].ImageTags[1].Type).To(Equal("FLOATING"))
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
})
