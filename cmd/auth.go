// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	csp2 "github.com/vmware-labs/marketplace-cli/v2/internal/csp"
)

//go:generate counterfeiter . TokenServices
type TokenServices interface {
	Redeem(refreshToken string) (*csp2.Claims, error)
}

//go:generate counterfeiter . TokenServicesInitializer
type TokenServicesInitializer func(cspHost string) (TokenServices, error)

var InitializeTokenServices TokenServicesInitializer = func(cspHost string) (TokenServices, error) {
	return csp2.NewTokenServices(cspHost)
}

func GetRefreshToken(cmd *cobra.Command, args []string) error {
	tokenServices, err := InitializeTokenServices(
		fmt.Sprintf("https://%s/", viper.GetString("csp.host")),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize token services: %w", err)
	}

	apiToken := viper.GetString("csp.api-token")
	if apiToken == "" {
		return fmt.Errorf("missing CSP API token")
	}

	claims, err := tokenServices.Redeem(apiToken)
	if err != nil {
		return fmt.Errorf("failed to exchange api token: %w", err)
	}

	viper.Set("csp.refresh-token", claims.Token)
	return nil
}

type CredentialsResponse struct {
	AccessID     string    `json:"accessId"`
	AccessKey    string    `json:"accessKey"`
	SessionToken string    `json:"sessionToken"`
	Expiration   time.Time `json:"expiration"`
}

func GetUploadCredentials(cmd *cobra.Command, args []string) error {
	requestURL := Marketplace.MakeURL("/aws/credentials/generate", url.Values{})
	requestURL.Host = Marketplace.APIHost
	response, err := Marketplace.Get(requestURL)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch credentials: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	credsResponse := &CredentialsResponse{}
	err = json.Unmarshal(body, credsResponse)
	if err != nil {
		return err
	}

	UploadCredentials.AccessKeyID = credsResponse.AccessID
	UploadCredentials.SecretAccessKey = credsResponse.AccessKey
	UploadCredentials.SessionToken = credsResponse.SessionToken
	UploadCredentials.Expires = credsResponse.Expiration

	return nil
}

// This adds this as a secret command, please remove from finished product.
func init() {
	rootCmd.AddCommand(AuthCmd)
}

var AuthCmd = &cobra.Command{
	Use:     "auth",
	Short:   "fetch and return a valid CSP refresh token",
	Hidden:  true,
	PreRunE: GetRefreshToken,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(viper.GetString("csp.refresh-token"))
	},
}
