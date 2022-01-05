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
		steps.When("running mkpcli container-image list --product redis-enterprise-kubernetes-operator-for-vmware-enterprise-pks --product-version 5.4.2-27")
		steps.Then("the command exits without error")
		steps.And("the table of container images is printed")
	})

	Scenario("Getting a single container images", func() {
		steps.When("running mkpcli container-image get --product redis-enterprise-kubernetes-operator-for-vmware-enterprise-pks --product-version 5.4.2-27 --image-repository https://hub.docker.com/u/redislabs")
		steps.Then("the command exits without error")
		steps.And("the table of the container image is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define, "production")

		define.Then(`^the table of container images is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("IMAGE"))
			Eventually(CommandSession.Out).Should(Say("TAGS"))
			Eventually(CommandSession.Out).Should(Say("DOWNLOADS"))
			Eventually(CommandSession.Out).Should(Say("https://hub.docker.com/u/redislabs"))

			Eventually(CommandSession.Out).Should(Say("Deployment instructions:"))
			Eventually(CommandSession.Out).Should(Say("Redis Enterprise for PKS is deployed and maintained using a Kubernetes Operator."))
		})

		define.Then(`^the table of the container image is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("https://hub.docker.com/u/redislabs"))
			Eventually(CommandSession.Out).Should(Say("Tags:"))
			Eventually(CommandSession.Out).Should(Say("TAG"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("DOWNLOADS"))

			Eventually(CommandSession.Out).Should(Say("Deployment instructions:"))
			Eventually(CommandSession.Out).Should(Say("Redis Enterprise for PKS is deployed and maintained using a Kubernetes Operator."))
		})
	})
})
