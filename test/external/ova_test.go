// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	"fmt"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("OVA", func() {
	steps := NewSteps()

	var (
		slug                = "test-ova-product-11"
		knownProductVersion = "0.0.0"
	)

	Scenario("Listing OVAs", func() {
		steps.When(fmt.Sprintf("running mkpcli ova list --product %s --product-version %s", slug, knownProductVersion))
		steps.Then("the command exits without error")
		steps.And("the table of OVAs is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the table of OVAs is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("NAME"))
			Eventually(CommandSession.Out).Should(Say("SIZE"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("FILES"))
			Eventually(CommandSession.Out).Should(Say("photon-hw13-uefi-4-1622584869708"))
			Eventually(CommandSession.Out).Should(Say("234018527"))
			Eventually(CommandSession.Out).Should(Say("vcsp.ovf"))
			Eventually(CommandSession.Out).Should(Say("4"))
		})
	})
})
