// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

type ListSubscriptionsResponse struct {
	Response *ListSubscriptionsResponsePayload `json:"response"`
}
type ListSubscriptionsResponsePayload struct {
	Message       string                 `json:"string"`
	StatusCode    int                    `json:"statuscode"`
	Subscriptions []*models.Subscription `json:"subscriptionsList"`
	Params        struct {
		SubscriptionCount int                  `json:"itemsnumber"`
		Pagination        *internal.Pagination `json:"pagination"`
	} `json:"params"`
}

func (m *Marketplace) ListSubscriptions() ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	totalSubscriptions := 1
	pagination := &internal.Pagination{
		Page:     1,
		PageSize: 20,
	}

	for ; len(subscriptions) < totalSubscriptions; pagination.Page++ {
		requestURL := m.MakeURL("/api/v1/subscriptions", nil)
		requestURL = pagination.Apply(requestURL)
		resp, err := m.Get(requestURL)
		if err != nil {
			return nil, fmt.Errorf("sending the request for the list of subscriptions failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("getting the list of subscriptions failed: (%d) %s", resp.StatusCode, resp.Status)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read the list of subscriptions: %w", err)
		}

		response := &ListSubscriptionsResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the list of subscriptions: %w", err)
		}
		totalSubscriptions = response.Response.Params.SubscriptionCount
		subscriptions = append(subscriptions, response.Response.Subscriptions...)
	}

	return subscriptions, nil
}
