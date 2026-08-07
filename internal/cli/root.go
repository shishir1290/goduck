package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "goduck",
	Short: "Goduck Backend DSL",
	Long:  "Goduck is a backend DSL that generates Go applications.",
}

func Execute() error {
	autoInstallExtensionsIfNeeded()
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installCmd)
}