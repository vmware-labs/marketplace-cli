// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output_test

import (
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

var _ = Describe("HumanOutput", func() {
	var writer *Buffer

	BeforeEach(func() {
		writer = NewBuffer()
	})

	Describe("RenderProduct", func() {
		It("renders the product", func() {
			humanOutput := output.NewHumanOutput(writer, "marketplace.example.com")

			product := &models.Product{
				ProductId:    "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
				Slug:         "hyperspace-database",
				DisplayName:  "Hyperspace Database",
				Status:       "theoretical",
				SolutionType: "HELMCHARTS",
				PublisherDetails: &models.Publisher{
					OrgDisplayName: "Astronomical Widgets",
				},
				Description: &models.Description{
					Summary:     "A database that's out of this world",
					Description: "Connecting to a database should be:<ul><li>Instant</li><li>Robust</li><li>Break the laws of causality</li></ul><br />Our database does just that!",
				},
				AllVersions: []*models.Version{
					{Number: "1.0.0"},
					{Number: "1.2.3"},
					{Number: "0.0.1"},
				},
			}

			err := humanOutput.RenderProduct(product, &models.Version{Number: "1.0.0"})
			Expect(err).ToNot(HaveOccurred())
			Expect(writer).To(Say("Name:      Hyperspace Database"))
			Expect(writer).To(Say("Publisher: Astronomical Widgets"))
			Expect(writer).To(Say("A database that's out of this world"))
			Expect(writer).To(Say(regexp.QuoteMeta("https://marketplace.example.com/services/details/hyperspace-database?slug=true")))
			Expect(writer).To(Say("Product Details:"))

			Expect(writer).To(Say("PRODUCT ID"))
			Expect(writer).To(Say("SLUG"))
			Expect(writer).To(Say("TYPE"))
			Expect(writer).To(Say("LATEST VERSION"))
			Expect(writer).To(Say("STATUS"))

			Expect(writer).To(Say("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"))
			Expect(writer).To(Say("hyperspace-database"))
			Expect(writer).To(Say("HELMCHARTS"))
			Expect(writer).To(Say("1.2.3"))
			Expect(writer).To(Say("theoretical"))

			Expect(writer).To(Say("Assets for 1.0.0:"))
			Expect(writer).To(Say("None"))

			Expect(writer).To(Say("Description:"))
			Expect(writer).To(Say(regexp.QuoteMeta("Connecting to a database should be:\n\n* Instant\n* Robust\n* Break the laws of causality\n\nOur database does just that!")))
		})
	})
})
