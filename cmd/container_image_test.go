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

var _ = Describe("ContainerImage", func() {
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
			cmd.ContainerImageProductSlug = "my-super-product"
			cmd.ContainerImageProductVersion = "1.2.3"
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
				Expect(output.RenderContainerImagesCallCount()).To(Equal(1))
				images := output.RenderContainerImagesArgsForCall(0)
				Expect(images.AppVersion).To(Equal("1.2.3"))
				Expect(images.DockerURLs).To(HaveLen(1))
				Expect(images.DockerURLs[0].ImageTags).To(HaveLen(2))
				Expect(images.DockerURLs[0].ImageTags[0].Tag).To(Equal("0.0.1"))
				Expect(images.DockerURLs[0].ImageTags[0].Type).To(Equal("FIXED"))
				Expect(images.DockerURLs[0].ImageTags[1].Tag).To(Equal("latest"))
				Expect(images.DockerURLs[0].ImageTags[1].Type).To(Equal("FLOATING"))
			})
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says there are no products", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
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
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: marketplace request failed: request failed"))
			})
		})

		Context("No product version found", func() {
			It("says that the version does not exist", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "9.9.9"
				err := cmd.ListContainerImageCmd.RunE(cmd.ListContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
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
			cmd.ContainerImageProductSlug = "my-super-product"
			cmd.ContainerImageProductVersion = "1.2.3"
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
				Expect(output.RenderContainerImageCallCount()).To(Equal(1))
				containerImage := output.RenderContainerImageArgsForCall(0)
				Expect(containerImage.Url).To(Equal("myId"))
				Expect(containerImage.ImageTags).To(ContainElement(&models.DockerImageTag{
					Tag:  "0.0.1",
					Type: "FIXED",
				}))
				Expect(containerImage.ImageTags).To(ContainElement(&models.DockerImageTag{
					Tag:  "latest",
					Type: "FLOATING",
				}))
			})
		})

		Context("No product found", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)
			})

			It("says that the product was not found", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "myId"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No version found", func() {
			It("says that the version does not exist", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "9.9.9"
				cmd.ImageRepository = "myId"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 9.9.9"))
			})
		})

		Context("No container images for version", func() {
			It("says that the version does not exist", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "thisImageDoesNotExist"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("my-super-product 1.2.3 does not have the container image \"thisImageDoesNotExist\""))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "myId"
				err := cmd.GetContainerImageCmd.RunE(cmd.GetContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: marketplace request failed: request failed"))
			})
		})
	})

	Describe("AttachContainerImageCmd", func() {
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

			cmd.AttachContainerImageCmd.SetOut(stdout)
			cmd.AttachContainerImageCmd.SetErr(stderr)
		})

		It("sends the right requests", func() {
			cmd.ContainerImageProductSlug = "my-super-product"
			cmd.ContainerImageProductVersion = "1.2.3"
			cmd.ImageRepository = "python"
			cmd.ImageTag = "1.2.3"
			cmd.ImageTagType = cmd.ImageTagTypeFixed
			err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
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
				Expect(output.RenderContainerImagesCallCount()).To(Equal(1))
				images := output.RenderContainerImagesArgsForCall(0)
				Expect(images.DockerURLs).To(HaveLen(2))
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

				cmd.AttachContainerImageCmd.SetOut(stdout)
				cmd.AttachContainerImageCmd.SetErr(stderr)
			})

			It("sends the right requests", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
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
					Expect(output.RenderContainerImagesCallCount()).To(Equal(1))
					images := output.RenderContainerImagesArgsForCall(0)
					Expect(images.DockerURLs[0].ImageTags).To(HaveLen(2))
					Expect(images.DockerURLs[0].ImageTags[0].Tag).To(Equal("latest"))
					Expect(images.DockerURLs[0].ImageTags[0].Type).To(Equal("FLOATING"))
					Expect(images.DockerURLs[0].ImageTags[1].Tag).To(Equal("5.5.5"))
					Expect(images.DockerURLs[0].ImageTags[1].Type).To(Equal("FIXED"))
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
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("No product version found", func() {
			It("says there are no versions", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "0.0.0"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" does not have a version 0.0.0"))
			})
		})

		Context("Container image already exists", func() {
			It("says the image already exists", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "latest"
				cmd.ImageTagType = cmd.ImageTagTypeFloating
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("my-super-product 1.2.3 already has the tag nginx:latest"))
			})
		})

		Context("invalid tag type", func() {
			It("prints the error", func() {
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = "fancy"
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid image tag type: FANCY. must be either \"FIXED\" or \"FLOATING\""))
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
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
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
				cmd.ContainerImageProductSlug = "my-super-product"
				cmd.ContainerImageProductVersion = "1.2.3"
				cmd.ImageRepository = "nginx"
				cmd.ImageTag = "5.5.5"
				cmd.ImageTagType = cmd.ImageTagTypeFixed
				err := cmd.AttachContainerImageCmd.RunE(cmd.AttachContainerImageCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("updating product \"my-super-product\" failed: (418)\nTeapots all the way down"))
			})
		})
	})
})
