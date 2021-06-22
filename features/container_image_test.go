// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

// +build enemy

package features_test

import (
	"fmt"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/features"
)

var _ = Describe("Container Image", func() {
	steps := NewSteps()

	Scenario("Listing container images", func() {
		steps.When(fmt.Sprintf("running mkpcli container-image list --product %s --product-version %s", ContainerProductSlug, ContainerProductVersion))
		steps.Then("the command exits without error")
		steps.And("the table of container images is printed")
	})

	Scenario("Getting a single container images", func() {
		steps.When(fmt.Sprintf("running mkpcli container-image get --product %s --product-version %s --image-repository harbor-repo.vmware.com/tanzu_isv_engineering/test-container-product", ContainerProductSlug, ContainerProductVersion))
		steps.Then("the command exits without error")
		steps.And("the table of the container image is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the table of container images is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("IMAGE"))
			Eventually(CommandSession.Out).Should(Say("TAGS"))
			Eventually(CommandSession.Out).Should(Say(fmt.Sprintf("harbor-repo.vmware.com/tanzu_isv_engineering/test-container-product *%s", ContainerProductVersion)))
			Eventually(CommandSession.Out).Should(Say("Deployment instructions:"))
		})

		define.Then(`^the table of the container image is printed$`, func() {
			Eventually(CommandSession.Out).Should(Say("TAG"))
			Eventually(CommandSession.Out).Should(Say("TYPE"))
			Eventually(CommandSession.Out).Should(Say(ContainerProductVersion))
			Eventually(CommandSession.Out).Should(Say("FIXED"))
		})
	})
})
