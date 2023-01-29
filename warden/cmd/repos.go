package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

var (
	reposCmd = &cobra.Command{
		Use:   "repos",
		Short: "Subcommands for repositories.yml files",
	}
)

func init() {

	AddRepositoriesFileFlag(reposCmd)

	rootCmd.AddCommand(reposCmd)
}
