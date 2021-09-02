// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Product", func() {
	var (
		httpClient  *pkgfakes.FakeHTTPClient
		marketplace *pkg.Marketplace
	)

	BeforeEach(func() {
		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			Client: httpClient,
		}
	})

	Describe("ListProduct", func() {
		BeforeEach(func() {
			product1 := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product1, "1.1.1")
			product2 := test.CreateFakeProduct(
				"",
				"My Other Product",
				"my-other-product",
				"PENDING")
			test.AddVerions(product2, "2.2.2")

			response := &pkg.ListProductResponse{
				Response: &pkg.ListProductResponsePayload{
					Products: []*models.Product{
						product1,
						product2,
					},
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
		})

		It("gets the list of products", func() {
			products, err := marketplace.ListProducts(false, "")
			Expect(err).ToNot(HaveOccurred())

			By("sending the right request", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/products"))
				Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":1,\"pageSize\":20}"))
			})

			Expect(products).To(HaveLen(2))
			Expect(products[0].Slug).To(Equal("my-super-product"))
			Expect(products[1].Slug).To(Equal("my-other-product"))
		})

		Context("with search term", func() {
			It("sends the request with the search term", func() {
				_, err := marketplace.ListProducts(false, "tanzu")
				Expect(err).ToNot(HaveOccurred())

				By("including the search term", func() {
					Expect(httpClient.DoCallCount()).To(Equal(1))
					request := httpClient.DoArgsForCall(0)
					Expect(request.URL.Query().Get("search")).To(Equal("tanzu"))
				})
			})
		})

		Context("Multiple pages of results", func() {
			BeforeEach(func() {
				var products []*models.Product
				for i := 0; i < 30; i++ {
					product := test.CreateFakeProduct(
						"",
						fmt.Sprintf("My Super Product %d", i),
						fmt.Sprintf("my-super-product-%d", i),
						"PENDING")
					test.AddVerions(product, "1.0.0")
					products = append(products, product)
				}

				response1 := &pkg.ListProductResponse{
					Response: &pkg.ListProductResponsePayload{
						Products:   products[:20],
						StatusCode: http.StatusOK,
						Params: struct {
							ProductCount int                  `json:"itemsnumber"`
							Pagination   *internal.Pagination `json:"pagination"`
						}{
							ProductCount: len(products),
							Pagination: &internal.Pagination{
								Enabled:  true,
								Page:     1,
								PageSize: 20,
							},
						},
						Message: "testing",
					},
				}
				response2 := &pkg.ListProductResponse{
					Response: &pkg.ListProductResponsePayload{
						Products:   products[20:],
						StatusCode: http.StatusOK,
						Params: struct {
							ProductCount int                  `json:"itemsnumber"`
							Pagination   *internal.Pagination `json:"pagination"`
						}{
							ProductCount: len(products),
							Pagination: &internal.Pagination{
								Enabled:  true,
								Page:     1,
								PageSize: 20,
							},
						},
						Message: "testing",
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
			})

			It("returns all results", func() {
				products, err := marketplace.ListProducts(false, "")
				Expect(err).ToNot(HaveOccurred())

				By("sending the correct requests", func() {
					Expect(httpClient.DoCallCount()).To(Equal(2))
					request := httpClient.DoArgsForCall(0)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.Path).To(Equal("/api/v1/products"))
					Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":1,\"pageSize\":20}"))

					request = httpClient.DoArgsForCall(1)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.Path).To(Equal("/api/v1/products"))
					Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":2,\"pageSize\":20}"))
				})

				By("returning all products", func() {
					Expect(products).To(HaveLen(30))
				})
			})
		})

		Context("Error fetching products", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				_, err := marketplace.ListProducts(false, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for the list of products failed: marketplace request failed: request failed"))
			})
		})

		Context("Unexpected status code", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Status:     http.StatusText(http.StatusTeapot),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.ListProducts(false, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting the list of products failed: (418) I'm a teapot"))
			})
		})

		Context("Un-parsable response", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader("This totally isn't a valid response")),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.ListProducts(false, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the list of products: invalid character 'T' looking for beginning of value"))
			})
		})
	})

	Describe("GetProduct", func() {
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "0.1.2", "1.2.3")
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
		})

		It("gets the product", func() {
			product, err := marketplace.GetProduct("my-super-product")
			Expect(err).ToNot(HaveOccurred())

			By("sending the correct request", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/products/my-super-product"))
				Expect(request.URL.Query().Get("increaseViewCount")).To(Equal("false"))
				Expect(request.URL.Query().Get("isSlug")).To(Equal("true"))
			})

			Expect(product.Slug).To(Equal("my-super-product"))
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says there are no products", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: marketplace request failed: request failed"))
			})
		})

		Context("Unexpected status code", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Status:     http.StatusText(http.StatusTeapot),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting product \"my-super-product\" failed: (418)"))
			})
		})

		Context("Un-parsable response", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader("This totally isn't a valid response")),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the response for product \"my-super-product\": invalid character 'T' looking for beginning of value"))
			})
		})
	})
})
