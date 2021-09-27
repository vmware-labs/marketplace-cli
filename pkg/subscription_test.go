// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
)

var _ = Describe("Subscription", func() {
	var (
		httpClient  *pkgfakes.FakeHTTPClient
		marketplace *pkg.Marketplace
	)

	BeforeEach(func() {
		httpClient = &pkgfakes.FakeHTTPClient{}
		marketplace = &pkg.Marketplace{
			Client: httpClient,
		}
	})

	Describe("ListSubscriptions", func() {
		BeforeEach(func() {
			response := &pkg.ListSubscriptionsResponse{
				Response: &pkg.ListSubscriptionsResponsePayload{
					Subscriptions: []*models.Subscription{
						{
							ID:          rand.Int(),
							ProductID:   uuid.New().String(),
							ProductName: "My Super Product",
							StatusText:  "awesome",
						},
						{
							ID:          rand.Int(),
							ProductID:   uuid.New().String(),
							ProductName: "My Basic Product",
							StatusText:  "online",
						},
					},
					StatusCode: http.StatusOK,
					Message:    "testing",
				},
			}
			responseBytes, err := json.Marshal(response)
			Expect(err).ToNot(HaveOccurred())

			httpClient.DoReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
			}, nil)
		})

		It("gets the list of subscriptions", func() {
			subscriptions, err := marketplace.ListSubscriptions()
			Expect(err).ToNot(HaveOccurred())

			By("sending the right request", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
				request := httpClient.DoArgsForCall(0)
				Expect(request.Method).To(Equal("GET"))
				Expect(request.URL.Path).To(Equal("/api/v1/subscriptions"))
				Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":1,\"pageSize\":20}"))
			})

			Expect(subscriptions).To(HaveLen(2))
			Expect(subscriptions[0].ProductName).To(Equal("My Super Product"))
			Expect(subscriptions[1].ProductName).To(Equal("My Basic Product"))
		})

		Context("Multiple pages of results", func() {
			BeforeEach(func() {
				var subcriptions []*models.Subscription
				for i := 0; i < 30; i++ {
					subcriptions = append(subcriptions, &models.Subscription{
						ID:          rand.Int(),
						ProductID:   uuid.New().String(),
						ProductName: fmt.Sprintf("My Super Product %d", i),
						StatusText:  "testable",
					})
				}

				response1 := &pkg.ListSubscriptionsResponse{
					Response: &pkg.ListSubscriptionsResponsePayload{
						Subscriptions: subcriptions[:20],
						StatusCode:    http.StatusOK,
						Params: struct {
							SubscriptionCount int                  `json:"itemsnumber"`
							Pagination        *internal.Pagination `json:"pagination"`
						}{
							SubscriptionCount: len(subcriptions),
							Pagination: &internal.Pagination{
								Enabled:  true,
								Page:     1,
								PageSize: 20,
							},
						},
						Message: "testing",
					},
				}
				response2 := &pkg.ListSubscriptionsResponse{
					Response: &pkg.ListSubscriptionsResponsePayload{
						Subscriptions: subcriptions[20:],
						StatusCode:    http.StatusOK,
						Params: struct {
							SubscriptionCount int                  `json:"itemsnumber"`
							Pagination        *internal.Pagination `json:"pagination"`
						}{
							SubscriptionCount: len(subcriptions),
							Pagination: &internal.Pagination{
								Enabled:  true,
								Page:     1,
								PageSize: 20,
							},
						},
						Message: "testing",
					},
				}
				responseBytes, err := json.Marshal(response1)
				Expect(err).ToNot(HaveOccurred())

				httpClient.DoReturnsOnCall(0, &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
				}, nil)

				responseBytes, err = json.Marshal(response2)
				Expect(err).ToNot(HaveOccurred())

				httpClient.DoReturnsOnCall(1, &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(responseBytes)),
				}, nil)
			})

			It("returns all results", func() {
				subscriptions, err := marketplace.ListSubscriptions()
				Expect(err).ToNot(HaveOccurred())

				By("sending the correct requests", func() {
					Expect(httpClient.DoCallCount()).To(Equal(2))
					request := httpClient.DoArgsForCall(0)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.Path).To(Equal("/api/v1/subscriptions"))
					Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":1,\"pageSize\":20}"))

					request = httpClient.DoArgsForCall(1)
					Expect(request.Method).To(Equal("GET"))
					Expect(request.URL.Path).To(Equal("/api/v1/subscriptions"))
					Expect(request.URL.Query().Get("pagination")).To(Equal("{\"page\":2,\"pageSize\":20}"))
				})

				By("returning all subscriptions", func() {
					Expect(subscriptions).To(HaveLen(30))
				})
			})
		})

		Context("Error fetching subscriptions", func() {
			BeforeEach(func() {
				httpClient.DoReturns(nil, fmt.Errorf("request failed"))
			})

			It("prints the error", func() {
				_, err := marketplace.ListSubscriptions()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("sending the request for the list of subscriptions failed: marketplace request failed: request failed"))
			})
		})

		Context("Unexpected status code", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusTeapot,
					Status:     http.StatusText(http.StatusTeapot),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.ListSubscriptions()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("getting the list of subscriptions failed: (418) I'm a teapot"))
			})
		})

		Context("Un-parsable response", func() {
			BeforeEach(func() {
				httpClient.DoReturns(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader("This totally isn't a valid response")),
				}, nil)
			})

			It("prints the error", func() {
				_, err := marketplace.ListSubscriptions()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to parse the list of subscriptions: invalid character 'T' looking for beginning of value"))
			})
		})
	})
})
