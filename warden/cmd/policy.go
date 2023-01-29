package cmd

import (
	"github.com/spf13/cobra"
)

var (
	policyCmd = &cobra.Command{
		Use:   "policy",
		Short: "Subcommands for policy.yml files",
	}
)

func init() {
	rootCmd.AddCommand(policyCmd)
}
