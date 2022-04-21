// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	"os"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Download", func() {
	const (
		ProductSlug    = "tanzu-kubenetes-grid-1-111-1-1"
		ProductVersion = "1.5.1"
	)

	steps := NewSteps()

	Scenario("Downloading an asset", func() {
		steps.Given("targeting the staging environment")
		steps.When("running mkpcli download --product " + ProductSlug + " --product-version " + ProductVersion + " --filter yq_linux_amd64 --filename yq --accept-eula")
		steps.Then("the command exits without error")
		steps.And("yq is downloaded")
	})

	Scenario("Download fails when there are multiple files", func() {
		steps.Given("targeting the staging environment")
		steps.When("running mkpcli download --product " + ProductSlug + " --product-version " + ProductVersion + " --accept-eula")
		steps.Then("the command exits with an error")
		steps.And("a message saying that there are multiple assets available to download")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^yq is downloaded$`, func() {
			_, err := os.Stat("yq")
			Expect(err).ToNot(HaveOccurred())
		}, func() {
			_ = os.Remove("yq")
		})

		define.Then(`^a message saying that there are multiple assets available to download$`, func() {
			Eventually(CommandSession.Err).Should(Say("product " + ProductSlug + " " + ProductVersion + " has multiple downloadable assets, please use the --filter parameter"))
		})
	})
})
