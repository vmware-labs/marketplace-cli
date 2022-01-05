// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
)

var _ = Describe("Marketplace", func() {
	var (
		httpClient  *pkgfakes.FakeHTTPClient
		marketplace *pkg.Marketplace
	)

	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			Client: httpClient,
			Host:   "marketplace.vmware.example",
		}
		httpClient.DoReturns(&http.Response{
			StatusCode: http.StatusTeapot,
		}, nil)
	})

	Describe("MakeURL", func() {
		It("sets the scheme and host", func() {
			marketplaceURL := marketplace.MakeURL(
				"/path/to/products/",
				url.Values{
					"color": []string{"red", "blue"},
				},
			)
			Expect(marketplaceURL.Scheme).To(Equal("https"))
			Expect(marketplaceURL.Host).To(Equal("marketplace.vmware.example"))
			Expect(marketplaceURL.Path).To(Equal("/path/to/products/"))
			Expect(marketplaceURL.RawQuery).To(Equal("color=red&color=blue"))
		})

		Context("nil values", func() {
			It("still works", func() {
				marketplaceURL := marketplace.MakeURL("/there/are/no/options", nil)
				Expect(marketplaceURL.Scheme).To(Equal("https"))
				Expect(marketplaceURL.Host).To(Equal("marketplace.vmware.example"))
				Expect(marketplaceURL.Path).To(Equal("/there/are/no/options"))
				Expect(marketplaceURL.RawQuery).To(Equal(""))
			})
		})
	})

	Describe("SendRequest", func() {
		It("sends a valid request", func() {
			content := strings.NewReader("everything totally passed")
			response, err := marketplace.SendRequest(
				"POST",
				marketplace.MakeURL(
					"/api/v1/unit-tests",
					url.Values{
						"color": []string{"blue", "green"},
					},
				),
				map[string]string{
					"Content-Type": "text/plain",
				},
				content,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusTeapot))

			Expect(httpClient.DoCallCount()).To(Equal(1))
			request := httpClient.DoArgsForCall(0)

			Expect(request.Method).To(Equal("POST"))

			By("building the right url", func() {
				Expect(request.URL.Scheme).To(Equal("https"))
				Expect(request.URL.Host).To(Equal("marketplace.vmware.example"))
				Expect(request.URL.Path).To(Equal("/api/v1/unit-tests"))
				Expect(request.URL.Query().Encode()).To(Equal("color=blue&color=green"))
			})

			By("setting the right headers", func() {
				Expect(request.Header.Get("Accept")).To(Equal("application/json"))
				Expect(request.Header.Get("csp-auth-token")).To(Equal("secrets"))
				Expect(request.Header.Get("Content-Type")).To(Equal("text/plain"))
			})

			By("including the right content", func() {
				Expect(ioutil.ReadAll(request.Body)).To(Equal([]byte("everything totally passed")))
			})
		})
	})

	var _ = Describe("Get", func() {
		It("sends a valid request", func() {
			response, err := marketplace.Get(
				marketplace.MakeURL(
					"/api/v1/unit-tests",
					url.Values{
						"color": []string{"blue", "green"},
					},
				),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusTeapot))

			Expect(httpClient.DoCallCount()).To(Equal(1))
			request := httpClient.DoArgsForCall(0)

			Expect(request.Method).To(Equal("GET"))

			By("building the right url", func() {
				Expect(request.URL.Scheme).To(Equal("https"))
				Expect(request.URL.Host).To(Equal("marketplace.vmware.example"))
				Expect(request.URL.Path).To(Equal("/api/v1/unit-tests"))
				Expect(request.URL.Query().Encode()).To(Equal("color=blue&color=green"))
			})

			By("setting the right headers", func() {
				Expect(request.Header.Get("Accept")).To(Equal("application/json"))
				Expect(request.Header.Get("csp-auth-token")).To(Equal("secrets"))
			})
		})
	})

	Describe("Put", func() {
		It("sends a valid request", func() {
			content := strings.NewReader("everything totally passed")
			response, err := marketplace.Put(
				marketplace.MakeURL(
					"/api/v1/unit-tests",
					url.Values{
						"color": []string{"blue", "green"},
					},
				),
				content,
				"text/plain",
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusTeapot))

			Expect(httpClient.DoCallCount()).To(Equal(1))
			request := httpClient.DoArgsForCall(0)

			Expect(request.Method).To(Equal("PUT"))

			By("building the right url", func() {
				Expect(request.URL.Scheme).To(Equal("https"))
				Expect(request.URL.Host).To(Equal("marketplace.vmware.example"))
				Expect(request.URL.Path).To(Equal("/api/v1/unit-tests"))
				Expect(request.URL.Query().Encode()).To(Equal("color=blue&color=green"))
			})

			By("setting the right headers", func() {
				Expect(request.Header.Get("Accept")).To(Equal("application/json"))
				Expect(request.Header.Get("csp-auth-token")).To(Equal("secrets"))
				Expect(request.Header.Get("Content-Type")).To(Equal("text/plain"))
			})

			By("including the right content", func() {
				Expect(ioutil.ReadAll(request.Body)).To(Equal([]byte("everything totally passed")))
			})
		})
	})
})
