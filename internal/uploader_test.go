// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

var _ = Describe("MakeUniqueFilename", func() {
	It("Returns a new unique filename", func() {
		Expect(internal.MakeUniqueFilename("test.binary")).To(MatchRegexp("test-[0-9]*.binary"))
		Expect(internal.MakeUniqueFilename("two.dots.test")).To(MatchRegexp("two.dots-[0-9]*.test"))
		Expect(internal.MakeUniqueFilename("no-dots")).To(MatchRegexp("no-dots-[0-9]*"))
	})
})
