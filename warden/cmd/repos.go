package cmd

import (
	"github.com/spf13/cobra"
)

var (
	reposCmd = &cobra.Command{
		Use:   "repos",
		Short: "Subcommands for repositories.yml files",
	}
)

func init() {
	rootCmd.AddCommand(reposCmd)
}
