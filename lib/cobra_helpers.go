package lib

import (
	"errors"
	"os"

	. "github.com/spf13/cobra"
)

func FileArg() PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires a single file as an argument")
		}

		_, err := os.Stat(args[0])
		return err
	}
}
