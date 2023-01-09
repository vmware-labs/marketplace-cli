// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	method     = "GET"
	payload    string
	useAPIHost = false
)

func init() {
	rootCmd.AddCommand(curlCmd)
	curlCmd.SetOut(curlCmd.OutOrStdout())
	curlCmd.Flags().StringVarP(&method, "method", "X", method, "HTTP verb to use")
	curlCmd.Flags().StringVar(&payload, "payload", "", "JSON file containing the payload to send as a request body")
	curlCmd.Flags().BoolVar(&useAPIHost, "use-api-host", false, "Send request to the API host, rather than the gateway host")
}

var curlCmd = &cobra.Command{
	Use:     "curl [/api/v1/path]",
	Long:    "Sends an HTTP request to the Marketplace",
	Example: fmt.Sprintf("%s curl /api/v1/products", AppName),
	Hidden:  true,
	PreRunE: GetRefreshToken,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		inputURL, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		host := Marketplace.GetHost()
		if useAPIHost {
			host = Marketplace.GetAPIHost()
		}

		requestURL := pkg.MakeURL(host, inputURL.Path, inputURL.Query())

		cmd.PrintErrf("Sending %s request to %s...\n", method, requestURL.String())

		var content io.Reader
		headers := map[string]string{}
		if payload != "" {
			payloadBytes, err := os.ReadFile(payload)
			if err != nil {
				return fmt.Errorf("failed to read payload file: %w", err)
			}
			content = bytes.NewReader(payloadBytes)
			headers["Content-Type"] = "application/json"
		}

		resp, err := Client.SendRequest(method, requestURL, headers, content)
		if err != nil {
			return err
		}

		cmd.PrintErrf("Response status %d\n", resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		cmd.PrintErrln("Body:")
		cmd.Println(string(body))
		return nil
	},
}
