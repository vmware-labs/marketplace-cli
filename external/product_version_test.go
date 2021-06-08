// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

// +build external

package external_test

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

	Scenario("Listing product versions", func() {
		steps.When(fmt.Sprintf("running mkpcli product-version list --product %s", ContainerProductSlug))
		steps.Then("the command exits without error")
		steps.And("the table of product versions is printed")
	})

	Scenario("Getting a single product version", func() {
		steps.When(fmt.Sprintf("running mkpcli product-version get --product %s --product-version %s", ContainerProductSlug, ContainerProductVersion))
		steps.Then("the command exits without error")
		steps.And("the table of the product version is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the table of product versions is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Versions:"))
			Eventually(CommandSession.Out).Should(Say("NUMBER"))
			Eventually(CommandSession.Out).Should(Say("STATUS"))

			Eventually(CommandSession.Out).Should(Say(ContainerProductVersion))
			Eventually(CommandSession.Out).Should(Say("PENDING"))
			Eventually(CommandSession.Out).Should(Say(ContainerProductVersionButOlder))
			Eventually(CommandSession.Out).Should(Say("PENDING"))
		})

		define.Then(`^the table of the product version is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("Version %s:", ContainerProductVersion)))
			Eventually(CommandSession.Out).Should(Say("IMAGE"))
			Eventually(CommandSession.Out).Should(Say("TAGS"))
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("harbor-repo.vmware.com/tanzu_isv_engineering/test-container-product *%s", ContainerProductVersion)))
			Eventually(CommandSession.Out).Should(Say("Deployment instructions:"))
		})
	})
})
