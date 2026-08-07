package cli

import (
	"fmt"
	"os"

	"github.com/shishir1290/goduck/internal/compiler"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build <file.duck>",
	Short: "Build a .duck application",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		filename := args[0]

		if _, err := os.Stat(filename); err != nil {
			return fmt.Errorf("file not found: %s", filename)
		}

		_, err := compiler.Build(filename)
		return err
	},
}