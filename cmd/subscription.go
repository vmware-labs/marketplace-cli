// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(SubscriptionCmd)
	SubscriptionCmd.AddCommand(ListSubscriptionsCmd)
}

var SubscriptionCmd = &cobra.Command{
	Use:       "subscription",
	Aliases:   []string{"subscriptions"},
	Short:     "List subscriptions",
	Long:      "List subscriptions",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{ListSubscriptionsCmd.Use},
}

var ListSubscriptionsCmd = &cobra.Command{
	Use:    "list",
	Short:  "List subscriptions",
	Long:   "Lists subscriptions",
	Args:   cobra.NoArgs,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		subscriptions, err := Marketplace.ListSubscriptions()
		if err != nil {
			return err
		}

		return Output.RenderSubscriptions(subscriptions)
	},
}
