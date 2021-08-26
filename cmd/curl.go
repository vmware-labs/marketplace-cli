// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(curlCmd)
	curlCmd.SetOut(curlCmd.OutOrStdout())
}

var curlCmd = &cobra.Command{
	Use:     "curl",
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

		cmd.PrintErrf("Sending request to %s...\n", requestURL.String())
		resp, err := Marketplace.Get(requestURL)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("request failed (%d)", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		cmd.Println(string(body))
		return nil
	},
}
