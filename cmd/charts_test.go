// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/lib"
	"github.com/vmware-labs/marketplace-cli/v2/lib/libfakes"
	"github.com/vmware-labs/marketplace-cli/v2/models"
)

var _ = Describe("Charts", func() {
	var (
		stdout *Buffer
		stderr *Buffer

		originalHttpClient lib.HTTPClient
		httpClient         *libfakes.FakeHTTPClient
	)

	BeforeEach(func() {
		stdout = NewBuffer()
		stderr = NewBuffer()

		originalHttpClient = lib.Client
		httpClient = &libfakes.FakeHTTPClient{}
		lib.Client = httpClient
	})

	AfterEach(func() {
		lib.Client = originalHttpClient
	})

	Describe("ListChartsCmd", func() {
		BeforeEach(func() {
			product := CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			AddVerions(product, "1.2.3", "2.3.4")
			product.ChartVersions = []*models.ChartVersion{
				{
					Id:         "my-chart",
					Version:    "1.0.0",
					AppVersion: "1.2.3",
					Repo: &models.Repo{
						Name: "nitbami",
						Url:  "https://charts.nitbami.com/nitbami",
					},
				},
			}
			response := &cmd.GetProductResponse{
				Response: &cmd.GetProductResponsePayload{
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
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "1.2.3"
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
				Expect(stdout).To(Say("ID        VERSION  URL  REPOSITORY"))
				Expect(stdout).To(Say("my-chart  1.0.0         nitbami https://charts.nitbami.com/nitbami"))
			})
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says there are no products", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
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
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: request failed"))
			})
		})

		Context("No product version found", func() {
			It("says that the version does not exist", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
			})
		})

		Context("No chart", func() {
			It("says there are no charts", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "2.3.4"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())
				Expect(stdout).To(Say("product \"my-super-product\" 2.3.4 does not have any charts"))
			})
		})
	})

	Describe("CreateChartCmd", func() {
		var productId string
		BeforeEach(func() {
			chart1 := &models.ChartVersion{
				Id:         "mydatabase",
				TarUrl:     "https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz",
				Version:    "1.0.0",
				AppVersion: "1.2.3",
				Repo: &models.Repo{
					Name: "nitbami",
					Url:  "https://charts.nitbami.com/nitbami",
				},
			}
			chart2 := &models.ChartVersion{
				Id:         "mydatabase",
				TarUrl:     "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz",
				Version:    "2.0.0",
				AppVersion: "1.2.3",
				Repo: &models.Repo{
					Name: "nitbami",
					Url:  "https://charts.nitbami.com/nitbami",
				},
			}

			productId = uuid.New().String()
			product := CreateFakeProduct(
				productId,
				"My Super Product",
				"my-super-product",
				"PENDING")
			AddVerions(product, "1.2.3")
			product.ChartVersions = []*models.ChartVersion{chart1}
			response1 := &cmd.GetProductResponse{
				Response: &cmd.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			updatedProduct := CreateFakeProduct(
				productId,
				"My Super Product",
				"my-super-product",
				"PENDING")
			AddVerions(updatedProduct, "1.2.3")
			updatedProduct.ChartVersions = []*models.ChartVersion{chart1, chart2}
			response2 := &cmd.GetProductResponse{
				Response: &cmd.GetProductResponsePayload{
					Data:       updatedProduct,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			responseBytes, err := json.Marshal(response1)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil)

			responseBytes, err = json.Marshal(response2)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturnsOnCall(1, &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil)

			cmd.CreateChartCmd.SetOut(stdout)
			cmd.CreateChartCmd.SetErr(stderr)
		})

		It("sends the right requests", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "1.2.3"
			cmd.ChartName = "mydatabase"
			cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
			cmd.ChartVersion = "2.0.0"
			cmd.ChartRepositoryName = "nitbami"
			cmd.ChartRepositoryURL = "https://charts.nitbami.com/nitbami"

			err := cmd.CreateChartCmd.RunE(cmd.CreateChartCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			Expect(httpClient.DoCallCount()).To(Equal(2))
			By("first, getting the product", func() {
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/products/my-super-product"))
				Expect(request.URL.Query().Get("increaseViewCount")).To(Equal("false"))
				Expect(request.URL.Query().Get("isSlug")).To(Equal("true"))
			})

			By("second, sending the new product", func() {
				request := httpClient.DoArgsForCall(1)
				Expect(request.Method).To(Equal("PUT"))
				Expect(request.URL.Path).To(Equal(fmt.Sprintf("/api/v1/products/%s", productId)))
				updatedProduct := &models.Product{}
				requestBody, err := ioutil.ReadAll(request.Body)
				Expect(err).ToNot(HaveOccurred())
				err = json.Unmarshal(requestBody, updatedProduct)
				Expect(err).ToNot(HaveOccurred())

				Expect(updatedProduct.DeploymentTypes).To(ContainElement("HELM"))
				Expect(updatedProduct.ChartVersions).To(ContainElements(
					&models.ChartVersion{
						Id:             "mydatabase",
						TarUrl:         "https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz",
						Version:        "1.0.0",
						AppVersion:     "1.2.3",
						InstallOptions: "",
						Repo: &models.Repo{
							Name: "nitbami",
							Url:  "https://charts.nitbami.com/nitbami",
						},
					},
					&models.ChartVersion{
						Id:             "mydatabase",
						TarUrl:         "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz",
						Version:        "2.0.0",
						AppVersion:     "1.2.3",
						InstallOptions: "",
						Repo: &models.Repo{
							Name: "nitbami",
							Url:  "https://charts.nitbami.com/nitbami",
						},
					},
				))
			})

			By("outputting the response", func() {
				Expect(stdout).To(Say("ID          VERSION  URL                                                             REPOSITORY"))
				Expect(stdout).To(Say("mydatabase  1.0.0    https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz  nitbami https://charts.nitbami.com/nitbami"))
				Expect(stdout).To(Say("mydatabase  2.0.0    https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz  nitbami https://charts.nitbami.com/nitbami"))
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
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ChartName = "mydatabase"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
				cmd.ChartVersion = "2.0.0"
				cmd.ChartRepositoryName = "nitbami"
				cmd.ChartRepositoryURL = "https://charts.nitbami.com/nitbami"
				err := cmd.CreateChartCmd.RunE(cmd.CreateChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No product version found", func() {
			It("says there are no versions", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "0.0.0"
				cmd.ChartName = "mydatabase"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
				cmd.ChartVersion = "2.0.0"
				cmd.ChartRepositoryName = "nitbami"
				cmd.ChartRepositoryURL = "https://charts.nitbami.com/nitbami"
				err := cmd.CreateChartCmd.RunE(cmd.CreateChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 0.0.0, please add it first"))
			})
		})

		Context("Error putting product", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(1,
					&http.Response{
						StatusCode: http.StatusTeapot,
						Body:       ioutil.NopCloser(strings.NewReader("Teapots all the way down")),
					}, nil)
			})
			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ChartName = "mydatabase"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
				cmd.ChartVersion = "2.0.0"
				cmd.ChartRepositoryName = "nitbami"
				cmd.ChartRepositoryURL = "https://charts.nitbami.com/nitbami"
				err := cmd.CreateChartCmd.RunE(cmd.CreateChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("updating product \"my-super-product\" failed: (418)\nTeapots all the way down"))
			})
		})
	})
})
