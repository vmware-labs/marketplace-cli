// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

var _ = Describe("Pagination", func() {
	Describe("Apply", func() {
		It("modifies a URL to add pagination with a very specific encoding format", func() {
			pagination := internal.Pagination{
				Page:     1,
				PageSize: 25,
			}

			By("appending to an existing query", func() {
				baseURL := &url.URL{
					Scheme: "https",
					Host:   "marketplace.vmware.com",
					Path:   "/api/products",
				}
				baseURL.RawQuery = url.Values{
					"price": []string{"free"},
				}.Encode()
				paginatedURL := pagination.Apply(baseURL)
				Expect(paginatedURL.RawQuery).To(Equal("price=free&pagination={%22page%22:1,%22pageSize%22:25}"))
				Expect(paginatedURL.String()).To(Equal("https://marketplace.vmware.com/api/products?price=free&pagination={%22page%22:1,%22pageSize%22:25}"))
			})

			By("working with urls without an existing query", func() {
				baseURL := &url.URL{
					Scheme: "https",
					Host:   "marketplace.vmware.com",
					Path:   "/api/products",
				}
				paginatedURL := pagination.Apply(baseURL)
				Expect(paginatedURL.RawQuery).To(Equal("pagination={%22page%22:1,%22pageSize%22:25}"))
				Expect(paginatedURL.String()).To(Equal("https://marketplace.vmware.com/api/products?pagination={%22page%22:1,%22pageSize%22:25}"))
			})
		})
	})

	Describe("QueryString", func() {
		It("returns the specific string that works in the query parameter", func() {
			pagination := internal.Pagination{
				Page:     1,
				PageSize: 25,
			}
			queryString := pagination.QueryString()
			Expect(queryString).To(Equal("pagination={%22page%22:1,%22pageSize%22:25}"))
		})
	})
})
