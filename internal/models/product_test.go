// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("PrepForUpdate", func() {
	It("prepares the product for update", func() {
		product := test.CreateFakeProduct("", "Hyperspace Database", "hyperspace-database", models.SolutionTypeChart)
		test.AddVersions(product, "1.0.0")
		product.Versions = append(product.Versions, &models.Version{
			Number:       "2.0.0",
			Details:      "Details for 2.0.0",
			Status:       "ACTIVE",
			Instructions: "Instructions for 2.0.0",
		})
		product.EncryptionDetails = &models.ProductEncryptionDetails{
			List: []string{"alpha", "bravo", "delta"},
		}

		product.PrepForUpdate()
		By("converting the encryption details list to the encryption hash", func() {
			Expect(product.Encryption).ToNot(BeNil())
			Expect(product.Encryption.List).ToNot(BeNil())
			Expect(product.Encryption.List).To(HaveKeyWithValue("alpha", true))
			Expect(product.Encryption.List).To(HaveKeyWithValue("bravo", true))
			Expect(product.Encryption.List).To(HaveKeyWithValue("delta", true))
		})

		By("Ensuring that both the versions list and the all versions list are in sync", func() {
			Expect(product.AllVersions).To(HaveLen(2))
			Expect(product.AllVersions[0].Number).To(Equal("1.0.0"))
			Expect(product.AllVersions[1].Number).To(Equal("2.0.0"))
			Expect(product.Versions).To(HaveLen(2))
			Expect(product.Versions[0].Number).To(Equal("1.0.0"))
			Expect(product.Versions[1].Number).To(Equal("2.0.0"))
		})
	})
})
