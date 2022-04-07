// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(ConfigCmd)
	ConfigCmd.SetOut(ConfigCmd.OutOrStdout())
}

var ConfigCmd = &cobra.Command{
	Use:    "config",
	Short:  "prints the current config",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := viper.AllSettings()
		formattedConfig, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return err
		}
		cmd.Println(string(formattedConfig))
		return nil
	},
}
