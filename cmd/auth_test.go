// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	. "github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/cmdfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/csp"
)

var _ = Describe("Auth", func() {
	Describe("GetRefreshToken", func() {
		var (
			initializer   *cmdfakes.FakeTokenServicesInitializer
			tokenServices *cmdfakes.FakeTokenServices
		)

		BeforeEach(func() {
			tokenServices = &cmdfakes.FakeTokenServices{}

			initializer = &cmdfakes.FakeTokenServicesInitializer{}
			initializer.Returns(tokenServices, nil)
			InitializeTokenServices = initializer.Spy
		})

		BeforeEach(func() {
			viper.Set("csp.api-token", "my-csp-api-token")
			viper.Set("csp.host", "console.cloud.vmware.com.example")
			tokenServices.RedeemReturns(&csp.Claims{
				Token: "my-refresh-token",
			}, nil)
		})

		It("gets the refresh token and puts it into viper", func() {
			err := GetRefreshToken(nil, []string{})
			Expect(err).ToNot(HaveOccurred())
			Expect(viper.GetString("csp.refresh-token")).To(Equal("my-refresh-token"))

			Expect(initializer.CallCount()).To(Equal(1))
			Expect(initializer.ArgsForCall(0)).To(Equal("https://console.cloud.vmware.com.example/"))

			Expect(tokenServices.RedeemCallCount()).To(Equal(1))
			Expect(tokenServices.RedeemArgsForCall(0)).To(Equal("my-csp-api-token"))
		})

		Context("fails to initialize token services", func() {
			BeforeEach(func() {
				initializer.Returns(nil, fmt.Errorf("initializer failed"))
			})

			It("returns an error", func() {
				err := GetRefreshToken(nil, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to initialize token services: initializer failed"))
			})
		})

		Context("fails to exchange api token", func() {
			BeforeEach(func() {
				tokenServices.RedeemReturns(nil, fmt.Errorf("redeem failed"))
			})

			It("returns an error", func() {
				err := GetRefreshToken(nil, []string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to exchange api token: redeem failed"))
			})
		})
	})
})
