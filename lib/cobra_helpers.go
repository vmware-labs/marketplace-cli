// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package lib

import (
	"fmt"
	"os"

	. "github.com/spf13/cobra"
)

func FileArg() PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires a single file as an argument")
		}

		_, err := os.Stat(args[0])
		return err
	}
}
