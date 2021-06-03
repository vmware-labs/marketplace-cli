package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	. "gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/lib"
)

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: fmt.Sprintf("%s is a CLI interface for the VMware Tanzu Marketplace", AppName),
	Long: fmt.Sprintf(`%s is a CLI interface for the VMware Tanzu Marketplace,
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

	_ = viper.BindEnv("marketplace.host", "MARKETPLACE_HOST")
	rootCmd.PersistentFlags().String(
		"marketplace-host",
		"",
		"Host for the Marketplace API",
	)
	_ = rootCmd.PersistentFlags().MarkHidden("marketplace-host")
	_ = viper.BindPFlag("marketplace.host", rootCmd.PersistentFlags().Lookup("marketplace-host"))
	if viper.GetString("marketplace.host") == "" {
		viper.Set("marketplace.host", "gtw.marketplace.cloud.vmware.com")
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
