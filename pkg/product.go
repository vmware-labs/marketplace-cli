// Copyright 2022 VMware, Inc.
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

type VersionDoesNotExistError struct {
	Product string
	Version string
}

func (e *VersionDoesNotExistError) Error() string {
	return fmt.Sprintf("product \"%s\" does not have version %s", e.Product, e.Version)
}

func (e *VersionDoesNotExistError) Is(otherError error) bool {
	_, ok := otherError.(*VersionDoesNotExistError)
	return ok
}

type ListProductResponse struct {
	Response *ListProductResponsePayload `json:"response"`
}

type ListProductResponseParams struct {
	Filters      map[string][]interface{} `json:"filters"`
	Pagination   *internal.Pagination     `json:"pagination"`
	ProductCount int                      `json:"itemsnumber"`
	Search       string                   `json:"search"`
	SortingList  []*internal.Sorting      `json:"sortinglist"`
	SelectList   []interface{}            `json:"selectlist"`
}

type ListProductResponsePayload struct {
	Message          string                     `json:"message"`
	StatusCode       int                        `json:"statuscode"`
	Products         []*models.Product          `json:"dataList"`
	AvailableFilters interface{}                `json:"availablefilters"`
	Params           *ListProductResponseParams `json:"params"`
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
	sorting := &internal.Sorting{
		Order:     1,
		Key:       internal.SortKeyDisplayName,
		Direction: internal.SortDirectionAscending,
	}

	var progressBar *progressbar.ProgressBar
	for ; firstTime || len(products) < totalProducts; pagination.Page++ {
		requestURL := MakeURL(m.GetHost(), "/api/v1/products", values)
		ApplyParameters(requestURL, pagination, sorting)
		resp, err := m.Client.Get(requestURL)
		if err != nil {
			return nil, fmt.Errorf("sending the request for the list of products failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				body = []byte{}
			}
			return nil, fmt.Errorf("getting the list of products failed: (%d) %s: %s", resp.StatusCode, resp.Status, body)
		}

		response := &ListProductResponse{}
		err = m.DecodeJson(resp.Body, response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the list of products: %w", err)
		}

		// Return immediately if we get an empty list.
		// On empty lists, we cannot necessarily trust response.Response.Params.ProductCount
		// See: https://github.com/vmware-labs/marketplace-cli/issues/62
		if len(response.Response.Products) == 0 {
			return products, nil
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

type VersionSpecificDetailsRequestPayload struct {
	ProductId     string `json:"productId"`
	VersionNumber string `json:"versionNumber"`
}

type VersionSpecificDetailsPayloadResponse struct {
	Response *VersionSpecificDetailsPayload `json:"response"`
}

type VersionSpecificDetailsPayload struct {
	Message    string                                `json:"message"`
	StatusCode int                                   `json:"statuscode"`
	Data       *models.VersionSpecificProductDetails `json:"data"`
}

func (m *Marketplace) GetProduct(slug string) (*models.Product, error) {
	isSlug := true
	_, err := uuid.Parse(slug)
	if err == nil {
		isSlug = false
	}

	requestURL := MakeURL(
		m.GetHost(),
		fmt.Sprintf("/api/v1/products/%s", slug),
		url.Values{
			"increaseViewCount": []string{"false"},
			"isSlug":            []string{strconv.FormatBool(isSlug)},
		},
	)

	resp, err := m.Client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("sending the request for product %s failed: %w", slug, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product %s not found", slug)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return nil, fmt.Errorf("getting product %s failed: (%d)\n%s", slug, resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("getting product %s failed: (%d)", slug, resp.StatusCode)
	}

	response := &GetProductResponse{}
	err = m.DecodeJson(resp.Body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response for product %s: %w", slug, err)
	}

	product := response.Response.Data
	if product.HasVersion("") {
		product.LatestVersion = product.GetLatestVersion().Number
	}
	return product, nil
}

func (m *Marketplace) getVersionDetails(product *models.Product, version string) (*models.VersionSpecificProductDetails, error) {
	requestURL := MakeURL(
		m.GetHost(),
		fmt.Sprintf("/api/v1/products/%s/version-details", product.ProductId),
		url.Values{
			"versionNumber": []string{version},
		},
	)

	payload := &VersionSpecificDetailsRequestPayload{
		ProductId:     product.ProductId,
		VersionNumber: version,
	}

	resp, err := m.Client.PostJSON(requestURL, payload)
	if err != nil {
		return nil, fmt.Errorf("sending the product version details request for %s %s failed: %w", product.Slug, version, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product version details for %s %s not found", product.Slug, version)
	}

	// Workaround: Always ignore 400 errors, because they often indicate that there is no version-specific details
	if resp.StatusCode == http.StatusBadRequest {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("getting product version details for %s %s failed: (%d)", product.Slug, version, resp.StatusCode)
		}

		return nil, fmt.Errorf("getting product version details for %s %s failed: (%d)\n%s", product.Slug, version, resp.StatusCode, string(body))
	}

	response := &VersionSpecificDetailsPayloadResponse{}
	err = m.DecodeJson(resp.Body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response for product %s %s: %w", product.Slug, version, err)
	}
	return response.Response.Data, nil
}

func (m *Marketplace) GetProductWithVersion(slug, version string) (*models.Product, *models.Version, error) {
	product, err := m.GetProduct(slug)
	if err != nil {
		return nil, nil, err
	}

	if !product.HasVersion(version) {
		return product, nil, &VersionDoesNotExistError{Product: slug, Version: version}
	}
	versionObject := product.GetVersion(version)

	versionDetails, err := m.getVersionDetails(product, versionObject.Number)
	if err != nil {
		return nil, nil, err
	}

	if versionDetails != nil {
		product.UpdateWithVersionSpecificDetails(versionObject.Number, versionDetails)
	}

	product.CurrentVersion = versionObject.Number
	return product, versionObject, nil
}

func (m *Marketplace) PutProduct(product *models.Product, versionUpdate bool) (*models.Product, error) {
	encoded, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	requestURL := MakeURL(
		m.GetHost(),
		fmt.Sprintf("/api/v1/products/%s", product.ProductId),
		url.Values{
			"archivepreviousversion": []string{"false"},
			"isversionupdate":        []string{strconv.FormatBool(versionUpdate)},
		},
	)

	resp, err := m.Client.Put(requestURL, bytes.NewReader(encoded), "application/json")
	if err != nil {
		return nil, fmt.Errorf("sending the update for product \"%s\" failed: %w", product.Slug, err)
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("you do not have permission to modify the product \"%s\"", product.Slug)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			body = []byte{}
		}
		return nil, fmt.Errorf("updating product \"%s\" failed: (%d)\n%s", product.Slug, resp.StatusCode, body)
	}

	response := &GetProductResponse{}
	err = m.DecodeJson(resp.Body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the response for product \"%s\": %w", product.Slug, err)
	}
	return response.Response.Data, nil
}
