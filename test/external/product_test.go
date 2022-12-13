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

var _ = Describe("Product", func() {
	steps := NewSteps()

	Scenario("Listing products", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli product list --all-orgs --search-text " + Nginx)
		steps.Then("the command exits without error")
		steps.And("the list of products is printed")
	})

	Scenario("Listing product versions", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli product list-versions --product " + Nginx)
		steps.Then("the command exits without error")
		steps.And("the list of product versions is printed")
	})

	Scenario("Getting product details", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli product get --product " + Nginx)
		steps.Then("the command exits without error")
		steps.And("the product details are printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the list of products is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("All products from all organizations filtered by \"nginx\""))
			Eventually(CommandSession.Out).Should(Say("SLUG"))
			Eventually(CommandSession.Out).Should(Say("NAME"))
			Eventually(CommandSession.Out).Should(Say("PUBLISHER"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("LATEST VERSION"))

			Eventually(CommandSession.Out).Should(Say("nginx"))
			Eventually(CommandSession.Out).Should(Say("NGINX Open Source Helm Chart packaged by Bitnami"))
			Eventually(CommandSession.Out).Should(Say("Bitnami"))
			Eventually(CommandSession.Out).Should(Say("HELMCHARTS"))

			Eventually(CommandSession.Out).Should(Say(`Total count: \d`))
		})

		define.Then(`^the list of product versions is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Versions for NGINX Open Source Helm Chart packaged by Bitnami:"))
			Eventually(CommandSession.Out).Should(Say("NUMBER"))
			Eventually(CommandSession.Out).Should(Say("STATUS"))

			Eventually(CommandSession.Out).Should(Say("1.21.3_0"))
			Eventually(CommandSession.Out).Should(Say("ACTIVE"))
			Eventually(CommandSession.Out).Should(Say("1.21.2_0"))
			Eventually(CommandSession.Out).Should(Say("ACTIVE"))
			Eventually(CommandSession.Out).Should(Say("1.21.1_0"))
			Eventually(CommandSession.Out).Should(Say("ACTIVE"))
		})

		define.Then(`^the product details are printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Name:"))
			Eventually(CommandSession.Out).Should(Say("NGINX Open Source Helm Chart packaged by Bitnami"))
			Eventually(CommandSession.Out).Should(Say("Publisher:"))
			Eventually(CommandSession.Out).Should(Say("Bitnami"))

			Eventually(CommandSession.Out).Should(Say("Up-to-date, secure, and ready to deploy on Kubernetes."))
			Eventually(CommandSession.Out).Should(Say(`https://marketplace.cloud.vmware.com/services/details/nginx\?slug=true`))

			Eventually(CommandSession.Out).Should(Say("Product Details:"))
			Eventually(CommandSession.Out).Should(Say("PRODUCT ID"))
			Eventually(CommandSession.Out).Should(Say("SLUG"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("LATEST VERSION"))
			Eventually(CommandSession.Out).Should(Say("89431c5d-ddb7-45df-a544-2c81a370e17b"))
			Eventually(CommandSession.Out).Should(Say("nginx"))
			Eventually(CommandSession.Out).Should(Say("HELMCHARTS"))

			Eventually(CommandSession.Out).Should(Say("Description:"))
			Eventually(CommandSession.Out).Should(Say("NGINX Open Source is a lightweight and high-performance server"))
		})
	})
})
