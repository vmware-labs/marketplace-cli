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

var _ = Describe("Product Version", func() {
	steps := NewSteps()

	var (
		knownProductVersion        = "0.0.35"
		anotherKnownProductVersion = "0.0.36"
	)

	Scenario("Listing product versions", func() {
		steps.When("running mkpcli product-version list --product test-container-product2")
		steps.Then("the command exits without error")
		steps.And("the table of product versions is printed")
	})

	Scenario("Getting a single product version", func() {
		steps.When(fmt.Sprintf("running mkpcli product-version get --product test-container-product2 --product-version %s", knownProductVersion))
		steps.Then("the command exits without error")
		steps.And("the table of the product version is printed")
	})

	steps.Define(func(define Definitions) {

		DefineCommonSteps(define)

		define.Then(`^the table of product versions is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Versions:"))
			Eventually(CommandSession.Out).Should(Say("NUMBER"))
			Eventually(CommandSession.Out).Should(Say("STATUS"))

			Eventually(CommandSession.Out).Should(Say(knownProductVersion))
			Eventually(CommandSession.Out).Should(Say(anotherKnownProductVersion))
			Eventually(CommandSession.Out).Should(Say("PENDING"))

			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("Version %s:", knownProductVersion)))
			Eventually(CommandSession.Out).Should(Say("IMAGE"))
			Eventually(CommandSession.Out).Should(Say("TAGS"))
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("harbor-repo.vmware.com/tanzu_isv_engineering/test-container-product *%s", knownProductVersion)))
			Eventually(CommandSession.Out).Should(Say("Deployment instructions:"))

			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("Version %s:", anotherKnownProductVersion)))
			Eventually(CommandSession.Out).Should(Say("IMAGE"))
			Eventually(CommandSession.Out).Should(Say("TAGS"))
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("harbor-repo.vmware.com/tanzu_isv_engineering/test-container-product *%s", anotherKnownProductVersion)))
			Eventually(CommandSession.Out).Should(Say("Deployment instructions:"))
		})

		define.Then(`^the table of the product version is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("Version %s:", knownProductVersion)))
			Eventually(CommandSession.Out).Should(Say("IMAGE"))
			Eventually(CommandSession.Out).Should(Say("TAGS"))
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("harbor-repo.vmware.com/tanzu_isv_engineering/test-container-product *%s", knownProductVersion)))
			Eventually(CommandSession.Out).Should(Say("Deployment instructions:"))
		})
	})
})
