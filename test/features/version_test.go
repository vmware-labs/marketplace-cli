// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package features_test

import (
	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Report version", func() {
	steps := NewSteps()

	Scenario("version command reports version", func() {
		steps.When("running mkpcli version")
		steps.Then("the command exits without error")
		steps.And("the version is printed")
	})

	Scenario("version command does not require CSP API token", func() {
		steps.Given("the environment variable CSP_API_TOKEN is not set")
		steps.When("running mkpcli version")
		steps.Then("the command exits without error")
		steps.And("the version is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the version is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("mkpcli version: 1.2.3"))
		})
	})
})
