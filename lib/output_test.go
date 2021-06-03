// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package lib_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/lib"
)

var _ = Describe("NewTableWriter", func() {
	It("prints a table", func() {
		output := NewBuffer()
		tableWriter := NewTableWriter(output, "name", "kind")
		tableWriter.Append([]string{"Pete", "Developer"})
		tableWriter.Append([]string{"Apples", "fruit"})
		tableWriter.Render()

		Expect(output).To(Say("  NAME    KIND"))
		Expect(output).To(Say("  Pete    Developer"))
		Expect(output).To(Say("  Apples  fruit"))
	})
})

type BadWriter struct{}

func (w *BadWriter) Write(p []byte) (n int, err error) { return 0, errors.New("bad write") }

var _ = Describe("PrintJson", func() {
	It("prints the object as JSON", func() {
		data := struct {
			A             string `json:"a"`
			Awesome       bool   `json:"awesome"`
			MisMatchedKey []int  `json:"numbers"`
		}{
			A:             "some data",
			Awesome:       true,
			MisMatchedKey: []int{1, 2, 3},
		}

		out := NewBuffer()
		err := PrintJson(out, data)
		Expect(err).ToNot(HaveOccurred())
		Expect(out).To(Say("{\"a\":\"some data\",\"awesome\":true,\"numbers\":\\[1,2,3]\\}"))
	})

	Context("Failed to marshal json", func() {
		It("returns an error", func() {
			out := NewBuffer()
			err := PrintJson(out, func() string { return "this is not json-able!" })
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("json: unsupported type: func() string"))
		})
	})

	Context("Failed to marshal json", func() {
		It("returns an error", func() {
			out := &BadWriter{}
			err := PrintJson(out, false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("bad write"))
		})
	})
})
