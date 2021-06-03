// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

// +build enemy

package features_test

import (
	"fmt"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/features"
)

var _ = Describe("Chart", func() {
	steps := NewSteps()

	var (
		knownProductVersion = "0.0.21"
	)

	Scenario("Listing charts", func() {
		steps.When(fmt.Sprintf("running mkpcli chart list --product test-chart-product1 --product-version %s", knownProductVersion))
		steps.Then("the command exits without error")
		steps.And("the table of charts is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the table of charts is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("ID"))
			Eventually(CommandSession.Out).Should(Say("VERSION"))
			Eventually(CommandSession.Out).Should(Say("URL"))
			Eventually(CommandSession.Out).Should(Say("REPOSITORY"))

			Eventually(CommandSession.Out).Should(Say(knownProductVersion))
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("https://harbor-repo.vmware.com/chartrepo/tanzu_isv_engineering/charts/test-chart-product-%s.tgz", knownProductVersion)))
			Eventually(CommandSession.Out).Should(Say("https://harbor-repo.vmware.com/chartrepo/tanzu_isv_engineering tanzu_isv_engineering"))
		})
	})
})
