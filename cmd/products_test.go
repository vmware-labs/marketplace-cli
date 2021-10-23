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

var _ = Describe("Products", func() {
	var (
		stdout     *Buffer
		stderr     *Buffer
		httpClient *pkgfakes.FakeHTTPClient
		output     *outputfakes.FakeFormat
	)

	BeforeEach(func() {
		stdout = NewBuffer()
		stderr = NewBuffer()

		httpClient = &pkgfakes.FakeHTTPClient{}
		output = &outputfakes.FakeFormat{}
		cmd.Output = output
		cmd.Marketplace = &pkg.Marketplace{
			Client: httpClient,
			Output: stderr,
		}
	})
	Describe("ListProductsCmd", func() {
		BeforeEach(func() {
			products := []*models.Product{
				test.CreateFakeProduct(
					"",
					"My Super Product",
					"my-super-product",
					"PENDING"),
				test.CreateFakeProduct(
					"",
					"My Other Product",
					"my-other-product",
					"PENDING"),
			}
			response := &pkg.ListProductResponse{
				Response: &pkg.ListProductResponsePayload{
					Products:   products,
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

			cmd.ListProductsCmd.SetOut(stdout)
			cmd.ListProductsCmd.SetErr(stderr)
		})

		It("sends the right request", func() {
			err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
			Expect(err).ToNot(HaveOccurred())

			By("sending the correct request", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/products"))
				Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":1,\"pageSize\":20}"))
			})

			By("outputting the response", func() {
				Expect(output.RenderProductsCallCount()).To(Equal(1))
				products := output.RenderProductsArgsForCall(0)
				Expect(products).To(HaveLen(2))
				Expect(products[0].Slug).To(Equal("my-super-product"))
				Expect(products[1].Slug).To(Equal("my-other-product"))
			})
		})

		Context("Error fetching products", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for the list of products failed: marketplace request failed: request failed"))
			})
		})

		Context("Unexpected status code", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Status:     http.StatusText(http.StatusTeapot),
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("Teapots!"))),
				}, nil)
			})

			It("prints the error", func() {
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting the list of products failed: (418) I'm a teapot: Teapots!"))
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
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the list of products: invalid character 'T' looking for beginning of value"))
			})
		})
	})

	Describe("GetProductCmd", func() {
		BeforeEach(func() {
			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")

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

			cmd.GetProductCmd.SetOut(stdout)
			cmd.GetProductCmd.SetErr(stderr)
		})

		It("sends the right request", func() {
			cmd.ProductSlug = "my-super-product"
			err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
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
				Expect(output.RenderProductCallCount()).To(Equal(1))
				product := output.RenderProductArgsForCall(0)
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(product.DisplayName).To(Equal("My Super Product"))
			})
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says that the product was not found", func() {
				cmd.ProductSlug = "my-super-product"
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: marketplace request failed: request failed"))
			})
		})
	})

	Describe("AddProductVersionCmd", func() {
		var productID string
		BeforeEach(func() {
			productID = uuid.New().String()
			product := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "0.1.2", "1.2.3")
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
			test.AddVerions(updatedProduct, "0.1.2", "1.2.3", "9.9.9")

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

			cmd.AddProductVersionCmd.SetOut(stdout)
			cmd.AddProductVersionCmd.SetErr(stderr)
		})

		It("sends the right requests", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "9.9.9"
			err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
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
			})

			By("outputting the response", func() {
				Expect(output.RenderVersionsCallCount()).To(Equal(1))
				product := output.RenderVersionsArgsForCall(0)
				Expect(product.Slug).To(Equal("my-super-product"))
				Expect(product.AllVersions).To(HaveLen(3))
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
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("Version already exists", func() {
			It("says that the version already exists", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" already has version 1.2.3"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(0, nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: marketplace request failed: request failed"))
			})
		})

		Context("Unexpected status code", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(0,
					&http.Response{
						StatusCode: http.StatusTeapot,
						Status:     http.StatusText(http.StatusTeapot),
					}, nil)
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting product \"my-super-product\" failed: (418)"))
			})
		})

		Context("Un-parsable response", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(0,
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader("This totally isn't a valid response")),
					}, nil)
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the response for product \"my-super-product\": invalid character 'T' looking for beginning of value"))
			})
		})

		Context("No permission to update product", func() {
			BeforeEach(func() {
				httpClient.DoReturnsOnCall(1,
					&http.Response{
						StatusCode: http.StatusForbidden,
						Body:       ioutil.NopCloser(strings.NewReader("{\"response\":{\"message\":\"User is not authorized to perform this action\"}}\n")),
					}, nil)
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("you do not have permission to modify the product \"my-super-product\""))
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
				cmd.ProductVersion = "9.9.9"
				err := cmd.AddProductVersionCmd.RunE(cmd.AddProductVersionCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("updating product \"my-super-product\" failed: (418)\nTeapots all the way down"))
			})
		})
	})

	Describe("ListProductVersionsCmd", func() {
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

			cmd.ListProductVersionsCmd.SetOut(stdout)
			cmd.ListProductVersionsCmd.SetErr(stderr)
		})

		It("sends the right request", func() {
			cmd.ProductSlug = "my-super-product"
			err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
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
				Expect(output.RenderVersionsCallCount()).To(Equal(1))
				product := output.RenderVersionsArgsForCall(0)
				Expect(product.Slug).To(Equal("my-super-product"))
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
				err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
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
				cmd.ProductSlug = "my-super-product"
				err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
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
				cmd.ProductSlug = "my-super-product"
				err := cmd.ListProductVersionsCmd.RunE(cmd.ListProductVersionsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the response for product \"my-super-product\": invalid character 'T' looking for beginning of value"))
			})
		})
	})
})
