// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Product Version", func() {
	steps := NewSteps()

	Scenario("Listing product versions", func() {
		steps.When("running mkpcli product-version list --product nginx")
		steps.Then("the command exits without error")
		steps.And("the list of product versions is printed")
	})

	Scenario("Getting a single product version", func() {
		steps.When("running mkpcli product-version get --product nginx --product-version 1.21.1_0")
		steps.Then("the command exits without error")
		steps.And("the product version is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define, "production")

		define.Then(`^the list of product versions is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Versions:"))
			Eventually(CommandSession.Out).Should(Say("NUMBER"))
			Eventually(CommandSession.Out).Should(Say("STATUS"))

			Eventually(CommandSession.Out).Should(Say("1.21.3_0"))
			Eventually(CommandSession.Out).Should(Say("ACTIVE"))
			Eventually(CommandSession.Out).Should(Say("1.21.2_0"))
			Eventually(CommandSession.Out).Should(Say("ACTIVE"))
			Eventually(CommandSession.Out).Should(Say("1.21.1_0"))
			Eventually(CommandSession.Out).Should(Say("ACTIVE"))
		})

		define.Then(`^the product version is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Version 1.21.1_0"))
		})
	})
})
