// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package csp_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal/csp"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("CSP Token Services", func() {
	var (
		client        *pkgfakes.FakeHTTPClient
		tokenServices *csp.TokenServices
	)

	BeforeEach(func() {
		client = &pkgfakes.FakeHTTPClient{}
		tokenServices = &csp.TokenServices{
			CSPHost: "csp.example.com",
			Client:  client,
		}
	})

	Describe("Redeem", func() {
		XIt("exchanges the token for the JWT token claims", func() {

		})

		When("sending the request fails", func() {
			BeforeEach(func() {
				client.PostFormReturns(nil, errors.New("failed to send request"))
			})
			It("returns an error", func() {
				_, err := tokenServices.Redeem("my-csp-api-token")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to redeem token: failed to send request"))
			})
		})

		When("the request returns an unparseable response", func() {
			BeforeEach(func() {
				client.PostFormReturns(test.MakeFailingBodyResponse("bad-response-body"), nil)
			})
			It("returns an error", func() {
				_, err := tokenServices.Redeem("my-csp-api-token")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse redeem response: bad-response-body"))
			})
		})

		When("the response indicates that the token is invalid", func() {
			BeforeEach(func() {
				responseBody := csp.RedeemResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "invalid_grant: Invalid refresh token: xxxxx-token",
				}
				response := test.MakeJSONResponse(responseBody)
				response.StatusCode = http.StatusBadRequest
				client.PostFormReturns(response, nil)
			})
			It("returns an error", func() {
				_, err := tokenServices.Redeem("my-csp-api-token")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("the CSP API token is invalid or expired"))
			})
		})

		When("the response shows some other error", func() {
			BeforeEach(func() {
				responseBody := csp.RedeemResponse{
					StatusCode: http.StatusTeapot,
					Message:    "teapots!",
				}
				response := test.MakeJSONResponse(responseBody)
				response.Status = http.StatusText(http.StatusTeapot)
				response.StatusCode = http.StatusTeapot
				client.PostFormReturns(response, nil)
			})
			It("returns an error", func() {
				_, err := tokenServices.Redeem("my-csp-api-token")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to exchange refresh token for access token: I'm a teapot: teapots!"))
			})
		})

		XWhen("the response is not a valid token", func() {

		})
	})
})
