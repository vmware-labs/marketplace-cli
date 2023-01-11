// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var _ = Describe("Filter", func() {
	Describe("ListProductFilter", func() {
		Context("no filters", func() {
			It("returns an empty filter", func() {
				filter := pkg.ListProductFilter{}
				Expect(filter.QueryString()).To(Equal("filters={}"))
			})
		})

		Context("organization id provided", func() {
			It("adds a filter with the organization id", func() {
				filter := pkg.ListProductFilter{
					OrgIds: []string{"my-org-id"},
				}
				Expect(filter.QueryString()).To(Equal("filters={%22Publishers%22:[%22my-org-id%22]}"))
			})
		})

		Context("text filter", func() {
			It("adds a filter with the text", func() {
				filter := pkg.ListProductFilter{
					Text: "tanzu",
				}
				Expect(filter.QueryString()).To(Equal("filters={%22search%22:%22tanzu%22}"))
			})
		})
	})
})
