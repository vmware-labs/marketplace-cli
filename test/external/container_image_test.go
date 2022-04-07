// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Container Image", func() {
	steps := NewSteps()

	Scenario("Listing container images", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli container-image list --product redis-enterprise-kubernetes-operator-for-vmware-enterprise-pks --product-version 5.4.2-27")
		steps.Then("the command exits without error")
		steps.And("the table of container images is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the table of container images is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("IMAGE"))
			Eventually(CommandSession.Out).Should(Say("TAGS"))
			Eventually(CommandSession.Out).Should(Say("DOWNLOADS"))
			Eventually(CommandSession.Out).Should(Say("https://hub.docker.com/u/redislabs"))
		})
	})
})
