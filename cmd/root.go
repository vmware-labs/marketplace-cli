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
enabling users to view, get, and manage their Marketplace products.`, AppName),
	PersistentPreRunE: RunSerially(
		func(cmd *cobra.Command, args []string) error {
			Client = pkg.NewClient(
				os.Stderr,
				viper.GetBool("debugging.enabled"),
				viper.GetBool("debugging.print-request-payloads"),
				viper.GetBool("debugging.print-response-payloads"),
			)

			Marketplace = &pkg.Marketplace{
				Host:          viper.GetString("marketplace.host"),
				APIHost:       viper.GetString("marketplace.api-host"),
				UIHost:        viper.GetString("marketplace.ui-host"),
				StorageBucket: viper.GetString("marketplace.storage.bucket"),
				StorageRegion: viper.GetString("marketplace.storage.region"),
				Client:        Client,
				Output:        os.Stderr,
			}

			if viper.GetBool("marketplace.strict-decoding") {
				Marketplace.EnableStrictDecoding()
			}
			return nil
		},
		ValidateOutputFormatFlag,
	),
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

	viper.SetDefault("debugging.print-response-payloads", false)

	viper.SetDefault("csp.api-token", "")
	_ = viper.BindEnv("csp.api-token", "CSP_API_TOKEN")
	rootCmd.PersistentFlags().String("csp-api-token", "", "VMware Cloud Service Platform API Token, used for authenticating to the VMware Marketplace [$CSP_API_TOKEN]")
	_ = viper.BindPFlag("csp.api-token", rootCmd.PersistentFlags().Lookup("csp-api-token"))

	viper.SetDefault("csp.host", "console.cloud.vmware.com")
	_ = viper.BindEnv("csp.host", "CSP_HOST")
	rootCmd.PersistentFlags().String("csp-host", "console.cloud.vmware.com", "Host for VMware Cloud Service Platform")
	_ = rootCmd.PersistentFlags().MarkHidden("csp-host")
	_ = viper.BindPFlag("csp.host", rootCmd.PersistentFlags().Lookup("csp-host"))

	_ = viper.BindEnv("marketplace.host", "MKPCLI_HOST")
	_ = viper.BindEnv("marketplace.api-host", "MKPCLI_API_HOST")
	_ = viper.BindEnv("marketplace.ui-host", "MKPCLI_UI_HOST")
	_ = viper.BindEnv("marketplace.storage.bucket", "MKPCLI_STORAGE_BUCKET")
	_ = viper.BindEnv("marketplace.storage.region", "MKPCLI_STORAGE_REGION")

	if os.Getenv("MARKETPLACE_ENV") == "staging" {
		viper.SetDefault("marketplace.host", "gtwstg.market.csp.vmware.com")
		viper.SetDefault("marketplace.api-host", "apistg.market.csp.vmware.com")
		viper.SetDefault("marketplace.ui-host", "stg.market.csp.vmware.com")
		viper.SetDefault("marketplace.storage.bucket", "cspmarketplacestage")
		viper.SetDefault("marketplace.storage.region", "us-east-2")
	} else {
		viper.SetDefault("marketplace.host", "gtw.marketplace.cloud.vmware.com")
		viper.SetDefault("marketplace.api-host", "api.marketplace.cloud.vmware.com")
		viper.SetDefault("marketplace.ui-host", "marketplace.cloud.vmware.com")
		viper.SetDefault("marketplace.storage.bucket", "cspmarketplaceprd")
		viper.SetDefault("marketplace.storage.region", "us-west-2")
	}

	viper.SetDefault("marketplace.strict-decoding", false)
	_ = viper.BindEnv("marketplace.strict-decoding", "MKPCLI_STRICT_DECODING")

	viper.SetDefault("output_format", output.FormatHuman)
	_ = viper.BindEnv("output_format", "MKPCLI_OUTPUT")
	rootCmd.PersistentFlags().StringP("output", "o", output.FormatHuman, fmt.Sprintf("Output format. One of %s. [$MKPCLI_OUTPUT]", strings.Join(output.SupportedOutputs, "|")))
	_ = viper.BindPFlag("output_format", rootCmd.PersistentFlags().Lookup("output"))

	viper.SetDefault("skip_ssl_validation", "false")
	_ = viper.BindEnv("skip_ssl_validation", "MKPCLI_SKIP_SSL_VALIDATION")
	rootCmd.PersistentFlags().Bool("skip-ssl-validation", false, "Skip SSL certificate validation during HTTP requests")
	_ = rootCmd.PersistentFlags().MarkHidden("skip-ssl-validation")
	_ = viper.BindPFlag("skip_ssl_validation", rootCmd.PersistentFlags().Lookup("skip-ssl-validation"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
