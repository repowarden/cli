package cmd

import "github.com/spf13/cobra"

var policyFileFl string
var repositoriesFileFl string

func AddPolicyFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&policyFileFl, "policyFile", "", "file containing rules (default is ./policy.y[a]ml)")
}

func AddRepositoriesFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&repositoriesFileFl, "repositoriesFile", "", "file containing rules (default is ./repositories.y[a]ml)")
}
