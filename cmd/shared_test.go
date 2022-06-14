// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
)

var _ = Describe("Products", func() {
	Describe("ValidateAssetTypeFilter", func() {
		It("allows valid asset types", func() {
			cmd.ListAssetsByType = "addon"
			Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
			cmd.ListAssetsByType = "chart"
			Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
			cmd.ListAssetsByType = "image"
			Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
			cmd.ListAssetsByType = "metafile"
			Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
			cmd.ListAssetsByType = "vm"
			Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
		})

		When("an invalid asset type is used", func() {
			It("returns an error", func() {
				cmd.ListAssetsByType = "dogfood"
				err := cmd.ValidateAssetTypeFilter(nil, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unknown asset type: dogfood\nPlease use one of addon, chart, image, metafile, vm"))
			})
		})
	})
})
