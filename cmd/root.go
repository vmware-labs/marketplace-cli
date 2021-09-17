// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

func RunSerially(funcs ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range funcs {
			err := fn(cmd, args)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func EnableDebugging(command *cobra.Command, _ []string) error {
	if debug {
		Marketplace.Client = pkg.EnableDebugging(debugRequestPayloads, Marketplace.Client, command.ErrOrStderr())
	}
	return nil
}

func ValidateOutputFormatFlag(command *cobra.Command, _ []string) error {
	if OutputFormat == output.FormatHuman {
		Output = output.NewHumanOutput(command.OutOrStdout(), Marketplace.UIHost)
	} else if OutputFormat == output.FormatJSON {
		Output = output.NewJSONOutput(command.OutOrStdout())
	} else if OutputFormat == output.FormatYAML {
		Output = output.NewYAMLOutput(command.OutOrStdout())
	} else {
		return fmt.Errorf("output format not supported: %s", OutputFormat)
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: fmt.Sprintf("%s is a CLI interface for the VMware Marketplace", AppName),
	Long: fmt.Sprintf(`%s is a CLI interface for the VMware Marketplace,
enabling users to view, get, and manage their Marketplace entries.`, AppName),
	PersistentPreRunE: RunSerially(EnableDebugging, ValidateOutputFormatFlag, GetRefreshToken),
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	_ = rootCmd.PersistentFlags().MarkHidden("debug")

	rootCmd.PersistentFlags().BoolVar(&debugRequestPayloads, "debug-request-payloads", false, "Also print request payloads")
	_ = rootCmd.PersistentFlags().MarkHidden("debug-request-payloads")

	_ = viper.BindEnv("csp.api-token", "CSP_API_TOKEN")
	rootCmd.PersistentFlags().String(
		"csp-api-token",
		"",
		"VMware Cloud Service Platform API Token, used for authentication to the VMware Marketplace",
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

	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", output.FormatHuman, fmt.Sprintf("Output format. One of %s", strings.Join(output.SupportedOutputs, "|")))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
