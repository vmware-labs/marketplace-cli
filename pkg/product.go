// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
)

type ListProductResponse struct {
	Response *ListProductResponsePayload `json:"response"`
}
type ListProductResponsePayload struct {
	Message    string            `json:"string"`
	StatusCode int               `json:"statuscode"`
	Products   []*models.Product `json:"dataList"`
	Params     struct {
		ProductCount int                  `json:"itemsnumber"`
		Pagination   *internal.Pagination `json:"pagination"`
	} `json:"params"`
}

func (m *Marketplace) ListProducts(allOrgs bool, searchTerm string) ([]*models.Product, error) {
	values := url.Values{
		"managed": []string{strconv.FormatBool(!allOrgs)},
	}
	if searchTerm != "" {
		values.Set("search", searchTerm)
	}

	var products []*models.Product
	firstTime := true
	totalProducts := 0
	pagination := &internal.Pagination{
		Page:     1,
		PageSize: 20,
	}

	var progressBar *progressbar.ProgressBar
	for ; firstTime || len(products) < totalProducts; pagination.Page++ {
		requestURL := m.MakeURL("/api/v1/products", values)
		requestURL = pagination.Apply(requestURL)
		resp, err := m.Get(requestURL)
		if err != nil {
			return nil, fmt.Errorf("sending the request for the list of products failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("getting the list of products failed: (%d) %s", resp.StatusCode, resp.Status)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read the list of products: %w", err)
		}

		response := &ListProductResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the list of products: %w", err)
		}
		products = append(products, response.Response.Products...)

		if firstTime {
			totalProducts = response.Response.Params.ProductCount
			progressBar = m.makeRequestProgressBar(totalProducts)
			firstTime = false
		}
		_ = progressBar.Add(len(response.Response.Products))
	}

	return products, nil
}

func (m *Marketplace) makeRequestProgressBar(max int) *progressbar.ProgressBar {
	progressBar := progressbar.NewOptions(
		max,
		progressbar.OptionSetDescription("Getting the list of products"),
		progressbar.OptionSetWriter(m.Output),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("products"),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	_ = progressBar.RenderBlank()
	return progressBar
}

type GetProductResponse struct {
	Response *GetProductResponsePayload `json:"response"`
}
type GetProductResponsePayload struct {
	Message    string          `json:"message"`
	StatusCode int             `json:"statuscode"`
	Data       *models.Product `json:"data"`
}

func (m *Marketplace) GetProduct(slug string) (*models.Product, error) {
	isSlug := true
	_, err := uuid.Parse(slug)
	if err == nil {
		isSlug = false
	}

	requestURL := m.MakeURL(
		fmt.Sprintf("/api/v1/products/%s", slug),
		url.Values{
			"increaseViewCount": []string{"false"},
			"isSlug":            []string{strconv.FormatBool(isSlug)},
		},
	)

	resp, err := m.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("sending the request for product \"%s\" failed: %w", slug, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product \"%s\" not found", slug)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting product \"%s\" failed: (%d)", slug, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response for product \"%s\": %w", slug, err)
	}

	response := &GetProductResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response for product \"%s\": %w", slug, err)
	}
	return response.Response.Data, nil
}

func (m *Marketplace) GetProductWithVersion(slug, version string) (*models.Product, *models.Version, error) {
	product, err := m.GetProduct(slug)
	if err != nil {
		return nil, nil, err
	}

	if !product.HasVersion(version) {
		return nil, nil, fmt.Errorf("product \"%s\" does not have a version %s", slug, version)
	}

	return product, product.GetVersion(version), nil
}

func (m *Marketplace) PutProduct(product *models.Product, versionUpdate bool) (*models.Product, error) {
	encoded, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	requestURL := m.MakeURL(
		fmt.Sprintf("/api/v1/products/%s", product.ProductId),
		url.Values{
			"archivepreviousversion": []string{"false"},
			"isversionupdate":        []string{strconv.FormatBool(versionUpdate)},
		},
	)

	resp, err := m.Put(requestURL, bytes.NewReader(encoded), "application/json")
	if err != nil {
		return nil, fmt.Errorf("sending the update for product \"%s\" failed: %w", product.Slug, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the update response for product \"%s\": %w", product.Slug, err)
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("you do not have permission to modify the product \"%s\"", product.Slug)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("updating product \"%s\" failed: (%d)\n%s", product.Slug, resp.StatusCode, body)
	}

	response := &GetProductResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response for product \"%s\": %w", product.Slug, err)
	}
	return response.Response.Data, nil
}
