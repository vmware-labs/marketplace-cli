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

var _ = Describe("Virtual Machine", func() {
	steps := NewSteps()

	Scenario("Listing virtual machine files", func() {
		steps.When("running mkpcli vm list --product nginxstack --product-version 1.21.0_1")
		steps.Then("the command exits without error")
		steps.And("the table of virtual machine files is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define, "production")

		define.Then(`^the table of virtual machine files is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("ID"))
			Eventually(CommandSession.Out).Should(Say("NAME"))
			Eventually(CommandSession.Out).Should(Say("STATUS"))
			Eventually(CommandSession.Out).Should(Say("SIZE"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say("FILES"))
			Eventually(CommandSession.Out).Should(Say("68c1663c-4f1e-4a8e-a719-0aef7d6bae94"))
			Eventually(CommandSession.Out).Should(Say("bitnami-nginx-1.21.0-1-linux-centos-7-x86_64-nami"))
			Eventually(CommandSession.Out).Should(Say("ACTIVE"))
			Eventually(CommandSession.Out).Should(Say("1.28 GB"))
			Eventually(CommandSession.Out).Should(Say("vcsp.ovf"))
			Eventually(CommandSession.Out).Should(Say("3"))
		})
	})
})
