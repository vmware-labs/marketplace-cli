// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/internal/csp"
)

//go:generate counterfeiter . TokenServices
type TokenServices interface {
	Redeem(refreshToken string) (*csp.Claims, error)
}

//go:generate counterfeiter . TokenServicesInitializer
type TokenServicesInitializer func(cspHost string) (TokenServices, error)

var InitializeTokenServices TokenServicesInitializer = func(cspHost string) (TokenServices, error) {
	return csp.NewTokenServices(cspHost)
}

func GetRefreshToken(cmd *cobra.Command, args []string) error {
	tokenServices, err := InitializeTokenServices(
		fmt.Sprintf("https://%s", viper.GetString("csp.host")),
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
