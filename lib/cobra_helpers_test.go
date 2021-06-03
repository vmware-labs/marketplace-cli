// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package lib_test

import (
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vmware-labs/marketplace-cli/v2/lib"
)

var _ = Describe("FileArg", func() {
	Context("no args", func() {
		It("returns an error", func() {
			FileArgFn := FileArg()
			err := FileArgFn(nil, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("requires a single file as an argument"))
		})
	})

	Context("file does not exist", func() {
		It("returns an error", func() {
			FileArgFn := FileArg()
			err := FileArgFn(nil, []string{"a", "a"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("requires a single file as an argument"))
		})
	})

	Context("too many args", func() {
		It("returns an error", func() {
			FileArgFn := FileArg()
			err := FileArgFn(nil, []string{"/this/file/does/not/exist"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("stat /this/file/does/not/exist: no such file or directory"))
		})
	})

	Context("file exists", func() {
		It("returns without error", func() {
			_, file, _, _ := runtime.Caller(0)

			FileArgFn := FileArg()
			err := FileArgFn(nil, []string{file})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
