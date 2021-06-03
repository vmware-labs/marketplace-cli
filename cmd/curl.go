package cmd

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	. "gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/lib"
)

func init() {
	rootCmd.AddCommand(CurlCmd)
}

var CurlCmd = &cobra.Command{
	Use:     "curl",
	Hidden:  true,
	PreRunE: GetRefreshToken,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		req, err := MakeGetRequest(args[0], url.Values{})
		if err != nil {
			return err
		}

		resp, err := Client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return errors.Errorf("request failed (%d)", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response")
		}

		cmd.Println(string(body))
		return nil
	},
}
