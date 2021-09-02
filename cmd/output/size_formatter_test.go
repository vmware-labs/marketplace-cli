// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vmware-labs/marketplace-cli/v2/cmd/output"
)

var _ = Describe("FormatSize", func() {
	It("prints the object as JSON", func() {
		Expect(FormatSize(123)).To(Equal("123 B"))
		Expect(FormatSize(1234)).To(Equal("1.23 KB"))
		Expect(FormatSize(12345)).To(Equal("12.3 KB"))
		Expect(FormatSize(123456)).To(Equal("123 KB"))
		Expect(FormatSize(1234567)).To(Equal("1.23 MB"))
		Expect(FormatSize(12345678)).To(Equal("12.3 MB"))
		Expect(FormatSize(123456789)).To(Equal("123 MB"))
		Expect(FormatSize(1234567890)).To(Equal("1.23 GB"))
		Expect(FormatSize(12345678901)).To(Equal("12.3 GB"))
		Expect(FormatSize(123456789012)).To(Equal("123 GB"))
		Expect(FormatSize(1234567890123)).To(Equal("1.23 TB"))
		Expect(FormatSize(12345678901234)).To(Equal("12.3 TB"))
		Expect(FormatSize(123456789012345)).To(Equal("123 TB"))
		Expect(FormatSize(1234567890123456)).To(Equal("1.23 PB"))
		Expect(FormatSize(12345678901234567)).To(Equal("12.3 PB"))
		Expect(FormatSize(123456789012345678)).To(Equal("123 PB"))
	})
})
