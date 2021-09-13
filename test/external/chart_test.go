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

var _ = Describe("Chart", func() {
	steps := NewSteps()

	Scenario("Listing charts", func() {
		steps.When("running mkpcli chart list --product nginx --product-version 1.21.1_0")
		steps.Then("the command exits without error")
		steps.And("the table of charts is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define, "production")

		define.Then(`^the table of charts is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("ID"))
			Eventually(CommandSession.Out).Should(Say("VERSION"))
			Eventually(CommandSession.Out).Should(Say("URL"))
			Eventually(CommandSession.Out).Should(Say("REPOSITORY"))

			Eventually(CommandSession.Out).Should(Say("9b5a4eb0-d42e-4c14-bbba-2a94ac1fe1f9"))
			Eventually(CommandSession.Out).Should(Say("9.3.6"))
			Eventually(CommandSession.Out).Should(Say("https://charts.bitnami.com/bitnami/nginx-9.3.6.tgz"))
			Eventually(CommandSession.Out).Should(Say("Bitnami charts repo @ Github https://github.com/bitnami/charts/tree/master/bitnami/nginx"))
		})
	})
})
