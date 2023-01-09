// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output"
)

var _ = Describe("EncodedOutput", func() {
	var writer *Buffer

	BeforeEach(func() {
		writer = NewBuffer()
	})

	Describe("Print", func() {

		Context("JSON", func() {
			It("prints the object as JSON", func() {
				encodedOutput := output.NewJSONOutput(writer)
				data := struct {
					A             string `json:"a"`
					Awesome       bool   `json:"awesome"`
					MisMatchedKey []int  `json:"numbers"`
				}{
					A:             "some data",
					Awesome:       true,
					MisMatchedKey: []int{1, 2, 3},
				}

				err := encodedOutput.Print(data)
				Expect(err).ToNot(HaveOccurred())
				Expect(writer).To(Say("{\"a\":\"some data\",\"awesome\":true,\"numbers\":\\[1,2,3]\\}"))
			})
		})

		Context("YAML", func() {
			It("prints the object as YAML", func() {
				encodedOutput := output.NewYAMLOutput(writer)
				data := struct {
					A             string `json:"a"`
					Awesome       bool   `json:"awesome"`
					MisMatchedKey []int  `json:"numbers"`
				}{
					A:             "some data",
					Awesome:       true,
					MisMatchedKey: []int{1, 2, 3},
				}

				err := encodedOutput.Print(data)
				Expect(err).ToNot(HaveOccurred())
				Expect(writer).To(Say("a: some data\nawesome: true\nmismatchedkey:\n    - 1\n    - 2\n    - 3\n"))
			})
		})

		Context("Encoding fails", func() {
			It("returns an error", func() {
				encodedOutput := output.NewJSONOutput(writer)
				encodedOutput.Marshall = func(v interface{}) ([]byte, error) {
					return nil, errors.New("encoding went bad")
				}

				err := encodedOutput.Print([]string{"a", "b"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("encoding went bad"))
			})
		})
	})
})
