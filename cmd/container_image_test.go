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
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("ContainerImage", func() {
	var (
		stdout     *Buffer
		stderr     *Buffer
		httpClient *pkgfakes.FakeHTTPClient
	)

	BeforeEach(func() {
		httpClient = &pkgfakes.FakeHTTPClient{}
		cmd.Marketplace = &pkg.Marketplace{
			Client: httpClient,
		}
		stdout = NewBuffer()
		stderr = NewBuffer()
	})

	Describe("ListContainerImageCmd", func() {
		BeforeEach(func() {
			container := test.CreateFakeContainerImage(
				"myId",
				"0.0.1",
				"latest",
			)

			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")
			test.AddContainerImages(product, "1.2.3", "Machine wash cold with like colors", container)
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

			cmd.ListContainerImageCmd.SetOut(stdout)
			cmd.ListContainerImageCmd.SetErr(stderr)
		})

		It("outputs the container images", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "1.2.3"
			err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
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
				Expect(stdout).To(Say("IMAGE  TAGS"))
				Expect(stdout).To(Say("myId   0.0.1, latest"))
				Expect(stdout).To(Say("Deployment instructions:"))
				Expect(stdout).To(Say("Machine wash cold with like colors"))
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
				err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
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
				err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: request failed"))
			})
		})

		Context("No product version found", func() {
			It("says that the version does not exist", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
			})
		})

		Context("No container images", func() {
			It("says there are no container images", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "2.3.4"
				err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
				Expect(err).ToNot(HaveOccurred())
				Expect(stdout).To(Say("product \"my-super-product\" 2.3.4 does not have any container images"))
			})
		})
	})

	Describe("GetContainerImageCmd", func() {
		BeforeEach(func() {
			container := test.CreateFakeContainerImage(
				"myId",
				"0.0.1",
				"latest",
			)

			product := test.CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3", "2.3.4")
			test.AddContainerImages(product, "1.2.3", "Machine wash cold with like colors", container)
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

			cmd.GetContainerImageCmd.SetOut(stdout)
			cmd.GetContainerImageCmd.SetErr(stderr)
		})

		It("sends the right request", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "1.2.3"
			cmd.ImageRepository = "myId"
			err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
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
				Expect(stdout).To(Say("TAG     TYPE"))
				Expect(stdout).To(Say("0.0.1   FIXED"))
				Expect(stdout).To(Say("latest  FLOATING"))
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
				cmd.ProductVersion = "1.2.3"
				cmd.ImageRepository = "myId"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No version found", func() {
			It("says that the version does not exist", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "9.9.9"
				cmd.ImageRepository = "myId"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
			})
		})

		Context("No container images for version", func() {
			It("says that the version does not exist", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ImageRepository = "thisImageDoesNotExist"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" 1.2.3 does not have the container image \"thisImageDoesNotExist\""))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ImageRepository = "myId"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: request failed"))
			})
		})
	})

	Describe("CreateContainerImageCmd", func() {
		var productID string
		BeforeEach(func() {
			nginx := test.CreateFakeContainerImage("nginx", "latest")
			python := test.CreateFakeContainerImage("python", "1.2.3")

			productID = uuid.New().String()
			product := test.CreateFakeProduct(
				productID,
				"My Super Product",
				"my-super-product",
				"PENDING")
			test.AddVerions(product, "1.2.3")
			test.AddContainerImages(product, "1.2.3", "Machine wash cold with like colors", nginx)
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
			test.AddContainerImages(updatedProduct, "1.2.3", "Machine wash cold with like colors", nginx, python)
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

			cmd.CreateContainerImageCmd.SetOut(stdout)
			cmd.CreateContainerImageCmd.SetErr(stderr)
		})

		It("sends the right requests", func() {
			cmd.ProductSlug = "my-super-product"
			cmd.ProductVersion = "1.2.3"
			cmd.ImageRepository = "python"
			cmd.ImageTag = "1.2.3"
			cmd.ImageTagType = cmd.ImageTagTypeFixed
			err := cmd.CreateContainerImageCmd.RunE(cmd.CreateContainerImageCmd, []string{""})
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
				Expect(stdout).To(Say("IMAGE   TAGS"))
				Expect(stdout).To(Say("nginx   latest"))
				Expect(stdout).To(Say("python  1.2.3"))
				Expect(stdout).To(Say("Deployment instructions:"))
				Expect(stdout).To(Say("Machine wash cold with like colors"))
			})
		})

		Context("Adding a new tag to an existing container image", func() {
			BeforeEach(func() {
				nginx := test.CreateFakeContainerImage("nginx", "latest")
				nginxUpdated := test.CreateFakeContainerImage("nginx", "latest", "5.5.5")

				productID = uuid.New().String()
				product := test.CreateFakeProduct(
					productID,
					"My Super Product",
					"my-super-product",
					"PENDING")
				test.AddVerions(product, "1.2.3")
				test.AddContainerImages(product, "1.2.3", "Machine wash cold with like colors", nginx)
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
				test.AddContainerImages(updatedProduct, "1.2.3", "Machine wash cold with like colors", nginxUpdated)
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

				cmd.CreateContainerImageCmd.SetOut(stdout)
				cmd.CreateContainerImageCmd.SetErr(stderr)
			})

			It("sends the right requests", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.CreateContainerImageCmd.RunE(cmd.CreateContainerImageCmd, []string{""})
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
					Expect(stdout).To(Say("IMAGE  TAGS"))
					Expect(stdout).To(Say("nginx  latest"))
					Expect(stdout).To(Say("Deployment instructions:"))
					Expect(stdout).To(Say("Machine wash cold with like colors"))
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
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.CreateContainerImageCmd.RunE(cmd.CreateContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No product version found", func() {
			It("says there are no versions", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "0.0.0"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.CreateContainerImageCmd.RunE(cmd.CreateContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 0.0.0, please add it first"))
			})
		})

		Context("Container image already exists", func() {
			It("says the image already exists", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "latest"
				cmd.ImageTagType = cmd.ImageTagTypeFloating
				err := cmd.CreateContainerImageCmd.RunE(cmd.CreateContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" 1.2.3 already has the container image nginx:latest"))
			})
		})

		Context("invalid tag type", func() {
			It("prints the error", func() {
				cmd.ProductSlug = "my-super-product"
				cmd.ProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = "fancy"
				err := cmd.CreateContainerImageCmd.RunE(cmd.CreateContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid image tag type: FANCY. must be either \"FIXED\" or \"FLOATING\""))
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
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.CreateContainerImageCmd.RunE(cmd.CreateContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("updating product \"my-super-product\" failed: (418)\nTeapots all the way down"))
			})
		})
	})
})
