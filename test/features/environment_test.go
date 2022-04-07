// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package features_test

import (
	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Marketplace environment variables", func() {
	steps := NewSteps()

	Scenario("Production environment", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has marketplace.host with the value gtw.marketplace.cloud.vmware.com")
		steps.And("the printed configuration has marketplace.api-host with the value api.marketplace.cloud.vmware.com")
		steps.And("the printed configuration has marketplace.ui-host with the value marketplace.cloud.vmware.com")
		steps.And("the printed configuration has marketplace.storage.bucket with the value cspmarketplaceprd")
		steps.And("the printed configuration has marketplace.storage.region with the value us-west-2")
	})

	Scenario("Staging environment", func() {
		steps.Given("targeting the staging environment")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has marketplace.host with the value gtwstg.market.csp.vmware.com")
		steps.And("the printed configuration has marketplace.api-host with the value apistg.market.csp.vmware.com")
		steps.And("the printed configuration has marketplace.ui-host with the value stg.market.csp.vmware.com")
		steps.And("the printed configuration has marketplace.storage.bucket with the value cspmarketplacestage")
		steps.And("the printed configuration has marketplace.storage.region with the value us-east-2")
	})

	Scenario("Overriding marketplace gateway host", func() {
		steps.Given("the environment variable MKPCLI_HOST is set to gtw.marketplace.example.com")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has marketplace.host with the value gtw.marketplace.example.com")
	})

	Scenario("Overriding marketplace API host", func() {
		steps.Given("the environment variable MKPCLI_API_HOST is set to api.marketplace.example.com")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has marketplace.api-host with the value api.marketplace.example.com")
	})

	Scenario("Overriding marketplace UI host", func() {
		steps.Given("the environment variable MKPCLI_UI_HOST is set to marketplace.example.com")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has marketplace.ui-host with the value marketplace.example.com")
	})

	Scenario("Overriding marketplace storage bucket", func() {
		steps.Given("the environment variable MKPCLI_STORAGE_BUCKET is set to exmaplestoragebucket")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has marketplace.storage.bucket with the value exmaplestoragebucket")
	})

	Scenario("Overriding marketplace storage region", func() {
		steps.Given("the environment variable MKPCLI_STORAGE_REGION is set to us-central-1")
		steps.When("running mkpcli config")
		steps.Then("the command exits without error")
		steps.And("the printed configuration has marketplace.storage.region with the value us-central-1")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)
	})
})
