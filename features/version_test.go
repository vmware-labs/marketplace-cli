// +build feature

package features_test

import (
	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/features"
)

var _ = Describe("Report version", func() {
	steps := NewSteps()

	Scenario("version command reports version", func() {
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
