// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package lib_test

import (
	"io/ioutil"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/lib"
)

var TestConfig = &lib.MarketplaceConfiguration{
	Host: "marketplace.vmware.example",
}

var _ = Describe("Pagination", func() {
	Describe("Apply", func() {
		It("modifies a URL to add pagination with a very specific encoding format", func() {
			pagination := lib.Pagination{
				Page:     1,
				PageSize: 25,
			}

			By("appending to an existing query", func() {
				baseUrl := &url.URL{
					Scheme: "https",
					Host:   "marketplace.vmware.com",
					Path:   "/api/products",
				}
				baseUrl.RawQuery = url.Values{
					"price": []string{"free"},
				}.Encode()
				paginatedUrl := pagination.Apply(baseUrl)
				Expect(paginatedUrl.RawQuery).To(Equal("price=free&pagination={%22page%22:1,%22pageSize%22:25}"))
				Expect(paginatedUrl.String()).To(Equal("https://marketplace.vmware.com/api/products?price=free&pagination={%22page%22:1,%22pageSize%22:25}"))
			})

			By("working with urls without an existing query", func() {
				baseUrl := &url.URL{
					Scheme: "https",
					Host:   "marketplace.vmware.com",
					Path:   "/api/products",
				}
				paginatedUrl := pagination.Apply(baseUrl)
				Expect(paginatedUrl.RawQuery).To(Equal("pagination={%22page%22:1,%22pageSize%22:25}"))
				Expect(paginatedUrl.String()).To(Equal("https://marketplace.vmware.com/api/products?pagination={%22page%22:1,%22pageSize%22:25}"))
			})
		})
	})
})

var _ = Describe("MakeRequest", func() {
	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
	})

	It("Makes a valid request object", func() {
		content := strings.NewReader("everything totally passed")
		request, err := TestConfig.MakeRequest(
			"POST",
			"/api/v1/unit-tests",
			url.Values{
				"color": []string{"blue", "green"},
			},
			map[string]string{
				"Content-Type": "text/plain",
			},
			content,
		)
		Expect(err).ToNot(HaveOccurred())

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

var _ = Describe("MakeGetRequest", func() {
	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
	})

	It("Makes a valid request object", func() {
		request, err := TestConfig.MakeGetRequest(
			"/api/v1/unit-tests",
			url.Values{
				"color": []string{"blue", "green"},
			},
		)
		Expect(err).ToNot(HaveOccurred())

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
