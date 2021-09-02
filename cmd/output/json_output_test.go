// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output"
)

var _ = Describe("JSONOutput", func() {
	var (
		jsonOutput *output.JSONOutput
		writer     *Buffer
	)

	BeforeEach(func() {
		writer = NewBuffer()
		jsonOutput = output.NewJSONOutput(writer)
	})

	Describe("PrintJSON", func() {

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

			err := jsonOutput.PrintJSON(data)
			Expect(err).ToNot(HaveOccurred())
			Expect(writer).To(Say("{\"a\":\"some data\",\"awesome\":true,\"numbers\":\\[1,2,3]\\}"))
		})

		Context("Failed to marshal json", func() {
			It("returns an error", func() {
				err := jsonOutput.PrintJSON(func() string { return "this is not json-able!" })
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("json: unsupported type: func() string"))
			})
		})
	})
})
