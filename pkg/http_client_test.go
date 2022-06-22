// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var _ = Describe("HTTP Client", func() {
	Describe("ApplyParameters", func() {
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
})
