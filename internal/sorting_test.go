// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

var _ = Describe("Sorting", func() {
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
