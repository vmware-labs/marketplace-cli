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

var _ = Describe("product list-assets", func() {
	const (
		ChartProductSlug             = "nginx"
		ChartProductVersion          = "1.21.1_0"
		ContainerImageProductSlug    = "cloudian-s3-compatible-object-storage-for-tkgs0-1-1"
		ContainerImageProductVersion = "1.2.1"
		VMProductSlug                = "nginxstack"
		VMProductVersion             = "1.21.0_1"
	)

	steps := NewSteps()

	Scenario("Listing charts", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli product list-assets --type chart --product " + ChartProductSlug + " --product-version " + ChartProductVersion)
		steps.Then("the command exits without error")
		steps.And("the table of charts is printed")
	})

	Scenario("Listing container images", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli product list-assets --type image --product " + ContainerImageProductSlug + " --product-version " + ContainerImageProductVersion)
		steps.Then("the command exits without error")
		steps.And("the table of container images is printed")
	})

	Scenario("Listing virtual machine files", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli product list-assets --type vm --product " + VMProductSlug + " --product-version " + VMProductVersion)
		steps.Then("the command exits without error")
		steps.And("the table of virtual machine files is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the table of charts is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Chart assets for NGINX Open Source Helm Chart packaged by Bitnami 1.21.1_0:"))
			Eventually(CommandSession.Out).Should(Say("NAME"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("VERSION"))
			Eventually(CommandSession.Out).Should(Say("SIZE"))
			Eventually(CommandSession.Out).Should(Say("DOWNLOADS"))

			Eventually(CommandSession.Out).Should(Say("https://charts.bitnami.com/bitnami/nginx-9.3.6.tgz"))
			Eventually(CommandSession.Out).Should(Say("Chart"))
			Eventually(CommandSession.Out).Should(Say("9.3.6"))
			Eventually(CommandSession.Out).Should(Say("0 B"))
		})

		define.Then(`^the table of container images is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("Container Image assets for Cloudian S3 compatible object storage for Tanzu 1.2.1:"))
			Eventually(CommandSession.Out).Should(Say("NAME"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("VERSION"))
			Eventually(CommandSession.Out).Should(Say("SIZE"))
			Eventually(CommandSession.Out).Should(Say("DOWNLOADS"))

			Eventually(CommandSession.Out).Should(Say("quay.io/cloudian/hyperstorec:v1.3.0rc1"))
			Eventually(CommandSession.Out).Should(Say("Container Image"))
			Eventually(CommandSession.Out).Should(Say("v1.3.0rc1"))
			Eventually(CommandSession.Out).Should(Say("2.38 GB"))
		})

		define.Then(`^the table of virtual machine files is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("VM assets for NGINX Open Source Virtual Appliance packaged by Bitnami 1.21.0_1:"))
			Eventually(CommandSession.Out).Should(Say("NAME"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("VERSION"))
			Eventually(CommandSession.Out).Should(Say("SIZE"))
			Eventually(CommandSession.Out).Should(Say("DOWNLOADS"))

			Eventually(CommandSession.Out).Should(Say("nginxstack"))
			Eventually(CommandSession.Out).Should(Say("VM"))
			Eventually(CommandSession.Out).Should(Say("1.28 GB"))
		})
	})
})
