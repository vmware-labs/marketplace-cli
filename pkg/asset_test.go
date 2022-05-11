// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"strconv"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Asset", func() {
	Describe("GetAssets", func() {
		Context("Chart", func() {
			var (
				product *models.Product
				chart   *models.ChartVersion
			)
			BeforeEach(func() {
				product = test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1")
				chart = &models.ChartVersion{
					Id:         uuid.New().String(),
					Version:    "1.0.0",
					AppVersion: "1",
					HelmTarUrl: "https://mychart.registry.examle.com/hyperspace-database-1.0.0.tgz",
					Repo: &models.Repo{
						Name: "My registry",
					},
					Size:                           5000,
					DownloadCount:                  10,
					IsUpdatedInMarketplaceRegistry: true,
				}
				product.ChartVersions = []*models.ChartVersion{chart}
			})

			It("returns the chart", func() {
				assets := pkg.GetAssets(product, "1")
				Expect(assets).To(HaveLen(1))

				Expect(assets[0].DisplayName).To(Equal("https://mychart.registry.examle.com/hyperspace-database-1.0.0.tgz"))
				Expect(assets[0].Filename).To(Equal("chart.tgz"))
				Expect(assets[0].Version).To(Equal("1.0.0"))
				Expect(assets[0].Type).To(Equal("Chart"))
				Expect(strconv.FormatInt(assets[0].Size, 10)).To(Equal("5000"))
				Expect(strconv.FormatInt(assets[0].Downloads, 10)).To(Equal("10"))
				Expect(assets[0].Downloadable).To(BeTrue())

				Expect(assets[0].DownloadRequestPayload.ProductId).To(Equal(product.ProductId))
				Expect(assets[0].DownloadRequestPayload.AppVersion).To(Equal("1"))
				Expect(assets[0].DownloadRequestPayload.ChartVersion).To(Equal("1.0.0"))
			})
		})

		Context("Container image", func() {
			var (
				product        *models.Product
				containerImage *models.DockerURLDetails
			)
			BeforeEach(func() {
				product = test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1")
				containerImage = test.CreateFakeContainerImage("astrowidgets/hyperspacedb", "1", "imaginary")
				test.AddContainerImages(product, "1", "docker run it", containerImage)
			})

			It("returns the container image", func() {
				assets := pkg.GetAssets(product, "1")
				Expect(assets).To(HaveLen(2))

				By("having an entry for each tag", func() {
					Expect(assets[0].DisplayName).To(Equal("astrowidgets/hyperspacedb:1"))
					Expect(assets[0].Filename).To(Equal("image.tar"))
					Expect(assets[0].Version).To(Equal("1"))
					Expect(assets[0].Type).To(Equal("Container Image"))
					Expect(strconv.FormatInt(assets[0].Size, 10)).To(Equal("12345"))
					Expect(strconv.FormatInt(assets[0].Downloads, 10)).To(Equal("15"))
					Expect(assets[0].Downloadable).To(BeTrue())

					Expect(assets[0].DownloadRequestPayload.ProductId).To(Equal(product.ProductId))
					Expect(assets[0].DownloadRequestPayload.AppVersion).To(Equal("1"))
					Expect(assets[0].DownloadRequestPayload.DockerlinkVersionID).To(Equal(product.DockerLinkVersions[0].ID))
					Expect(assets[0].DownloadRequestPayload.DockerUrlId).To(Equal(containerImage.ID))
					Expect(assets[0].DownloadRequestPayload.ImageTagId).To(Equal(containerImage.ImageTags[0].ID))

					Expect(assets[1].DisplayName).To(Equal("astrowidgets/hyperspacedb:imaginary"))
					Expect(assets[1].Filename).To(Equal("image.tar"))
					Expect(assets[1].Version).To(Equal("imaginary"))
					Expect(assets[1].Type).To(Equal("Container Image"))
					Expect(strconv.FormatInt(assets[1].Size, 10)).To(Equal("12345"))
					Expect(strconv.FormatInt(assets[1].Downloads, 10)).To(Equal("15"))
					Expect(assets[1].Downloadable).To(BeTrue())

					Expect(assets[1].DownloadRequestPayload.ProductId).To(Equal(product.ProductId))
					Expect(assets[1].DownloadRequestPayload.AppVersion).To(Equal("1"))
					Expect(assets[1].DownloadRequestPayload.DockerlinkVersionID).To(Equal(product.DockerLinkVersions[0].ID))
					Expect(assets[1].DownloadRequestPayload.DockerUrlId).To(Equal(containerImage.ID))
					Expect(assets[1].DownloadRequestPayload.ImageTagId).To(Equal(containerImage.ImageTags[1].ID))
				})
			})
		})

		Context("VM and MetaFile", func() {
			var (
				product  *models.Product
				vm       *models.ProductDeploymentFile
				metafile *models.MetaFile
			)
			BeforeEach(func() {
				product = test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", "PENDING")
				test.AddVerions(product, "1")
				vm = test.CreateFakeOVA("hyperspace-database.ova", "1")
				product.ProductDeploymentFiles = append(product.ProductDeploymentFiles, vm)

				metafile = test.CreateFakeMetaFile("deploy.sh", "0.0.1", "1")
				product.MetaFiles = append(product.MetaFiles, metafile)
			})

			It("returns the aggregated list of attached assets", func() {
				assets := pkg.GetAssets(product, "1")
				Expect(assets).To(HaveLen(2))

				By("including the VM file", func() {
					Expect(assets[0].DisplayName).To(Equal("hyperspace-database.ova"))
					Expect(assets[0].Filename).To(Equal("hyperspace-database.ova"))
					Expect(assets[0].Version).To(BeEmpty())
					Expect(assets[0].Type).To(Equal("VM"))
					Expect(strconv.FormatInt(assets[0].Size, 10)).To(Equal("1000100"))
					Expect(strconv.FormatInt(assets[0].Downloads, 10)).To(Equal("20"))
					Expect(assets[0].Downloadable).To(BeTrue())

					Expect(assets[0].DownloadRequestPayload.ProductId).To(Equal(product.ProductId))
					Expect(assets[0].DownloadRequestPayload.AppVersion).To(Equal("1"))
					Expect(assets[0].DownloadRequestPayload.DeploymentFileId).To(Equal(vm.FileID))
				})

				By("including the meta file", func() {
					Expect(assets[1].DisplayName).To(Equal("deploy.sh"))
					Expect(assets[1].Filename).To(Equal("deploy.sh"))
					Expect(assets[1].Version).To(Equal("0.0.1"))
					Expect(assets[1].Type).To(Equal("MetaFile"))
					Expect(strconv.FormatInt(assets[1].Size, 10)).To(Equal("123"))
					Expect(strconv.FormatInt(assets[1].Downloads, 10)).To(Equal("25"))
					Expect(assets[0].Downloadable).To(BeTrue())

					Expect(assets[1].DownloadRequestPayload.ProductId).To(Equal(product.ProductId))
					Expect(assets[1].DownloadRequestPayload.AppVersion).To(Equal("1"))
					Expect(assets[1].DownloadRequestPayload.MetaFileID).To(Equal(metafile.ID))
					Expect(assets[1].DownloadRequestPayload.MetaFileObjectID).To(Equal(metafile.Objects[0].FileID))
				})
			})
		})
	})
})
