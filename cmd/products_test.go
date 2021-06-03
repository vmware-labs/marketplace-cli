package cmd_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/cmd"
	"gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/lib"
	"gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/lib/libfakes"
	"gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/models"
)

var _ = Describe("Products", func() {
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

	Describe("ListProductCmd", func() {
		BeforeEach(func() {
			product1 := CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			AddVerions(product1, "1.1.1")
			product2 := CreateFakeProduct(
				"",
				"My Other Product",
				"my-other-product",
				"PENDING")
			AddVerions(product2, "2.2.2")

			response := &cmd.ListProductResponse{
				Response: &cmd.ListProductResponsePayload{
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
				Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":0,\"pagesize\":20}"))
			})

			By("outputting the response", func() {
				Expect(stdout).To(Say("SLUG              NAME"))
				Expect(stdout).To(Say("my-super-product  My Super Product"))
				Expect(stdout).To(Say("my-other-product  My Other Product"))
				Expect(stdout).To(Say("TOTAL COUNT: 2"))
			})
		})

		Context("Error fetching products", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("request failed"))
			})

			It("prints the error", func() {
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for the list of products failed: request failed"))
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
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
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
				err := cmd.ListProductsCmd.RunE(cmd.ListProductsCmd, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the list of products: invalid character 'T' looking for beginning of value"))
			})
		})
	})

	Describe("GetProductCmd", func() {
		BeforeEach(func() {
			product := CreateFakeProduct(
				"",
				"My Super Product",
				"my-super-product",
				"PENDING")
			AddVerions(product, "0.1.2", "1.2.3")
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

			cmd.GetProductCmd.SetOut(stdout)
			cmd.GetProductCmd.SetErr(stderr)
		})

		It("outputs the product", func() {
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
				Expect(stdout).To(Say("  SLUG              NAME"))
				Expect(stdout).To(Say("  my-super-product  My Super Product"))
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
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("product \"my-super-product\" not found"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("request failed"))
			})

			It("prints the error", func() {
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for product \"my-super-product\" failed: request failed"))
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
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
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
				err := cmd.GetProductCmd.RunE(cmd.GetProductCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the response for product \"my-super-product\": invalid character 'T' looking for beginning of value"))
			})
		})
	})
})
