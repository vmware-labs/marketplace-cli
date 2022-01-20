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
		steps.When("running mkpcli --debug product get --product vmware-tanzu-rabbitmq1")
		steps.Then("the command exits without error")
		steps.And("the request is printed")
	})

	Scenario("Debugging enabled with environment variable", func() {
		steps.When("The environment variable MKPCLI_DEBUG is set to true")
		steps.And("running mkpcli product get --product vmware-tanzu-rabbitmq1")
		steps.Then("the command exits without error")
		steps.And("the request is printed")
	})

	Scenario("Debugging enabled with request payloads", func() {
		steps.When("running mkpcli --debug --debug-request-payloads download -p vmware-tanzu-rabbitmq1 -v 1.0.0")
		steps.Then("the command exits without error")
		steps.And("the container image is downloaded")
		steps.And("the requests are printed with request payloads")
	})

	Scenario("Debugging enabled with request payloads with environment variables", func() {
		steps.When("The environment variable MKPCLI_DEBUG is set to true")
		steps.And("The environment variable MKPCLI_DEBUG_REQUEST_PAYLOADS is set to true")
		steps.And("running mkpcli download -p vmware-tanzu-rabbitmq1 -v 1.0.0")
		steps.Then("the command exits without error")
		steps.And("the container image is downloaded")
		steps.And("the requests are printed with request payloads")
	})

	steps.Define(func(define Definitions) {
		DefineCommonSteps(define, "production")

		define.Then(`^the request is printed$`, func() {
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #0: GET https://gtw.marketplace.cloud.vmware.com/api/v1/products/vmware-tanzu-rabbitmq1?increaseViewCount=false&isSlug=true")))
			Eventually(CommandSession.Err).Should(Say("Request #0 Response: 200 OK"))
			Eventually(CommandSession.Out).Should(Say("Name:      VMware Tanzu RabbitMQ"))
			Eventually(CommandSession.Out).Should(Say("Publisher: VMware Inc"))
			Eventually(CommandSession.Out).Should(Say("Assets for 1.0.0:"))
			Eventually(CommandSession.Out).Should(Say("registry.pivotal.io/rabbitmq/vmware-tanzu-rabbitmq:2020.12"))
		})

		define.Then(`^the container image is downloaded$`, func() {
			err := os.Remove("image.tar")
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then(`^the requests are printed with request payloads$`, func() {
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #0: GET https://gtw.marketplace.cloud.vmware.com/api/v1/products/vmware-tanzu-rabbitmq1?increaseViewCount=false&isSlug=true")))
			Eventually(CommandSession.Err).Should(Say("Request #0 Response: 200 OK"))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #1: POST https://gtw.marketplace.cloud.vmware.com/api/v1/products/5f99a9d5-dbfd-4cfc-a564-b1a67d092b4f/download")))
			Eventually(CommandSession.Err).Should(Say("--- Start of request #1 body payload ---"))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("{\"productid\":\"5f99a9d5-dbfd-4cfc-a564-b1a67d092b4f\",\"appVersion\":\"1.0.0\",\"eulaAccepted\":true,\"dockerlinkVersionId\":\"d333021c-6e7d-4a15-b87a-2f66eda9d30c\",\"dockerUrlId\":\"f1702bcd-5634-4983-a652-653b6aedbe1d\",\"imageTagId\":\"83fa36ea-4ebc-4499-b68d-b95bc185dd65\"}")))
			Eventually(CommandSession.Err).Should(Say("--- End of request #1 body payload ---"))
			Eventually(CommandSession.Err).Should(Say("Request #1 Response: 200 OK"))
			Eventually(CommandSession.Err).Should(Say(regexp.QuoteMeta("Request #2: GET https://cmpprdcontainersolutions.s3.us-west-2.amazonaws.com/containerImageTars/")))
			Eventually(CommandSession.Err).Should(Say("Request #2 Response: 200 OK"))
		})
	})
})
