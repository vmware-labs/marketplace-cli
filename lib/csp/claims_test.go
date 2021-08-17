// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package csp_test

import (
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/lib/csp"
)

var _ = Describe("CSP Claims", func() {
	Describe("GetQualifiedUsername", func() {
		It("returns the username", func() {
			claims := &csp.Claims{
				Username: "john@yahoo.com",
			}

			Expect(claims.GetQualifiedUsername()).To(Equal("john@yahoo.com"))
		})

		Context("No domain in username", func() {
			It("returns a username using the domain field", func() {
				claims := &csp.Claims{
					Domain:   "example.com",
					Username: "john",
				}

				Expect(claims.GetQualifiedUsername()).To(Equal("john@example.com"))
			})
		})
	})

	Describe("IsOrgOwner", func() {
		It("returns true if the right perm is set", func() {
			claims := &csp.Claims{
				Perms: []string{
					"csp:org_owner",
					"other:perm",
				},
			}
			Expect(claims.IsOrgOwner()).To(BeTrue())

			claims = &csp.Claims{
				Perms: []string{
					"other:perm",
				},
			}
			Expect(claims.IsOrgOwner()).To(BeFalse())
		})
	})

	Describe("IsOrgOwner", func() {
		It("returns true if the right perm is set", func() {
			claims := &csp.Claims{
				Perms: []string{
					"csp:org_owner",
					"other:perm",
				},
			}
			Expect(claims.IsOrgOwner()).To(BeTrue())

			claims = &csp.Claims{
				Perms: []string{
					"other:perm",
				},
			}
			Expect(claims.IsOrgOwner()).To(BeFalse())
		})
	})

	Describe("IsPlatformOperator", func() {
		It("returns true if the right perm is set", func() {
			claims := &csp.Claims{
				Perms: []string{
					"csp:platform_operator",
					"other:perm",
				},
			}
			Expect(claims.IsPlatformOperator()).To(BeTrue())

			claims = &csp.Claims{
				Perms: []string{
					"other:perm",
				},
			}
			Expect(claims.IsPlatformOperator()).To(BeFalse())
		})
	})

	Describe("ParseWithClaims", func() {
		It("parses successfully", func() {
			accessToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2bXdhcmVpZDpiM2I4NDQwNS01Njg4LTQyN2QtOTllNS0zNTk3Njg1NmIzNTYiLCJhenAiOiJjc3BfZGV2X2ludGVybmFsX2NsaWVudF9pZCIsImRvbWFpbiI6InZtd2FyZWlkIiwiY29udGV4dCI6ImE2ODE1MzEyLTQwOTUtNGE5ZC05OGZkLTMwZTFjYjg3MzgxNCIsImlzcyI6Imh0dHBzOi8vZ2F6LWRldi5jc3AtdmlkbS1wcm9kLmNvbSIsInBlcm1zIjpbImNzcDpwbGF0Zm9ybV9vcGVyYXRvciJdLCJjb250ZXh0X25hbWUiOiJ2MjNoZzg3Ni05ODQzLTQxMmYtYTA3MS0zNzVpZzVoNjc0ZiIsImV4cCI6MTUzNjE3MTI2NSwiaWF0IjoxNTM2MTY5NDY1LCJqdGkiOiI1MjEzNzYyYy1mZDU0LTRiMDctYmZiZC1kMTEzMDJlZDA0MGMiLCJ1c2VybmFtZSI6InJtLmludGVncmF0aW9uLnRlc3QudXNlckBnbWFpbC5jb20ifQ.fPvEQ5WwlHhayPuFjMr2iNP8J2kndwp1YQ6EMAaHpUX3AoQ0OI5iXKGkYIDl7Gxm2cG5o3uiIoIiw00geNnd_-R6YoqnP7-S4H8oJjXFATG5ZH9RuvJah8RoNWtUFWKBpr6e_kynexF9bT0urPsNLHG43ALCZ_15EQZXZu0edUr0EojI8NeEBon0oiOm-lTB8Jsqnaj-XiHcHDqLF7_IBKy9ZeaGBqKWbimK-MIZzqxbQ6OmkWQoTGhGy4argGmmd-g1OptxOaQIP7CDShUJtTq4oixgkw4z473Cmhm6uQRvodeHb5JEBtr5m8-AAmgAoTX3tQRKGpO22eDxRdtcCg"

			// this is the decoded token from above
			// {
			// 	"sub": "vmwareid:b3b84405-5688-427d-99e5-35976856b356",
			// 	"azp": "csp_dev_internal_client_id",
			// 	"domain": "vmwareid",
			// 	"context": "a6815312-4095-4a9d-98fd-30e1cb873814",
			// 	"iss": "https://gaz-dev.csp-vidm-prod.com",
			// 	"perms": [
			// 		"csp:platform_operator"
			// 	],
			// 	"context_name": "v23hg876-9843-412f-a071-375ig5h674f",
			// 	"exp": 1536171265,
			// 	"iat": 1536169465,
			// 	"jti": "5213762c-fd54-4b07-bfbd-d11302ed040c",
			// 	"username": "rm.integration.test.user@gmail.com"
			// 	}
			claims := &csp.Claims{}
			_, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
				// token was just retrieved, no need to validate
				return "not a valid key anyway", nil
			})

			// I know the token has expired ok
			Expect(err).To(HaveOccurred())

			Expect(claims.ContextName).To(Equal("v23hg876-9843-412f-a071-375ig5h674f"))
			Expect(claims.Context).To(Equal("a6815312-4095-4a9d-98fd-30e1cb873814"))
			Expect(claims.Username).To(Equal("rm.integration.test.user@gmail.com"))
			Expect(claims.Perms).To(ContainElements(
				"csp:platform_operator",
			))
		})
	})
})
