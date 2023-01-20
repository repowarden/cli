package cmd

import "github.com/spf13/cobra"

var policiesFileFl string
var repositoriesFileFl string

func AddPoliciesFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&policiesFileFl, "policiesFile", "", "file containing rules (default is ./policies.y[a]ml)")
}

func AddRepositoriesFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&repositoriesFileFl, "repositoriesFile", "", "file containing rules (default is ./repositories.y[a]ml)")
}
