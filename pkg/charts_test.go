// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
	helmChart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

var _ = Describe("Charts", func() {
	var (
		chart               *helmChart.Chart
		chartDir            string
		chartPath           string
		chartLoader         *pkgfakes.FakeChartLoaderFunc
		previousChartLoader pkg.ChartLoaderFunc
		httpClient          *pkgfakes.FakeHTTPClient
		marketplace         *pkg.Marketplace
	)
	BeforeEach(func() {
		chart, chartPath, chartDir = test.CreateFakeChart("hyperspace-database-chart")

		chartLoader = &pkgfakes.FakeChartLoaderFunc{}
		chartLoader.Returns(chart, nil)
		previousChartLoader = pkg.ChartLoader
		pkg.ChartLoader = chartLoader.Spy

		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			Client: httpClient,
			Host:   "marketplace.vmware.example",
		}
	})
	AfterEach(func() {
		Expect(os.RemoveAll(chartDir)).To(Succeed())
		pkg.ChartLoader = previousChartLoader
	})

	Describe("LoadChart", func() {
		It("loads a chart from a path", func() {
			chart, err := pkg.LoadChart("path/to/local/chart.tgz")
			Expect(err).ToNot(HaveOccurred())

			Expect(chart.Repo.Name).To(Equal("hyperspace-database-chart"))

			Expect(chartLoader.CallCount()).To(Equal(1))
			Expect(chartLoader.ArgsForCall(0)).To(Equal("path/to/local/chart.tgz"))
		})

		When("loading the chart failed", func() {
			BeforeEach(func() {
				chartLoader.Returns(nil, errors.New("chart loader failed"))
			})
			It("returns an error", func() {
				_, err := pkg.LoadChart("path/to/local/chart.tgz")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to read chart: chart loader failed"))
			})
		})
	})

	Describe("DownloadChart", func() {
		BeforeEach(func() {
			// This does not need to be cleaned up because chartDir will be removed
			chartArchive, err := chartutil.Save(chart, chartDir)
			Expect(err).ToNot(HaveOccurred())
			chartBytes, err := ioutil.ReadFile(chartArchive)
			Expect(err).ToNot(HaveOccurred())
			httpClient.DoReturns(test.MakeBytesResponse(chartBytes), nil)
		})

		It("downloads a chart", func() {
			chartUrl, err := url.Parse("https://charts.example.com/my-chart.tgz")
			Expect(err).ToNot(HaveOccurred())
			chart, err := marketplace.DownloadChart(chartUrl)
			Expect(err).ToNot(HaveOccurred())

			By("requesting the chart", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				Expect(httpClient.DoArgsForCall(0).URL.String()).To(Equal("https://charts.example.com/my-chart.tgz"))
			})

			Expect(chart.Repo.Name).To(Equal("hyperspace-database-chart"))
		})

		When("downloading the chart fails", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("http request failed"))
			})
			It("returns an error", func() {
				chartUrl, err := url.Parse("https://charts.example.com/my-chart.tgz")
				Expect(err).ToNot(HaveOccurred())
				_, err = marketplace.DownloadChart(chartUrl)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to download chart: http request failed"))
			})
		})

		When("copying the chart fails", func() {
			BeforeEach(func() {
				httpClient.DoReturns(test.MakeFailingBodyResponse("read fail"), nil)
			})
			It("returns an error", func() {
				chartUrl, err := url.Parse("https://charts.example.com/my-chart.tgz")
				Expect(err).ToNot(HaveOccurred())
				_, err = marketplace.DownloadChart(chartUrl)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to save local chart: read fail"))
			})
		})

		When("loading the chart fails", func() {
			BeforeEach(func() {
				chartLoader.Returns(nil, errors.New("chart loading failed"))
			})
			It("returns an error", func() {
				chartUrl, err := url.Parse("https://charts.example.com/my-chart.tgz")
				Expect(err).ToNot(HaveOccurred())
				_, err = marketplace.DownloadChart(chartUrl)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to read chart: chart loading failed"))
			})
		})
	})

	Describe("AttachLocalChart", func() {
		var uploader *internalfakes.FakeUploader
		BeforeEach(func() {
			httpClient.PutStub = PutProductEchoResponse
			uploader = &internalfakes.FakeUploader{}
			uploader.UploadProductFileReturns("", "https://example.com/uploaded-chart.tgz", nil)
			marketplace.SetUploader(uploader)
		})

		It("uploads and attaches a local chart", func() {
			product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
			version := &models.Version{Number: "1.2.3"}
			test.AddVerions(product, "1.2.3")
			updatedProduct, err := marketplace.AttachLocalChart(chartPath, "helm install it", product, version)
			Expect(err).ToNot(HaveOccurred())

			By("loading the local chart", func() {
				Expect(chartLoader.CallCount()).To(Equal(1))
				Expect(chartLoader.ArgsForCall(0)).To(Equal(chartPath))
			})

			By("uploading the chart", func() {
				Expect(uploader.UploadProductFileCallCount()).To(Equal(1))
				Expect(uploader.UploadProductFileArgsForCall(0)).To(Equal(chartPath))
			})

			By("updating the product", func() {
				Expect(httpClient.PutCallCount()).To(Equal(1))
			})

			By("returning the updated product", func() {
				Expect(updatedProduct.ChartVersions).To(HaveLen(1))
				attachedChart := updatedProduct.ChartVersions[0]
				Expect(attachedChart.HelmTarUrl).To(Equal("https://example.com/uploaded-chart.tgz"))
				Expect(attachedChart.AppVersion).To(Equal("1.2.3"))
				Expect(attachedChart.Readme).To(Equal("helm install it"))
			})
		})

		When("loading the chart fails", func() {
			BeforeEach(func() {
				chartLoader.Returns(nil, errors.New("load chart failed"))
			})

			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
				version := &models.Version{Number: "1.2.3"}
				test.AddVerions(product, "1.2.3")
				_, err := marketplace.AttachLocalChart(chartPath, "helm install it", product, version)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to read chart: load chart failed"))
			})
		})

		When("getting the uploader fails", func() {
			BeforeEach(func() {
				marketplace.SetUploader(nil)
				httpClient.GetReturns(nil, errors.New("get uploader failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
				version := &models.Version{Number: "1.2.3"}
				test.AddVerions(product, "1.2.3")
				_, err := marketplace.AttachLocalChart(chartPath, "helm install it", product, version)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get upload credentials: get uploader failed"))
			})
		})

		When("uploading the chart fails", func() {
			BeforeEach(func() {
				uploader.UploadProductFileReturns("", "", errors.New("upload product file failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
				version := &models.Version{Number: "1.2.3"}
				test.AddVerions(product, "1.2.3")
				_, err := marketplace.AttachLocalChart(chartPath, "helm install it", product, version)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("upload product file failed"))
			})
		})

		When("updating the product fails", func() {
			BeforeEach(func() {
				httpClient.PutReturns(nil, errors.New("update product failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
				version := &models.Version{Number: "1.2.3"}
				test.AddVerions(product, "1.2.3")
				_, err := marketplace.AttachLocalChart(chartPath, "helm install it", product, version)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the update for product \"hyperspace-database\" failed: update product failed"))
			})
		})
	})

	Describe("AttachPublicChart", func() {
		var (
			chartBytes []byte
			chartUrl   *url.URL
		)
		BeforeEach(func() {
			var err error
			chartUrl, err = url.Parse("https://example.com/my-public-chart.tgz")
			Expect(err).ToNot(HaveOccurred())

			// This does not need to be cleaned up because chartDir will be removed
			chartArchive, err := chartutil.Save(chart, chartDir)
			Expect(err).ToNot(HaveOccurred())
			chartBytes, err = ioutil.ReadFile(chartArchive)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturns(&http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(chartBytes)),
			}, nil)
			httpClient.PutStub = PutProductEchoResponse
		})

		It("attaches a public chart", func() {
			product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
			version := &models.Version{Number: "1.2.3"}
			test.AddVerions(product, "1.2.3")
			updatedProduct, err := marketplace.AttachPublicChart(chartUrl, "helm install it", product, version)
			Expect(err).ToNot(HaveOccurred())

			By("downloading the chart", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				downloadRequest := httpClient.DoArgsForCall(0)
				Expect(downloadRequest.Method).To(Equal("GET"))
				Expect(downloadRequest.URL.String()).To(Equal("https://example.com/my-public-chart.tgz"))
			})

			By("updating the product", func() {
				Expect(httpClient.PutCallCount()).To(Equal(1))
				url, _, contentType := httpClient.PutArgsForCall(0)
				Expect(url.String()).To(ContainSubstring("https://marketplace.vmware.example/api/v1/products"))
				Expect(contentType).To(Equal("application/json"))
			})

			By("returning the updated product", func() {
				Expect(updatedProduct.ChartVersions).To(HaveLen(1))
				attachedChart := updatedProduct.ChartVersions[0]
				Expect(attachedChart.HelmTarUrl).To(Equal("https://example.com/my-public-chart.tgz"))
				Expect(attachedChart.AppVersion).To(Equal("1.2.3"))
				Expect(attachedChart.Readme).To(Equal("helm install it"))
			})
		})

		When("downloading the chart fails", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("download chart failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
				version := &models.Version{Number: "1.2.3"}
				test.AddVerions(product, "1.2.3")
				_, err := marketplace.AttachPublicChart(chartUrl, "helm install it", product, version)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to download chart: download chart failed"))
			})
		})

		When("updating the product fails", func() {
			BeforeEach(func() {
				httpClient.PutReturns(nil, errors.New("update product failed"))
			})
			It("returns an error", func() {
				product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
				version := &models.Version{Number: "1.2.3"}
				test.AddVerions(product, "1.2.3")
				_, err := marketplace.AttachPublicChart(chartUrl, "helm install it", product, version)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the update for product \"hyperspace-database\" failed: update product failed"))
			})
		})
	})
})
