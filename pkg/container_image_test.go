// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Container image", func() {
	var (
		httpClient  *pkgfakes.FakeHTTPClient
		marketplace *pkg.Marketplace
		uploader    *internalfakes.FakeUploader
	)

	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			Client: httpClient,
			Host:   "marketplace.vmware.example",
		}
		uploader = &internalfakes.FakeUploader{}
		marketplace.SetUploader(uploader)
	})

	Describe("AttachPublicContainerImage", func() {
		BeforeEach(func() {
			httpClient.DoStub = PutProductEchoResponse
		})

		It("updates the product with a public container image", func() {
			product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
			test.AddVerions(product, "1.2.3")

			updatedProduct, err := marketplace.AttachPublicContainerImage("nginx", "latest", "FLOATING", "docker run it", product, &models.Version{Number: "1.2.3"})
			Expect(err).ToNot(HaveOccurred())

			By("updating the product in the marketplace", func() {
				images := updatedProduct.GetContainerImagesForVersion("1.2.3")
				Expect(images).To(HaveLen(1))
				Expect(images[0].AppVersion).To(Equal("1.2.3"))
				Expect(images[0].DockerURLs).To(HaveLen(1))
				Expect(images[0].DockerURLs[0].Url).To(Equal("nginx"))
				Expect(images[0].DockerURLs[0].DeploymentInstruction).To(Equal("docker run it"))
				Expect(images[0].DockerURLs[0].DockerType).To(Equal(models.DockerTypeRegistry))
				Expect(images[0].DockerURLs[0].ImageTags).To(HaveLen(1))
				Expect(images[0].DockerURLs[0].ImageTags[0].Tag).To(Equal("latest"))
				Expect(images[0].DockerURLs[0].ImageTags[0].Type).To(Equal("FLOATING"))
			})
		})

		Context("the container and tag combo already exists for this version", func() {
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1.2.3")
				test.AddContainerImages(product, "1.2.3", "docker run it", &models.DockerURLDetails{
					Url: "nginx",
					ImageTags: []*models.DockerImageTag{
						{
							Tag:  "latest",
							Type: "FLOATING",
						},
					},
				})

				_, err := marketplace.AttachPublicContainerImage("nginx", "latest", "FLOATING", "docker run it", product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("hyperspace-database 1.2.3 already has the image nginx:latest"))
			})
		})

		Context("Updating the product fails", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("put product failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1.2.3")

				_, err := marketplace.AttachPublicContainerImage("nginx", "latest", "FLOATING", "docker run it", product, &models.Version{Number: "1.2.3"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the update for product \"hyperspace-database\" failed: marketplace request failed: put product failed"))
			})
		})
	})
})
