// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"errors"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/internalfakes"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
)

var _ = Describe("Upload", func() {
	var (
		httpClient  *pkgfakes.FakeHTTPClient
		marketplace *pkg.Marketplace
	)

	BeforeEach(func() {
		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			APIHost: "marketplace.api.example.com",
			Host:    "marketplace.example.com",
			Client:  httpClient,
		}
	})

	Describe("GetUploadCredentials", func() {
		BeforeEach(func() {
			response := &pkg.CredentialsResponse{
				AccessID:     "my-access-id",
				AccessKey:    "my-access-key",
				SessionToken: "my-session-token",
				Expiration:   time.Time{},
			}
			httpClient.DoReturns(MakeJSONResponse(response), nil)
		})

		It("gets the credentials", func() {
			creds, err := marketplace.GetUploadCredentials()
			Expect(err).ToNot(HaveOccurred())
			Expect(creds.AccessID).To(Equal("my-access-id"))
			Expect(creds.AccessKey).To(Equal("my-access-key"))
			Expect(creds.SessionToken).To(Equal("my-session-token"))

			By("requesting the creds from the Marketplace", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.String()).To(Equal("https://marketplace.api.example.com/aws/credentials/generate"))
			})
		})

		When("the credentials request fails", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("get credentials failed"))
			})
			It("returns an error", func() {
				_, err := marketplace.GetUploadCredentials()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("marketplace request failed: get credentials failed"))
			})
		})

		When("the credentials response is not 200 OK", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusTeapot,
				}, nil)
			})

			It("returns an error", func() {
				_, err := marketplace.GetUploadCredentials()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to fetch credentials: 418"))
			})
		})

		When("the credentials response is invalid", func() {
			BeforeEach(func() {
				httpClient.DoReturns(MakeStringResponse("this is not valid json"), nil)
			})
			It("returns an error", func() {
				_, err := marketplace.GetUploadCredentials()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid character 'h' in literal true (expecting 'r')"))
			})
		})
	})

	Describe("GetUploader", func() {
		BeforeEach(func() {
			response := &pkg.CredentialsResponse{
				AccessID:     "my-access-id",
				AccessKey:    "my-access-key",
				SessionToken: "my-session-token",
				Expiration:   time.Time{},
			}
			httpClient.DoReturns(MakeJSONResponse(response), nil)
		})
		It("creates an uploader with upload credentials", func() {
			uploader, err := marketplace.GetUploader("my-org")
			Expect(err).ToNot(HaveOccurred())

			By("requesting the upload credentials", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
			})
			Expect(uploader).ToNot(BeNil())
		})

		When("getting the credentials fails", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, errors.New("get credentials failed"))
			})
			It("returns an error", func() {
				_, err := marketplace.GetUploader("my-org")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to get upload credentials: marketplace request failed: get credentials failed"))
			})
		})

		When("the uploaded already exists", func() {
			var uploader internal.Uploader
			BeforeEach(func() {
				uploader = &internalfakes.FakeUploader{}
				marketplace.SetUploader(uploader)
			})
			It("returns that uploader", func() {
				Expect(marketplace.GetUploader("doesn't matter")).To(Equal(uploader))
				Expect(httpClient.DoCallCount()).To(Equal(0))
			})
		})
	})
})
