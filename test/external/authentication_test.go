// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	"fmt"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("authentication", func() {
	steps := NewSteps()

	Scenario("Expired token", func() {
		steps.Given("targeting the production environment")
		steps.When(fmt.Sprintf("running mkpcli --csp-api-token %s auth", "M_sfojHArrjx90lxUCmID2qhZw-I0WGlW5fThBuiQXwVtvy7UJq6XeKtAKzf8cFm"))
		steps.Then("the command exits with an error")
		steps.And("the expired token error message is printed")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)
		define.Then(`^the expired token error message is printed$`, func() {
			Eventually(CommandSession.Err).Should(Say("the CSP API token is invalid or expired"))
		})
	})
})
