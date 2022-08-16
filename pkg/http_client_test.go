// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
)

var _ = Describe("HTTP Client", func() {
	var (
		httpClient     *pkg.DebuggingClient
		performRequest *pkgfakes.FakePerformRequestFunc
	)

	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
		performRequest = &pkgfakes.FakePerformRequestFunc{}
		performRequest.Returns(&http.Response{
			StatusCode: http.StatusTeapot,
		}, nil)
		httpClient = pkg.NewClient(nil, false, false, false)
		httpClient.PerformRequest = performRequest.Spy
	})

	var _ = Describe("Get", func() {
		It("sends a valid request", func() {
			response, err := httpClient.Get(pkg.MakeURL(
				"marketplace.vmware.example",
				"/api/v1/unit-tests",
				url.Values{
					"color": []string{"blue", "green"},
				},
			))
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusTeapot))

			Expect(performRequest.CallCount()).To(Equal(1))
			request := performRequest.ArgsForCall(0)

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
			response, err := httpClient.Put(
				pkg.MakeURL(
					"marketplace.vmware.example",
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

			Expect(performRequest.CallCount()).To(Equal(1))
			request := performRequest.ArgsForCall(0)

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
				Expect(io.ReadAll(request.Body)).To(Equal([]byte("everything totally passed")))
			})
		})
	})
})

var _ = Describe("MakeURL", func() {
	It("sets the scheme and host", func() {
		url := pkg.MakeURL(
			"marketplace.vmware.example",
			"/path/to/products/",
			url.Values{
				"color": []string{"red", "blue"},
			},
		)
		Expect(url.Scheme).To(Equal("https"))
		Expect(url.Host).To(Equal("marketplace.vmware.example"))
		Expect(url.Path).To(Equal("/path/to/products/"))
		Expect(url.RawQuery).To(Equal("color=red&color=blue"))
	})

	Context("nil values", func() {
		It("still works", func() {
			url := pkg.MakeURL("marketplace.vmware.example", "/there/are/no/options", nil)
			Expect(url.Scheme).To(Equal("https"))
			Expect(url.Host).To(Equal("marketplace.vmware.example"))
			Expect(url.Path).To(Equal("/there/are/no/options"))
			Expect(url.RawQuery).To(Equal(""))
		})
	})
})

var _ = Describe("ApplyParameters", func() {
	var (
		pagination *internal.Pagination
		sorting    *internal.Sorting
	)

	BeforeEach(func() {
		pagination = &internal.Pagination{
			Page:     1,
			PageSize: 25,
		}
		sorting = &internal.Sorting{
			Key:       internal.SortKeyCreationDate,
			Direction: internal.SortDirectionDescending,
		}
	})

	It("combines the parameter objects with existing query string", func() {
		url, _ := url.Parse("https://example.com/path?testing=great&flag")
		Expect(url.RawQuery).To(Equal("testing=great&flag"))
		pkg.ApplyParameters(url, pagination, sorting)
		Expect(url.RawQuery).To(Equal("testing=great&flag&pagination={%22page%22:1,%22pageSize%22:25}&sortBy={%22key%22:%22createdOn%22,%22direction%22:%22DESC%22}"))
	})

	When("there is no existing query string", func() {
		It("sets the query string", func() {
			url, _ := url.Parse("https://example.com/path")
			Expect(url.RawQuery).To(Equal(""))
			pkg.ApplyParameters(url, sorting, pagination)
			Expect(url.RawQuery).To(Equal("sortBy={%22key%22:%22createdOn%22,%22direction%22:%22DESC%22}&pagination={%22page%22:1,%22pageSize%22:25}"))
		})
	})

	When("there are no parameters", func() {
		It("doesn't append to anything", func() {
			url, _ := url.Parse("https://example.com/path?testing=great&flag")
			Expect(url.RawQuery).To(Equal("testing=great&flag"))
			pkg.ApplyParameters(url)
			Expect(url.RawQuery).To(Equal("testing=great&flag"))
		})
	})
})
