// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

var _ = Describe("Sorting", func() {
	Describe("Apply", func() {
		It("modifies a URL to add sorting with a very specific encoding format", func() {
			sorting := &internal.Sorting{
				Key:       internal.SortKeyCreationDate,
				Direction: internal.SortDirectionDescending,
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
				paginatedURL := sorting.Apply(baseURL)
				Expect(paginatedURL.RawQuery).To(Equal("price=free&sortBy={%22key%22:%22createdOn%22,%22direction%22:%22DESC%22}"))
				Expect(paginatedURL.String()).To(Equal("https://marketplace.vmware.com/api/products?price=free&sortBy={%22key%22:%22createdOn%22,%22direction%22:%22DESC%22}"))
			})

			By("working with urls without an existing query", func() {
				baseURL := &url.URL{
					Scheme: "https",
					Host:   "marketplace.vmware.com",
					Path:   "/api/products",
				}
				paginatedURL := sorting.Apply(baseURL)
				Expect(paginatedURL.RawQuery).To(Equal("sortBy={%22key%22:%22createdOn%22,%22direction%22:%22DESC%22}"))
				Expect(paginatedURL.String()).To(Equal("https://marketplace.vmware.com/api/products?sortBy={%22key%22:%22createdOn%22,%22direction%22:%22DESC%22}"))
			})
		})
	})

	Describe("QueryString", func() {
		It("returns the specific string that works in the query parameter", func() {
			sorting := &internal.Sorting{
				Key:       internal.SortKeyUpdateDate,
				Direction: internal.SortDirectionAscending,
			}
			queryString := sorting.QueryString()
			Expect(queryString).To(Equal("sortBy={%22key%22:%22updatedOn%22,%22direction%22:%22ASC%22}"))
		})
	})
})
