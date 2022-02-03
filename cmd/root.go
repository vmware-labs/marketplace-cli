// Copyright 2022 VMware, Inc.
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
	if viper.GetBool("debugging.enabled") {
		Marketplace.EnableDebugging(viper.GetBool("debugging.print-request-payloads"), command.ErrOrStderr())
	}
	return nil
}

func ValidateOutputFormatFlag(command *cobra.Command, _ []string) error {
	outputFormat := viper.GetString("output_format")
	if outputFormat == output.FormatHuman {
		Output = output.NewHumanOutput(command.OutOrStdout(), Marketplace.GetUIHost())
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
	viper.SetDefault("debugging.enabled", false)
	_ = viper.BindEnv("debugging.enabled", "MKPCLI_DEBUG")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug output [$MKPCLI_DEBUG}")
	_ = rootCmd.PersistentFlags().MarkHidden("debug")
	_ = viper.BindPFlag("debugging.enabled", rootCmd.PersistentFlags().Lookup("debug"))

	viper.SetDefault("debugging.print-request-payloads", false)
	_ = viper.BindEnv("debugging.print-request-payloads", "MKPCLI_DEBUG_REQUEST_PAYLOADS")
	rootCmd.PersistentFlags().Bool("debug-request-payloads", false, "Also print request payloads [$MKPCLI_DEBUG_REQUEST_PAYLOADS]")
	_ = rootCmd.PersistentFlags().MarkHidden("debug-request-payloads")
	_ = viper.BindPFlag("debugging.print-request-payloads", rootCmd.PersistentFlags().Lookup("debug-request-payloads"))

	viper.SetDefault("csp.api-token", "")
	_ = viper.BindEnv("csp.api-token", "CSP_API_TOKEN")
	rootCmd.PersistentFlags().String("csp-api-token", "", "VMware Cloud Service Platform API Token, used for authenticating to the VMware Marketplace [$CSP_API_TOKEN]")
	_ = viper.BindPFlag("csp.api-token", rootCmd.PersistentFlags().Lookup("csp-api-token"))

	viper.SetDefault("csp.host", "console.cloud.vmware.com")
	_ = viper.BindEnv("csp.host", "CSP_HOST")
	rootCmd.PersistentFlags().String("csp-host", "console.cloud.vmware.com", "Host for VMware Cloud Service Platform")
	_ = rootCmd.PersistentFlags().MarkHidden("csp-host")
	_ = viper.BindPFlag("csp.host", rootCmd.PersistentFlags().Lookup("csp-host"))

	if os.Getenv("MARKETPLACE_ENV") == "staging" {
		Marketplace = &pkg.Marketplace{
			Host:          "gtwstg.market.csp.vmware.com",
			APIHost:       "apistg.market.csp.vmware.com",
			UIHost:        "stg.market.csp.vmware.com",
			StorageBucket: "cspmarketplacestage",
			StorageRegion: "us-east-2",
			Client:        pkg.NewClient(),
			Output:        os.Stderr,
		}
	} else {
		Marketplace = &pkg.Marketplace{
			Host:          "gtw.marketplace.cloud.vmware.com",
			APIHost:       "api.marketplace.cloud.vmware.com",
			UIHost:        "marketplace.cloud.vmware.com",
			StorageBucket: "cspmarketplaceprd",
			StorageRegion: "us-west-2",
			Client:        pkg.NewClient(),
			Output:        os.Stderr,
		}
	}

	viper.SetDefault("marketplace.strict-decoding", false)
	_ = viper.BindEnv("marketplace.strict-decoding", "MKPCLI_STRICT_DECODING")
	if viper.GetBool("marketplace.strict-decoding") {
		Marketplace.EnableStrictDecoding()
	}

	viper.SetDefault("output_format", output.FormatHuman)
	_ = viper.BindEnv("output_format", "MKPCLI_OUTPUT")
	rootCmd.PersistentFlags().StringP("output", "o", output.FormatHuman, fmt.Sprintf("Output format. One of %s. [$MKPCLI_OUTPUT]", strings.Join(output.SupportedOutputs, "|")))
	_ = viper.BindPFlag("output_format", rootCmd.PersistentFlags().Lookup("output"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
