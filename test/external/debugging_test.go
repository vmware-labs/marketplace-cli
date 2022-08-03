// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	"os"
	"regexp"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Debugging", func() {
	steps := NewSteps()

	Scenario("Debugging enabled", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli --debug product get --product nginx --product-version 1.22.0_150_r04")
		steps.Then("the command exits without error")
		steps.And("the request is printed")
	})

	Scenario("Debugging enabled with environment variable", func() {
		steps.Given("targeting the production environment")
		steps.And("the environment variable MKPCLI_DEBUG is set to true")
		steps.When("running mkpcli product get --product nginx --product-version 1.22.0_150_r04")
		steps.Then("the command exits without error")
		steps.And("the request is printed")
	})

	Scenario("Debugging enabled with request payloads", func() {
		steps.Given("targeting the production environment")
		steps.When("running mkpcli --debug --debug-request-payloads download -p nginx -v 1.22.0_150_r04 --filename chart.tgz --accept-eula")
		steps.Then("the command exits without error")
		steps.And("chart.tgz is downloaded")
		steps.And("the requests are printed with request payloads")
	})

	Scenario("Debugging enabled with request payloads with environment variables", func() {
		steps.Given("targeting the production environment")
		steps.And("the environment variable MKPCLI_DEBUG is set to true")
		steps.And("the environment variable MKPCLI_DEBUG_REQUEST_PAYLOADS is set to true")
		steps.When("running mkpcli download -p nginx -v 1.22.0_150_r04 --filename chart.tgz --accept-eula")
		steps.Then("the command exits without error")
		steps.And("chart.tgz is downloaded")
		steps.And("the requests are printed with request payloads")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define)

		define.Then(`^the request is printed$`, func() {
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #0: POST https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #0 Response: 200 OK")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #1: GET https://console.cloud.vmware.com/csp/gateway/am/api/auth/token-public-key")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #1 Response: 200 OK")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #2: GET https://gtw.marketplace.cloud.vmware.com/api/v1/products/nginx?increaseViewCount=false&isSlug=true")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #2 Response: 200 OK")))
			Eventually(CommandSession.Out).Should(Say(regexp.QuoteMeta("Name:      NGINX Open Source Helm Chart packaged by Bitnami")))
			Eventually(CommandSession.Out).Should(Say(regexp.QuoteMeta("Publisher: Bitnami")))
			Eventually(CommandSession.Out).Should(Say(regexp.QuoteMeta("Assets for 1.22.0_150_r04:")))
			Eventually(CommandSession.Out).Should(Say(regexp.QuoteMeta("https://charts.bitnami.com/bitnami/nginx-12.0.4.tgz")))
		})

		define.Then(`^the container image is downloaded$`, func() {
			err := os.Remove("image.tar")
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then(`^the requests are printed with request payloads$`, func() {
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #0: POST https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("--- Start of request #0 body payload ---")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("refresh_token=")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("--- End of request #0 body payload ---")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #0 Response: 200 OK")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #1: GET https://console.cloud.vmware.com/csp/gateway/am/api/auth/token-public-key")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #1 Response: 200 OK")))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #2: GET https://gtw.marketplace.cloud.vmware.com/api/v1/products/nginx?increaseViewCount=false&isSlug=true")))
			Eventually(CommandSession.Err).Should(Say("Request #2 Response: 200 OK"))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #3: POST https://gtw.marketplace.cloud.vmware.com/api/v1/products/89431c5d-ddb7-45df-a544-2c81a370e17b/version-details?versionNumber=1.22.0_150_r04")))
			Eventually(CommandSession.Err).Should(Say("Request #3 Response: 200 OK"))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #4: POST https://gtw.marketplace.cloud.vmware.com/api/v1/products/89431c5d-ddb7-45df-a544-2c81a370e17b/download")))
			Eventually(CommandSession.Err).Should(Say("--- Start of request #4 body payload ---"))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("{\"productid\":\"89431c5d-ddb7-45df-a544-2c81a370e17b\",\"appVersion\":\"1.22.0_150_r04\",\"eulaAccepted\":true,\"chartVersion\":\"12.0.4\"}")))
			Eventually(CommandSession.Err).Should(Say("--- End of request #4 body payload ---"))
			Eventually(CommandSession.Err).Should(Say("Request #4 Response: 200 OK"))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #5: GET https://charts.bitnami.com/bitnami/nginx-13.1.5.tgz")))
			Eventually(CommandSession.Err).Should(Say("Request #5 Response: 200 OK"))
		})
	})
})
