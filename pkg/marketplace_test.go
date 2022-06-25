// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var _ = Describe("Marketplace", func() {
	var marketplace *pkg.Marketplace

	BeforeEach(func() {
		marketplace = &pkg.Marketplace{}
	})

	Describe("DecodeJson", func() {
		type AnObject struct {
			A string `json:"a"`
			B string `json:"b"`
			C string `json:"c"`
		}

		It("parses JSON", func() {
			input := "{\"a\": \"Apple\", \"b\": \"Butter\", \"c\": \"Croissant\"}"
			output := &AnObject{}
			err := marketplace.DecodeJson(strings.NewReader(input), output)

			Expect(err).ToNot(HaveOccurred())
			Expect(output.A).To(Equal("Apple"))
			Expect(output.B).To(Equal("Butter"))
			Expect(output.C).To(Equal("Croissant"))
		})

		Context("Strict decoding enabled", func() {
			BeforeEach(func() {
				marketplace.EnableStrictDecoding()
			})

			It("throws an error on unknown fields", func() {
				input := "{\"extra\": \"How did this get here?\", \"a\": \"Apple\", \"b\": \"Butter\", \"c\": \"Croissant\"}"
				output := &AnObject{}
				err := marketplace.DecodeJson(strings.NewReader(input), output)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("json: unknown field \"extra\""))
			})
		})
	})
})
