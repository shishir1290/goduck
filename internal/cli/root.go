package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:     "goduck",
	Short:   "Goduck Backend DSL",
	Long:    "Goduck is a backend DSL that generates Go applications.",
	Version: "v0.0.1",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Goduck",
	Long:  `All software has versions. This is Goduck's`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("goduck version %s\n", rootCmd.Version)
	},
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
	rootCmd.AddCommand(versionCmd)
}