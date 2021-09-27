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
	if viper.GetBool("debug") {
		Marketplace.Client = pkg.EnableDebugging(viper.GetBool("debug-request-payloads"), Marketplace.Client, command.ErrOrStderr())
	}
	return nil
}

func ValidateOutputFormatFlag(command *cobra.Command, _ []string) error {
	outputFormat := viper.GetString("output_format")
	if outputFormat == output.FormatHuman {
		Output = output.NewHumanOutput(command.OutOrStdout(), Marketplace.UIHost)
	} else if outputFormat == output.FormatJSON {
		Output = output.NewJSONOutput(command.OutOrStdout())
	} else if outputFormat == output.FormatYAML {
		Output = output.NewYAMLOutput(command.OutOrStdout())
	} else {
		return fmt.Errorf("output format not supported: %s", outputFormat)
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
	// debug
	_ = viper.BindEnv("debug", "MKPCLI_DEBUG")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug output [$MKPCLI_DEBUG}")
	_ = rootCmd.PersistentFlags().MarkHidden("debug")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// debug-request-payloads
	_ = viper.BindEnv("debug", "MKPCLI_DEBUG_REQUEST_PAYLOADS")
	rootCmd.PersistentFlags().Bool("debug-request-payloads", false, "Also print request payloads [$MKPCLI_DEBUG_REQUEST_PAYLOADS]")
	_ = rootCmd.PersistentFlags().MarkHidden("debug-request-payloads")
	_ = viper.BindPFlag("debug-request-payloads", rootCmd.PersistentFlags().Lookup("debug-request-payloads"))

	// csp-api-token
	_ = viper.BindEnv("csp.api-token", "CSP_API_TOKEN")
	rootCmd.PersistentFlags().String("csp-api-token", "", "VMware Cloud Service Platform API Token, used for authenticating to the VMware Marketplace [$CSP_API_TOKEN]")
	_ = viper.BindPFlag("csp.api-token", rootCmd.PersistentFlags().Lookup("csp-api-token"))

	_ = viper.BindEnv("csp.host", "CSP_HOST")
	rootCmd.PersistentFlags().String("csp-host", "console.cloud.vmware.com", "Host for VMware Cloud Service Platform")
	_ = rootCmd.PersistentFlags().MarkHidden("csp-host")
	_ = viper.BindPFlag("csp.host", rootCmd.PersistentFlags().Lookup("csp-host"))

	Marketplace = pkg.ProductionConfig
	if os.Getenv("MARKETPLACE_ENV") == "staging" {
		Marketplace = pkg.StagingConfig
	}

	_ = viper.BindEnv("output_format", "MKPCLI_OUTPUT")
	rootCmd.PersistentFlags().StringP("output", "o", output.FormatHuman, fmt.Sprintf("Output format. One of %s. [$MKPCLI_OUTPUT]", strings.Join(output.SupportedOutputs, "|")))
	_ = viper.BindPFlag("output_format", rootCmd.PersistentFlags().Lookup("output"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
