package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/lib/csp"
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
		fmt.Sprintf("https://%s/", viper.GetString("csp.host")),
	)
	if err != nil {
		return errors.Wrap(err, "failed to initialize token services")
	}

	apiToken := viper.GetString("csp.api-token")
	if apiToken == "" {
		return errors.New("missing CSP API token")
	}

	claims, err := tokenServices.Redeem(apiToken)
	if err != nil {
		return errors.Wrap(err, "failed to exchange api token")
	}

	viper.Set("csp.refresh-token", claims.Token)
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
