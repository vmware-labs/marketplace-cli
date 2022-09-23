// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package external_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	. "github.com/bunniesandbeatings/goerkin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("retryable requests", func() {
	steps := NewSteps()

	Scenario("Does not fail", func() {
		steps.Given("targeting the staging environment")
		steps.And("using a csp proxy with 0 failures")
		steps.And("mkpcli will use the csp proxy instead of the default")
		steps.And("the environment variable MKPCLI_SKIP_SSL_VALIDATION is set to true")
		steps.When("running mkpcli --debug product list")
		steps.Then("the command exits without error")
	})

	Scenario("Fails but retries will succeed", func() {
		steps.Given("targeting the staging environment")
		steps.And("using a csp proxy with 2 failures")
		steps.And("mkpcli will use the csp proxy instead of the default")
		steps.And("the environment variable MKPCLI_SKIP_SSL_VALIDATION is set to true")
		steps.When("running mkpcli --debug product list")
		steps.Then("the command exits without error")
	})

	Scenario("Fails too many times", func() {
		steps.Given("targeting the staging environment")
		steps.And("using a csp proxy with 6 failures")
		steps.And("mkpcli will use the csp proxy instead of the default")
		steps.And("the environment variable MKPCLI_SKIP_SSL_VALIDATION is set to true")
		steps.When("running mkpcli product list")
		steps.Then("the command exits with an error")
	})

	steps.Define(func(define Definitions) {
		var cspProxyServer *httptest.Server
		DefineCommonSteps(define)

		define.Given(`^using a csp proxy with (\d) failures$`, func(totalFailureCountString string) {
			failureCount := 0
			totalFailureCount, err := strconv.Atoi(totalFailureCountString)
			Expect(err).ToNot(HaveOccurred())

			cspProxyServer = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if failureCount < totalFailureCount {
					w.WriteHeader(http.StatusServiceUnavailable)
					failureCount++
					return
				}

				url := fmt.Sprintf("https://console.cloud.vmware.com/%s?%s", r.URL.Path, r.URL.RawQuery)
				body, err := io.ReadAll(r.Body)
				Expect(err).ToNot(HaveOccurred())
				forwardedRequest, err := http.NewRequest(r.Method, url, io.NopCloser(bytes.NewReader(body)))
				Expect(err).ToNot(HaveOccurred())
				copyHeaders(r.Header, forwardedRequest.Header)

				resp, err := http.DefaultClient.Do(forwardedRequest)
				Expect(err).ToNot(HaveOccurred())

				// Prepare the response
				copyHeaders(resp.Header, w.Header())
				w.WriteHeader(resp.StatusCode)

				_, err = io.Copy(w, resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Body.Close()).To(Succeed())
			}))
			cspProxyServer.StartTLS()
		}, func() {
			cspProxyServer.Close()
		})

		define.Given(`^mkpcli will use the csp proxy instead of the default$`, func() {
			cpsProxyServerUrl, err := url.Parse(cspProxyServer.URL)
			Expect(err).ToNot(HaveOccurred())
			EnvVars = append(EnvVars, "CSP_HOST="+cpsProxyServerUrl.Host)
		})

		define.Then(`^the expired token error message is printed$`, func() {
			Eventually(CommandSession.Err).Should(Say("the CSP API token is invalid or expired"))
		})
	})
})

func copyHeaders(source, dest http.Header) {
	for headerKey, headerValues := range source {
		for _, headerValue := range headerValues {
			dest.Add(headerKey, headerValue)
		}
	}
}
