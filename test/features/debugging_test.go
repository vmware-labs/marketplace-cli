// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package features_test

import (
	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Debugging", func() {
	steps := NewSteps()

	Scenario("No debugging", func() {
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has debugging.enabled with the value false")
		steps.And("the printed configuration has debugging.print-request-payloads with the value false")
	})

	Scenario("Debug flag", func() {
		steps.When("running mkpcli config --debug")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has debugging.enabled with the value true")
		steps.And("the printed configuration has debugging.print-request-payloads with the value false")
	})

	Scenario("Debug environment variable", func() {
		steps.Given("the environment variable MKPCLI_DEBUG is set to true")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has debugging.enabled with the value true")
		steps.And("the printed configuration has debugging.print-request-payloads with the value false")
	})

	Scenario("Debug flag", func() {
		steps.When("running mkpcli config --debug --debug-request-payloads")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has debugging.enabled with the value true")
		steps.And("the printed configuration has debugging.print-request-payloads with the value true")
	})

	Scenario("Debug environment variable", func() {
		steps.Given("the environment variable MKPCLI_DEBUG is set to true")
		steps.And("the environment variable MKPCLI_DEBUG_REQUEST_PAYLOADS is set to true")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has debugging.enabled with the value true")
		steps.And("the printed configuration has debugging.print-request-payloads with the value true")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the version is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("mkpcli version: 1.2.3"))
		})
	})
})
