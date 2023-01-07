// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Download", func() {
	steps := NewSteps()

	Scenario("Downloading an asset", func() {
		steps.Given("targeting the staging environment")
		steps.When("running mkpcli download --product " + TKG + " --product-version " + TKGVersion + " --filter yq_linux_amd64 --filename yq --accept-eula")
		steps.Then("the command exits without error")
		steps.And("yq is downloaded")
	})

	Scenario("Download fails when there are multiple files", func() {
		steps.Given("targeting the staging environment")
		steps.When("running mkpcli download --product " + TKG + " --product-version " + TKGVersion + " --accept-eula")
		steps.Then("the command exits with an error")
		steps.And("a message saying that there are multiple assets available to download")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^a message saying that there are multiple assets available to download$`, func() {
			Eventually(CommandSession.Err).Should(Say("product " + TKG + " " + TKGVersion + " has multiple downloadable assets, please use the --filter parameter"))
		})
	})
})
