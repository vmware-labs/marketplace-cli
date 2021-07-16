// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	. "github.com/vmware-labs/marketplace-cli/v2/lib"
	"github.com/vmware-labs/marketplace-cli/v2/models"
)

func init() {
	rootCmd.AddCommand(ProductCmd)
	ProductCmd.AddCommand(ListProductsCmd)
	ProductCmd.AddCommand(GetProductCmd)
	ProductCmd.PersistentFlags().StringVarP(&OutputFormat, "output-format", "f", FormatTable, "Output format")

	ListProductsCmd.Flags().StringVar(&SearchTerm, "search-text", "", "Filter by text")

	GetProductCmd.Flags().StringVarP(&ProductSlug, "product", "p", "", "Product slug")
	_ = GetProductCmd.MarkFlagRequired("product")
}

var ProductCmd = &cobra.Command{
	Use:               "product",
	Aliases:           []string{"products"},
	Short:             "stuff related to products",
	Long:              "",
	Args:              cobra.OnlyValidArgs,
	ValidArgs:         []string{"get", "list"},
	PersistentPreRunE: GetRefreshToken,
}

type ListProductResponse struct {
	Response *ListProductResponsePayload `json:"response"`
}
type ListProductResponsePayload struct {
	Message    string            `json:"string"`
	StatusCode int               `json:"statuscode"`
	Products   []*models.Product `json:"dataList"`
}

var ListProductsCmd = &cobra.Command{
	Use:  "list",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		values := url.Values{
			"pagination": Pagination(0, 20),
			"ownOrg":     []string{"true"},
		}
		if SearchTerm != "" {
			values.Set("search", SearchTerm)
		}

		req, err := MarketplaceConfig.MakeGetRequest("/api/v1/products", values)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("preparing the request for the list of products failed: %w", err)
		}

		resp, err := Client.Do(req)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("sending the request for the list of products failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			cmd.SilenceUsage = true
			return fmt.Errorf("getting the list of products failed: (%d) %s", resp.StatusCode, resp.Status)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to read the list of products: %w", err)
		}

		response := &ListProductResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to parse the list of products: %w", err)
		}

		err = RenderProductList(OutputFormat, response.Response.Products, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to render the list of products: %w", err)
		}

		return nil
	},
}

type GetProductResponse struct {
	Response *GetProductResponsePayload `json:"response"`
}
type GetProductResponsePayload struct {
	Message    string          `json:"string"`
	StatusCode int             `json:"statuscode"`
	Data       *models.Product `json:"data"`
}

func GetProduct(slug string, response *GetProductResponse) error {
	req, err := MarketplaceConfig.MakeGetRequest(
		fmt.Sprintf("api/v1/products/%s", slug),
		url.Values{
			"increaseViewCount": []string{"false"},
			"isSlug":            []string{"true"},
		},
	)
	if err != nil {
		return err
	}

	resp, err := Client.Do(req)
	if err != nil {
		return fmt.Errorf("sending the request for product \"%s\" failed: %w", ProductSlug, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("product \"%s\" not found", ProductSlug)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("getting product \"%s\" failed: (%d)", ProductSlug, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response for product \"%s\": %w", ProductSlug, err)
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("failed to parse the response for product \"%s\": %w", ProductSlug, err)
	}
	return nil
}

func PutProduct(product *models.Product, versionUpdate bool, response *GetProductResponse) error {
	product.PrepForUpdate()
	encoded, err := json.Marshal(product)
	if err != nil {
		return err
	}

	req, err := MarketplaceConfig.MakeRequest(
		"PUT",
		fmt.Sprintf("/api/v1/products/%s", product.ProductId),
		url.Values{
			"archivepreviousversion": []string{"false"},
			"isversionupdate":        []string{strconv.FormatBool(versionUpdate)},
		},
		map[string]string{
			"Content-Type": "application/json",
		},
		bytes.NewReader(encoded),
	)

	resp, err := Client.Do(req)
	if err != nil {
		return fmt.Errorf("sending the update for product \"%s\" failed: %w", ProductSlug, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read the update response for product \"%s\": %w", ProductSlug, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("updating product \"%s\" failed: (%d)\n%s", ProductSlug, resp.StatusCode, body)
	}

	return json.Unmarshal(body, response)
}

var GetProductCmd = &cobra.Command{
	Use:  "get [product slug]",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		response := &GetProductResponse{}
		err := GetProduct(ProductSlug, response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		err = RenderProduct(OutputFormat, response.Response.Data, cmd.OutOrStdout())
		if err != nil {
			cmd.SilenceUsage = true
			return fmt.Errorf("failed to render the product: %w", err)
		}
		return nil
	},
}
