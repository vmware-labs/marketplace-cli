// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var Marketplace *pkg.Marketplace

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: fmt.Sprintf("%s is a CLI interface for the VMware Marketplace", AppName),
	Long: fmt.Sprintf(`%s is a CLI interface for the VMware Marketplace,
enabling users to view, get, and manage their Marketplace entries.`, AppName),
}

func init() {
	_ = viper.BindEnv("csp.api-token", "CSP_API_TOKEN")
	rootCmd.PersistentFlags().String(
		"csp-api-token",
		"",
		"VMware Cloud Service Platform API token, used for authentication to Marketplace",
	)
	_ = viper.BindPFlag("csp.api-token", rootCmd.PersistentFlags().Lookup("csp-api-token"))

	rootCmd.PersistentFlags().String(
		"csp-host",
		"console.cloud.vmware.com",
		"Host for CSP",
	)
	_ = rootCmd.PersistentFlags().MarkHidden("csp-host")
	_ = viper.BindPFlag("csp.host", rootCmd.PersistentFlags().Lookup("csp-host"))
	viper.SetDefault("csp.host", "console.cloud.vmware.com")

	Marketplace = pkg.ProductionConfig
	if os.Getenv("MARKETPLACE_ENV") == "staging" {
		Marketplace = pkg.StagingConfig
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
