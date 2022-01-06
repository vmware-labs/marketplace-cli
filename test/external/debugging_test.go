// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	"regexp"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Debugging", func() {
	steps := NewSteps()

	Scenario("Debugging enabled", func() {
		steps.When("running mkpcli --debug product list")
		steps.Then("the command exits without error")
		steps.And("requests are printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define, "staging")

		define.Then(`^requests are printed$`, func() {
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #0: GET https://gtwstg.market.csp.vmware.com/api/v1/products?managed=true&pagination={%22page%22:1,%22pageSize%22:20}")))
			Eventually(CommandSession.Err).Should(Say("Request #0 Response: 200 OK"))
			Eventually(CommandSession.Out).Should(Say("All products from Tanzu-ISV-ENG"))
		})
	})
})
