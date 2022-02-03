// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/spf13/cobra"
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
	Hidden:  true,
	PreRunE: GetRefreshToken,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		inputURL, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		requestURL := Marketplace.MakeURL(inputURL.Path, inputURL.Query())
		if useAPIHost {
			requestURL.Host = Marketplace.GetAPIHost()
		}

		cmd.PrintErrf("Sending %s request to %s...\n", method, requestURL.String())

		var content io.Reader
		headers := map[string]string{}
		if payload != "" {
			payloadBytes, err := ioutil.ReadFile(payload)
			if err != nil {
				return fmt.Errorf("failed to read payload file: %w", err)
			}
			content = bytes.NewReader(payloadBytes)
			headers["Content-Type"] = "application/json"
		}

		resp, err := Marketplace.SendRequest(method, requestURL, headers, content)
		if err != nil {
			return err
		}

		cmd.PrintErrf("Response status %d\n", resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		cmd.PrintErrln("Body:")
		cmd.Println(string(body))
		return nil
	},
}
