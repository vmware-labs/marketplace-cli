// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
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
		stdout     *Buffer
		stderr     *Buffer
		httpClient *pkgfakes.FakeHTTPClient
		output     *outputfakes.FakeFormat
		uploader   *internalfakes.FakeUploader
	)

	BeforeEach(func() {
		httpClient = &pkgfakes.FakeHTTPClient{}
		output = &outputfakes.FakeFormat{}
		uploader = &internalfakes.FakeUploader{}
		cmd.Output = output
		cmd.Marketplace = &pkg.Marketplace{
			Client:   httpClient,
			Uploader: uploader,
		}
		stdout = NewBuffer()
		stderr = NewBuffer()
	})

	Describe("ListChartsCmd", func() {
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")
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
			response := &pkg.GetProductResponse{
				Response: &pkg.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			responseBytes, err := json.Marshal(response)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil)

			cmd.ListChartsCmd.SetOut(stdout)
			cmd.ListChartsCmd.SetErr(stderr)
		})

		It("outputs the charts", func() {
			cmd.ChartProductSlug = "my-super-product"
			cmd.ChartProductVersion = "1.2.3"
			err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("sending the correct request", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/products/my-super-product"))
				Expect(request.URL.Query().Get("increaseViewCount")).To(Equal("false"))
				Expect(request.URL.Query().Get("isSlug")).To(Equal("true"))
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

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says there are no products", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("Error fetching products", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(0, nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: marketplace request failed: request failed"))
			})
		})

		Context("No product version found", func() {
			It("says that the version does not exist", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "9.9.9"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
			})
		})
	})

	Describe("GetChartCmd", func() {
		var chartId string
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")
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
			response := &pkg.GetProductResponse{
				Response: &pkg.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			responseBytes, err := json.Marshal(response)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil)

			cmd.GetChartCmd.SetOut(stdout)
			cmd.GetChartCmd.SetErr(stderr)
		})

		It("sends the right request", func() {
			cmd.ChartProductSlug = "my-super-product"
			cmd.ChartProductVersion = "1.2.3"
			cmd.ChartID = chartId
			err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("sending the correct request", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/products/my-super-product"))
				Expect(request.URL.Query().Get("increaseViewCount")).To(Equal("false"))
				Expect(request.URL.Query().Get("isSlug")).To(Equal("true"))
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

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says that the product was not found", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartID = chartId
				err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No version found", func() {
			It("says that the version does not exist", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "9.9.9"
				cmd.ChartID = chartId
				err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
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

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartID = chartId
				err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: marketplace request failed: request failed"))
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
			product := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3")
			product.ChartVersions = []*models.ChartVersion{existingChart}
			response := &pkg.GetProductResponse{
				Response: &pkg.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}
			responseBytes, err := json.Marshal(response)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil)

			cmd.AttachChartCmd.SetOut(stdout)
			cmd.AttachChartCmd.SetErr(stderr)
		})

		Context("chart in public URL", func() {
			var chartDir string

			BeforeEach(func() {
				var chartPath string
				chartPath, chartDir = createChart("mydatabase")
				chartData, err := ioutil.ReadFile(chartPath)
				Expect(err).ToNot(HaveOccurred())
				httpClient.DoReturnsOnCall(1, &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(chartData)),
				}, nil)

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

				product := test.CreateFakeProduct(
					productID,
					"My Super Product",
					"my-super-product",
					"PENDING")
				test.AddVerions(product, "1.2.3")
				product.ChartVersions = []*models.ChartVersion{existingChart, newChart}
				response := &pkg.GetProductResponse{
					Response: &pkg.GetProductResponsePayload{
						Data:       product,
						StatusCode: http.StatusOK,
						Message:    "testing",
					},
				}
				responseBytes, err := json.Marshal(response)
				Expect(err).ToNot(HaveOccurred())

				httpClient.DoReturnsOnCall(2, &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
				}, nil)
			})

			AfterEach(func() {
				err := os.RemoveAll(chartDir)
				Expect(err).ToNot(HaveOccurred())
			})

			It("sends the right requests", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"

				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())

				Expect(httpClient.DoCallCount()).To(Equal(3))
				By("first, getting the product", func() {
					request := httpClient.DoArgsForCall(0)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.Path).To(Equal("/api/v1/products/my-super-product"))
					Expect(request.URL.Query().Get("increaseViewCount")).To(Equal("false"))
					Expect(request.URL.Query().Get("isSlug")).To(Equal("true"))
				})

				By("second, downloading the chart", func() {
					request := httpClient.DoArgsForCall(1)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.String()).To(Equal("https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"))
				})

				By("third, attaching the chart", func() {
					request := httpClient.DoArgsForCall(2)
					Expect(request.Method).To(Equal("PUT"))
					Expect(request.URL.Path).To(Equal(fmt.Sprintf("/api/v1/products/%s", productID)))
					updatedProduct := &models.Product{}
					requestBody, err := ioutil.ReadAll(request.Body)
					Expect(err).ToNot(HaveOccurred())
					err = json.Unmarshal(requestBody, updatedProduct)
					Expect(err).ToNot(HaveOccurred())

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
					httpClient.DoReturnsOnCall(2,
						&http.Response{
							StatusCode: http.StatusTeapot,
							Body:       ioutil.NopCloser(strings.NewReader("Teapots all the way down")),
						}, nil)
				})
				It("prints the error", func() {
					cmd.ChartProductSlug = "my-super-product"
					cmd.ChartProductVersion = "1.2.3"
					cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("updating product \"my-super-product\" failed: (418)\nTeapots all the way down"))
				})
			})

			Context("No permission to update product", func() {
				BeforeEach(func() {
					httpClient.DoReturnsOnCall(2,
						&http.Response{
							StatusCode: http.StatusForbidden,
							Body:       ioutil.NopCloser(strings.NewReader("{\"response\":{\"message\":\"User is not authorized to perform this action\"}}\n")),
						}, nil)
				})
				It("prints the error", func() {
					cmd.ChartProductSlug = "my-super-product"
					cmd.ChartProductVersion = "1.2.3"
					cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"
					err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("you do not have permission to modify the product \"my-super-product\""))
				})
			})
		})

		Context("chart is local file", func() {
			var (
				chartDir  string
				chartPath string
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

				product := test.CreateFakeProduct(
					productID,
					"My Super Product",
					"my-super-product",
					"PENDING")
				test.AddVerions(product, "1.2.3")
				product.ChartVersions = []*models.ChartVersion{existingChart, newChart}

				uploader.UploadReturns("https://www.example.com/storage/mydatabase-0.1.0.tgz", nil)
				httpClient.DoReturnsOnCall(1, ResponseWithPayload(&pkg.CredentialsResponse{
					AccessID:     "access_id",
					AccessKey:    "access_key",
					SessionToken: "session_token",
					Expiration:   time.Now(),
				}), nil)

				httpClient.DoReturnsOnCall(2, ResponseWithPayload(&pkg.GetProductResponse{
					Response: &pkg.GetProductResponsePayload{
						Data:       product,
						StatusCode: http.StatusOK,
						Message:    "testing",
					},
				}), nil)
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

				Expect(httpClient.DoCallCount()).To(Equal(3))
				By("first, getting the product", func() {
					request := httpClient.DoArgsForCall(0)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.Path).To(Equal("/api/v1/products/my-super-product"))
					Expect(request.URL.Query().Get("increaseViewCount")).To(Equal("false"))
					Expect(request.URL.Query().Get("isSlug")).To(Equal("true"))
				})

				By("second, uploading the chart", func() {
					request := httpClient.DoArgsForCall(1)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.Path).To(Equal("/aws/credentials/generate"))

					Expect(uploader.UploadCallCount()).To(Equal(1))
					uploadedFile := uploader.UploadArgsForCall(0)
					Expect(uploadedFile).To(Equal(chartPath))
				})

				By("third, attaching the chart", func() {
					request := httpClient.DoArgsForCall(2)
					Expect(request.Method).To(Equal("PUT"))
					Expect(request.URL.Path).To(Equal(fmt.Sprintf("/api/v1/products/%s", productID)))
					updatedProduct := &models.Product{}
					requestBody, err := ioutil.ReadAll(request.Body)
					Expect(err).ToNot(HaveOccurred())
					err = json.Unmarshal(requestBody, updatedProduct)
					Expect(err).ToNot(HaveOccurred())

					Expect(updatedProduct.DeploymentTypes).To(ContainElement("HELM"))
					Expect(updatedProduct.ChartVersions).To(HaveLen(1))
					newChart := updatedProduct.ChartVersions[0]
					Expect(newChart.AppVersion).To(Equal("1.2.3"))
					Expect(newChart.Version).To(Equal("0.1.0"))
					Expect(newChart.HelmTarUrl).To(Equal("https://www.example.com/storage/mydatabase-0.1.0.tgz"))
					Expect(newChart.Repo.Name).To(Equal("mydatabase"))
				})

				By("outputting the response", func() {
					Expect(output.RenderChartsCallCount()).To(Equal(1))
					charts := output.RenderChartsArgsForCall(0)
					Expect(charts).To(HaveLen(2))
				})
			})
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(0,
					&http.Response{
						StatusCode: http.StatusNotFound,
					}, nil)
			})

			It("says that the product was not found", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"
				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No product version found", func() {
			It("says there are no versions", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "0.0.0"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-0.1.0.tgz"
				err := cmd.AttachChartCmd.RunE(cmd.AttachChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 0.0.0"))
			})
		})
	})
})
