// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

var _ = Describe("Pagination", func() {
	Describe("QueryString", func() {
		It("returns the specific string that works in the query parameter", func() {
			pagination := &internal.Pagination{
				Page:     1,
				PageSize: 25,
			}
			queryString := pagination.QueryString()
			Expect(queryString).To(Equal("pagination={%22page%22:1,%22pageSize%22:25}"))
		})
	})
})
