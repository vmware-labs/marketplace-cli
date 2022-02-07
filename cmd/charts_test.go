// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output/outputfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

func createChart(name string) (string, string) {
	chartDir, err := os.MkdirTemp("", "mkpcli-test-chart")
	Expect(err).ToNot(HaveOccurred())

	chartFile, err := chartutil.Create(name, chartDir)
	Expect(err).ToNot(HaveOccurred())

	testChart, err := loader.Load(chartFile)
	Expect(err).ToNot(HaveOccurred())

	chartPath, err := chartutil.Save(testChart, chartDir)
	Expect(err).ToNot(HaveOccurred())

	return chartPath, chartDir
}

var _ = Describe("Charts", func() {
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

	Describe("ListChartsCmd", func() {
		BeforeEach(func() {
			product = test.CreateFakeProduct("", "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(product, "1.2.3")
			product.ChartVersions = []*models.ChartVersion{
				{
					Id:         uuid.New().String(),
					Version:    "5.0.0",
					AppVersion: "1.2.3",
					Repo: &models.Repo{
						Name: "Bitnami charts repo @ Github",
						Url:  "https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus",
					},
					HelmTarUrl: "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
					TarUrl:     "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
				},
			}
		})

		It("outputs the charts", func() {
			cmd.ChartProductSlug = "my-super-product"
			cmd.ChartProductVersion = "1.2.3"
			err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderChartsCallCount()).To(Equal(1))
				charts := output.RenderChartsArgsForCall(0)
				Expect(charts).To(HaveLen(1))
				Expect(charts[0].Version).To(Equal("5.0.0"))
				Expect(charts[0].AppVersion).To(Equal("1.2.3"))
				Expect(charts[0].Repo.Name).To(Equal("Bitnami charts repo @ Github"))
				Expect(charts[0].Repo.Url).To(Equal("https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus"))
				Expect(charts[0].HelmTarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
				Expect(charts[0].TarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})

	Describe("GetChartCmd", func() {
		var chartId string
		BeforeEach(func() {
			product = test.CreateFakeProduct("", "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(product, "1.2.3")
			chartId = uuid.New().String()
			product.ChartVersions = []*models.ChartVersion{
				{
					Id:         chartId,
					Version:    "5.0.0",
					AppVersion: "1.2.3",
					Repo: &models.Repo{
						Name: "Bitnami charts repo @ Github",
						Url:  "https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus",
					},
					HelmTarUrl: "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
					TarUrl:     "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
				},
			}
		})

		It("outputs the chart", func() {
			cmd.ChartProductSlug = "my-super-product"
			cmd.ChartProductVersion = "1.2.3"
			cmd.ChartID = chartId
			err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderChartCallCount()).To(Equal(1))
				chart := output.RenderChartArgsForCall(0)
				Expect(chart.Id).To(Equal(chartId))
				Expect(chart.AppVersion).To(Equal("1.2.3"))
				Expect(chart.Version).To(Equal("5.0.0"))
				Expect(chart.Repo.Name).To(Equal("Bitnami charts repo @ Github"))
				Expect(chart.Repo.Url).To(Equal("https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus"))
				Expect(chart.HelmTarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
				Expect(chart.TarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartID = chartId
				err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})

		Context("No matching chart found", func() {
			It("says that the version does not exist", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartID = "this-chart-id-does-not-exist"
				err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("my-super-product 1.2.3 does not have the chart \"this-chart-id-does-not-exist\""))
			})
		})
	})

	Describe("AttachChartCmd", func() {
		var (
			existingChart *models.ChartVersion
			productID     string
		)
		BeforeEach(func() {
			existingChart = &models.ChartVersion{
				Id:         uuid.New().String(),
				HelmTarUrl: "https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz",
				Version:    "1.0.0",
				AppVersion: "1.2.3",
				Repo: &models.Repo{
					Name: "my-database",
				},
				IsExternalUrl: true,
			}

			productID = uuid.New().String()
			product = test.CreateFakeProduct(productID, "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(product, "1.2.3")
			product.ChartVersions = []*models.ChartVersion{existingChart}
		})

		Context("chart in public URL", func() {
			BeforeEach(func() {
				newChart := &models.ChartVersion{
					Id:         uuid.New().String(),
					HelmTarUrl: "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz",
					Version:    "0.1.0",
					Repo: &models.Repo{
						Name: "mydatabase",
					},
					IsExternalUrl: true,
				}
				marketplace.DownloadChartReturns(newChart, nil)

				updatedProduct := test.CreateFakeProduct(
					productID,
					"My Super Product",
					"my-super-product",
					"PENDING")
				test.AddVerions(updatedProduct, "1.2.3")
				updatedProduct.ChartVersions = []*models.ChartVersion{existingChart, newChart}
				marketplace.PutProductReturns(updatedProduct, nil)
			})

			It("outputs the new chart", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"

				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("getting the product details", func() {
					Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				})

				By("downloading the public chart", func() {
					Expect(marketplace.DownloadChartCallCount()).To(Equal(1))
					chartUrl := marketplace.DownloadChartArgsForCall(0)
					Expect(chartUrl.String()).To(Equal("https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"))
				})

				By("updating the product with the new chart object", func() {
					Expect(marketplace.PutProductCallCount()).To(Equal(1))
					updatedProduct, versionUpdate := marketplace.PutProductArgsForCall(0)
					Expect(versionUpdate).To(BeFalse())

					Expect(updatedProduct.DeploymentTypes).To(ContainElement("HELM"))
					Expect(updatedProduct.ChartVersions).To(HaveLen(1))
					newChart := updatedProduct.ChartVersions[0]
					Expect(newChart.AppVersion).To(Equal("1.2.3"))
					Expect(newChart.Version).To(Equal("0.1.0"))
					Expect(newChart.HelmTarUrl).To(Equal("https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"))
					Expect(newChart.Repo.Name).To(Equal("mydatabase"))
				})

				By("outputting the response", func() {
					Expect(output.RenderChartsCallCount()).To(Equal(1))
					charts := output.RenderChartsArgsForCall(0)
					Expect(charts).To(HaveLen(2))
				})
			})

			Context("Error putting product", func() {
				BeforeEach(func() {
					marketplace.PutProductReturns(nil, fmt.Errorf("put product failed"))
				})

				It("returns an error", func() {
					cmd.ChartProductSlug = "my-super-product"
					cmd.ChartProductVersion = "1.2.3"
					cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("put product failed"))
				})
			})
		})

		Context("chart is local file", func() {
			var (
				chartDir  string
				chartPath string
				uploader  *internalfakes.FakeUploader
			)

			BeforeEach(func() {
				chartPath, chartDir = createChart("mydatabase")

				newChart := &models.ChartVersion{
					Id:         uuid.New().String(),
					HelmTarUrl: "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz",
					Version:    "0.1.0",
					AppVersion: "1.2.3",
					Repo: &models.Repo{
						Name: "mydatabase",
					},
					IsExternalUrl: true,
				}

				updatedProduct := test.CreateFakeProduct(productID, "My Super Product", "my-super-product", "PENDING")
				test.AddVerions(updatedProduct, "1.2.3")
				updatedProduct.ChartVersions = []*models.ChartVersion{existingChart, newChart}
				marketplace.PutProductReturns(updatedProduct, nil)
				marketplace.GetUploadCredentialsReturns(&pkg.CredentialsResponse{}, nil)
				uploader = &internalfakes.FakeUploader{}
				uploader.UploadProductFileReturns("mydatabase-0.1.0.tgz", "https://marketplace.example.vmware.com/uploader/mydatabase-0.1.0.tgz", nil)
				marketplace.GetUploaderReturns(uploader)
			})

			AfterEach(func() {
				err := os.RemoveAll(chartDir)
				Expect(err).ToNot(HaveOccurred())
			})

			It("sends the right requests", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartURL = chartPath

				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				By("getting the product details", func() {
					Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
				})

				By("uploadeding the chart", func() {
					Expect(uploader.UploadProductFileCallCount()).To(Equal(1))
					uploadedChartURL := uploader.UploadProductFileArgsForCall(0)
					Expect(uploadedChartURL).To(Equal(chartPath))
				})

				By("updating the product with the new chart object", func() {
					Expect(marketplace.PutProductCallCount()).To(Equal(1))
					updatedProduct, versionUpdate := marketplace.PutProductArgsForCall(0)
					Expect(versionUpdate).To(BeFalse())

					Expect(updatedProduct.DeploymentTypes).To(ContainElement("HELM"))
					Expect(updatedProduct.ChartVersions).To(HaveLen(1))
					newChart := updatedProduct.ChartVersions[0]
					Expect(newChart.AppVersion).To(Equal("1.2.3"))
					Expect(newChart.Version).To(Equal("0.1.0"))
					Expect(newChart.HelmTarUrl).To(Equal("https://marketplace.example.vmware.com/uploader/mydatabase-0.1.0.tgz"))
					Expect(newChart.Repo.Name).To(Equal("mydatabase"))
				})

				By("outputting the response", func() {
					Expect(output.RenderChartsCallCount()).To(Equal(1))
					charts := output.RenderChartsArgsForCall(0)
					Expect(charts).To(HaveLen(2))
				})
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
			})

			It("returns an error", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"
				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})
})
