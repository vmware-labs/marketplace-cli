// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
)

var _ = Describe("ValidateAssetTypeFilter", func() {
	It("allows valid asset types", func() {
		cmd.AssetType = "chart"
		Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
		cmd.AssetType = "image"
		Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
		cmd.AssetType = "metafile"
		Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
		cmd.AssetType = "other"
		Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
		cmd.AssetType = "vm"
		Expect(cmd.ValidateAssetTypeFilter(nil, nil)).To(Succeed())
	})

	When("an invalid asset type is used", func() {
		It("returns an error", func() {
			cmd.AssetType = "dogfood"
			err := cmd.ValidateAssetTypeFilter(nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Unknown asset type: dogfood\nPlease use one of chart, image, metafile, other, vm"))
		})
	})
})

var _ = Describe("ValidateMetaFileType", func() {
	It("allows valid meta file types", func() {
		cmd.MetaFileType = "cli"
		Expect(cmd.ValidateMetaFileType(nil, nil)).To(Succeed())
		cmd.MetaFileType = "config"
		Expect(cmd.ValidateMetaFileType(nil, nil)).To(Succeed())
		cmd.MetaFileType = "other"
		Expect(cmd.ValidateMetaFileType(nil, nil)).To(Succeed())
	})

	When("an invalid meta file type is used", func() {
		It("returns an error", func() {
			cmd.MetaFileType = "dogfood"
			err := cmd.ValidateMetaFileType(nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Unknown meta file type: dogfood\nPlease use one of cli, config, other"))
		})
	})
})
