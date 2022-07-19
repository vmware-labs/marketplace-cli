// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Product", func() {
	var (
		stderr      *Buffer
		httpClient  *pkgfakes.FakeHTTPClient
		marketplace *pkg.Marketplace
	)

	BeforeEach(func() {
		stderr = NewBuffer()
		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			Client: httpClient,
			Output: stderr,
		}
	})

	Describe("ListProduct", func() {
		BeforeEach(func() {
			product1 := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				models.SolutionTypeImage)
			test.AddVerions(product1, "1.1.1")
			product2 := test.CreateFakeProduct(
				"",
				"My Other Product",
				"my-other-product",
				models.SolutionTypeImage)
			test.AddVerions(product2, "2.2.2")
			products := []*models.Product{
				product1,
				product2,
			}

			response := &pkg.ListProductResponse{
				Response: &pkg.ListProductResponsePayload{
					Products: products,
					Params: &pkg.ListProductResponseParams{
						ProductCount: len(products),
					},
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}
			httpClient.GetReturns(test.MakeJSONResponse(response), nil)
		})

		It("gets the list of products", func() {
			products, err := marketplace.ListProducts(false, "")
			Expect(err).ToNot(HaveOccurred())

			By("sending the right request", func() {
				Expect(httpClient.GetCallCount()).To(Equal(1))
				url := httpClient.GetArgsForCall(0)
				Expect(url.Path).To(Equal("/api/v1/products"))
				Expect(url.Query().Get("pagination")).To(Equal("{\"page\":1,\"pageSize\":20}"))
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
					Expect(httpClient.GetCallCount()).To(Equal(1))
					url := httpClient.GetArgsForCall(0)
					Expect(url.Query().Get("search")).To(Equal("tanzu"))
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
						models.SolutionTypeImage)
					test.AddVerions(product, "1.0.0")
					products = append(products, product)
				}

				response1 := &pkg.ListProductResponse{
					Response: &pkg.ListProductResponsePayload{
						Products:   products[:20],
						StatusCode: http.StatusOK,
						Params: &pkg.ListProductResponseParams{
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
						Params: &pkg.ListProductResponseParams{
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
				httpClient.GetReturnsOnCall(0, test.MakeJSONResponse(response1), nil)
				httpClient.GetReturnsOnCall(1, test.MakeJSONResponse(response2), nil)
			})

			It("returns all results", func() {
				products, err := marketplace.ListProducts(false, "")
				Expect(err).ToNot(HaveOccurred())

				By("sending the correct requests", func() {
					Expect(httpClient.GetCallCount()).To(Equal(2))
					url := httpClient.GetArgsForCall(0)
					Expect(url.Path).To(Equal("/api/v1/products"))
					Expect(url.Query().Get("pagination")).To(Equal("{\"page\":1,\"pageSize\":20}"))

					url = httpClient.GetArgsForCall(1)
					Expect(url.Path).To(Equal("/api/v1/products"))
					Expect(url.Query().Get("pagination")).To(Equal("{\"page\":2,\"pageSize\":20}"))
				})

				By("returning all products", func() {
					Expect(products).To(HaveLen(30))
				})
			})
		})

		Context("Error fetching products", func() {
			BeforeEach(func() {
				httpClient.GetReturns(nil, errors.New("request failed"))
			})

			It("prints the error", func() {
				_, err := marketplace.ListProducts(false, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for the list of products failed: request failed"))
			})
		})

		Context("Unexpected status code", func() {
			BeforeEach(func() {
				httpClient.GetReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Status:     http.StatusText(http.StatusTeapot),
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("Teapots!"))),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.ListProducts(false, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting the list of products failed: (418) I'm a teapot: Teapots!"))
			})
		})

		Context("Un-parsable response", func() {
			BeforeEach(func() {
				httpClient.GetReturns(test.MakeStringResponse("This totally isn't a valid response"), nil)
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
				models.SolutionTypeImage)
			test.AddVerions(product, "0.1.2", "1.2.3")
			response := &pkg.GetProductResponse{
				Response: &pkg.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			httpClient.GetReturns(test.MakeJSONResponse(response), nil)
		})

		It("gets the product", func() {
			product, err := marketplace.GetProduct("my-super-product")
			Expect(err).ToNot(HaveOccurred())

			By("sending the correct request", func() {
				Expect(httpClient.GetCallCount()).To(Equal(1))
				url := httpClient.GetArgsForCall(0)
				Expect(url.Path).To(Equal("/api/v1/products/my-super-product"))
				Expect(url.Query().Get("increaseViewCount")).To(Equal("false"))
				Expect(url.Query().Get("isSlug")).To(Equal("true"))
			})

			Expect(product.Slug).To(Equal("my-super-product"))
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.GetReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says there are no products", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product my-super-product not found"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.GetReturns(nil, errors.New("request failed"))
			})

			It("prints the error", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product my-super-product failed: request failed"))
			})
		})

		Context("Unexpected status code", func() {
			BeforeEach(func() {
				httpClient.GetReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Status:     http.StatusText(http.StatusTeapot),
					Body:       ioutil.NopCloser(strings.NewReader("Teapots all the way down")),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting product my-super-product failed: (418)\nTeapots all the way down"))
			})
		})

		Context("Un-parsable response", func() {
			BeforeEach(func() {
				httpClient.GetReturns(test.MakeStringResponse("This totally isn't a valid response"), nil)
			})

			It("prints the error", func() {
				_, err := marketplace.GetProduct("my-super-product")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the response for product my-super-product: invalid character 'T' looking for beginning of value"))
			})
		})
	})

	Describe("GetProductWithVersion", func() {
		var productId string
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				models.SolutionTypeImage)
			productId = product.ProductId
			product.EulaURL = "https://example.com/eula.txt"
			product.OpenSourceDisclosure = &models.OpenSourceDisclosureURLS{
				SourceCodePackageURL: "https://github.com/vmware-labs/marketplace-cli",
			}
			test.AddVerions(product, "0.1.2", "1.2.3")
			response := &pkg.GetProductResponse{
				Response: &pkg.GetProductResponsePayload{
					Data:       product,
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			httpClient.GetReturns(test.MakeJSONResponse(response), nil)

			versionSpecificDetails := &pkg.VersionSpecificDetailsPayloadResponse{
				Response: &pkg.VersionSpecificDetailsPayload{
					Data: &models.VersionSpecificProductDetails{
						EulaURL: "https://example.com/eula-from-a-version.txt",
						OpenSourceDisclosure: &models.OpenSourceDisclosureURLS{
							LicenseDisclosureURL: "https://example.com/osl.txt",
						},
					},
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}

			httpClient.PostJSONReturns(test.MakeJSONResponse(versionSpecificDetails), nil)
		})

		It("returns the product with version specific details", func() {
			product, version, err := marketplace.GetProductWithVersion("my-super-product", "0.1.2")
			Expect(err).ToNot(HaveOccurred())

			Expect(product.Slug).To(Equal("my-super-product"))
			Expect(version.Number).To(Equal("0.1.2"))

			By("sending the correct requests", func() {
				Expect(httpClient.GetCallCount()).To(Equal(1))
				url := httpClient.GetArgsForCall(0)
				Expect(url.Path).To(Equal("/api/v1/products/my-super-product"))
				Expect(url.Query().Get("increaseViewCount")).To(Equal("false"))
				Expect(url.Query().Get("isSlug")).To(Equal("true"))

				Expect(httpClient.PostJSONCallCount()).To(Equal(1))
				url, content := httpClient.PostJSONArgsForCall(0)
				Expect(url.Path).To(Equal(fmt.Sprintf("/api/v1/products/%s/version-details", productId)))
				Expect(url.Query().Get("versionNumber")).To(Equal("0.1.2"))
				payload := content.(*pkg.VersionSpecificDetailsRequestPayload)
				Expect(payload.ProductId).To(Equal(productId))
				Expect(payload.VersionNumber).To(Equal("0.1.2"))
			})

			By("updating the product with the version specific details", func() {
				Expect(product.CurrentVersion).To(Equal("0.1.2"))
				Expect(product.EulaURL).To(Equal("https://example.com/eula-from-a-version.txt"))
				Expect(product.OpenSourceDisclosure.SourceCodePackageURL).To(BeEmpty())
				Expect(product.OpenSourceDisclosure.LicenseDisclosureURL).To(Equal("https://example.com/osl.txt"))
			})
		})

		Context("there was an error getting the product", func() {
			BeforeEach(func() {
				httpClient.GetReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("get product failed"))),
				}, nil)
			})
			It("returns the error", func() {
				_, _, err := marketplace.GetProductWithVersion("my-super-product", "0.1.2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting product my-super-product failed: (418)\nget product failed"))
			})
		})

		Context("the product does not have that version", func() {
			It("returns a version does not exist error", func() {
				product, _, err := marketplace.GetProductWithVersion("my-super-product", "9.9.9")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have version 9.9.9"))

				By("still returning the product", func() {
					Expect(product.Slug).To(Equal("my-super-product"))
				})
			})
		})

		Context("there was an error getting the version specific details", func() {
			BeforeEach(func() {
				httpClient.PostJSONReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("get version specific details failed"))),
				}, nil)
			})
			It("returns an error", func() {
				_, _, err := marketplace.GetProductWithVersion("my-super-product", "0.1.2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting product version details for my-super-product 0.1.2 failed: (418)\nget version specific details failed"))
			})
		})
		Context("there is no version specific details", func() {
			BeforeEach(func() {
				httpClient.PostJSONReturns(&http.Response{
					StatusCode: http.StatusNotFound,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("version specific details not found"))),
				}, nil)
			})
			It("returns an error", func() {
				_, _, err := marketplace.GetProductWithVersion("my-super-product", "0.1.2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product version details for my-super-product 0.1.2 not found"))
			})
		})
		Context("version specific details request returns bad request", func() {
			BeforeEach(func() {
				httpClient.PostJSONReturns(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("bad version specific details request"))),
				}, nil)
			})
			It("returns the product without version specific details", func() {

			})
		})
		Context("version specific details returns bad data", func() {
			BeforeEach(func() {
				httpClient.PostJSONReturns(test.MakeFailingBodyResponse("bad response body"), nil)
			})

			It("returns an error", func() {
				_, _, err := marketplace.GetProductWithVersion("my-super-product", "0.1.2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the response for product my-super-product 0.1.2: bad response body"))
			})
		})
		Context("version specific details returns malformed json", func() {
			BeforeEach(func() {
				httpClient.PostJSONReturns(test.MakeBytesResponse([]byte("}}} this is bad json! {{{")), nil)
			})

			It("returns an error", func() {
				_, _, err := marketplace.GetProductWithVersion("my-super-product", "0.1.2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to parse the response for product my-super-product 0.1.2:"))
			})
		})
	})
})
