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
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output/outputfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Charts", func() {
	var (
		stdout     *Buffer
		stderr     *Buffer
		httpClient *pkgfakes.FakeHTTPClient
		output     *outputfakes.FakeFormat
	)

	BeforeEach(func() {
		httpClient = &pkgfakes.FakeHTTPClient{}
		output = &outputfakes.FakeFormat{}
		cmd.Output = output
		cmd.Marketplace = &pkg.Marketplace{
			Client: httpClient,
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

	Describe("CreateChartCmd", func() {
		var productID string
		BeforeEach(func() {
			chart1 := &models.ChartVersion{
				Id:         uuid.New().String(),
				HelmTarUrl: "https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz",
				TarUrl:     "https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz",
				Version:    "1.0.0",
				AppVersion: "1.2.3",
				Repo: &models.Repo{
					Name: "Bitnami charts repo @ Github",
					Url:  "https://github.com/bitnami/charts/tree/master/bitnami/mydatabase",
				},
			}
			chart2 := &models.ChartVersion{
				Id:         uuid.New().String(),
				HelmTarUrl: "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz",
				TarUrl:     "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz",
				Version:    "2.0.0",
				AppVersion: "1.2.3",
				Repo: &models.Repo{
					Name: "Bitnami charts repo @ Github",
					Url:  "https://github.com/bitnami/charts/tree/master/bitnami/mydatabase",
				},
			}

			productID = uuid.New().String()
			product := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3")
			product.ChartVersions = []*models.ChartVersion{chart1}
			response1 := &pkg.GetProductResponse{
				Response: &pkg.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			updatedProduct := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(updatedProduct, "1.2.3")
			updatedProduct.ChartVersions = []*models.ChartVersion{chart1, chart2}
			response2 := &pkg.GetProductResponse{
				Response: &pkg.GetProductResponsePayload{
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
			cmd.ChartProductSlug = "my-super-product"
			cmd.ChartProductVersion = "1.2.3"
			cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
			cmd.ChartVersion = "2.0.0"
			cmd.ChartRepositoryName = "Bitnami charts repo @ Github"
			cmd.ChartRepositoryURL = "https://github.com/bitnami/charts/tree/master/bitnami/mydatabse"

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
				Expect(request.URL.Path).To(Equal(fmt.Sprintf("/api/v1/products/%s", productID)))
				updatedProduct := &models.Product{}
				requestBody, err := ioutil.ReadAll(request.Body)
				Expect(err).ToNot(HaveOccurred())
				err = json.Unmarshal(requestBody, updatedProduct)
				Expect(err).ToNot(HaveOccurred())

				Expect(updatedProduct.DeploymentTypes).To(ContainElement("HELM"))
				//Expect(updatedProduct.ChartVersions).To(ContainElements(
				//	&models.ChartVersion{
				//		HelmTarUrl: "https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz",
				//		TarUrl:     "https://charts.nitbami.com/nitbami/charts/mydatabase-1.0.0.tgz",
				//		Version:    "1.0.0",
				//		AppVersion: "1.2.3",
				//		Repo: &models.Repo{
				//			Name: "Bitnami charts repo @ Github",
				//			Url:  "https://github.com/bitnami/charts/tree/master/bitnami/mydatabase",
				//		},
				//	},
				//	&models.ChartVersion{
				//		HelmTarUrl: "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz",
				//		TarUrl:     "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz",
				//		Version:    "2.0.0",
				//		AppVersion: "1.2.3",
				//		Repo: &models.Repo{
				//			Name: "Bitnami charts repo @ Github",
				//			Url:  "https://github.com/bitnami/charts/tree/master/bitnami/mydatabase",
				//		},
				//	},
				//))
			})

			By("outputting the response", func() {
				Expect(output.RenderChartsCallCount()).To(Equal(1))
				charts := output.RenderChartsArgsForCall(0)
				Expect(charts).To(HaveLen(2))
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
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
				cmd.ChartVersion = "2.0.0"
				cmd.ChartRepositoryName = "Bitnami charts repo @ Github"
				cmd.ChartRepositoryURL = "https://github.com/bitnami/charts/tree/master/bitnami/mydatabse"
				err := cmd.CreateChartCmd.RunE(cmd.CreateChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No product version found", func() {
			It("says there are no versions", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "0.0.0"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
				cmd.ChartVersion = "2.0.0"
				cmd.ChartRepositoryName = "Bitnami charts repo @ Github"
				cmd.ChartRepositoryURL = "https://github.com/bitnami/charts/tree/master/bitnami/mydatabse"
				err := cmd.CreateChartCmd.RunE(cmd.CreateChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 0.0.0"))
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
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartURL = "https://charts.nitbami.com/nitbami/charts/mydatabase-2.0.0.tgz"
				cmd.ChartVersion = "2.0.0"
				cmd.ChartRepositoryName = "Bitnami charts repo @ Github"
				cmd.ChartRepositoryURL = "https://github.com/bitnami/charts/tree/master/bitnami/mydatabse"
				err := cmd.CreateChartCmd.RunE(cmd.CreateChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("updating product \"my-super-product\" failed: (418)\nTeapots all the way down"))
			})
		})
	})
})
