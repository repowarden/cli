package cmd

import "github.com/spf13/cobra"

var childrenFl bool

func AddChildrenFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().BoolVar(&childrenFl, "children", true, "whether or not to include child groups, default 'true'")
}

var groupFl string

func AddGroupFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&groupFl, "group", "all", "which group to filter repositories by, default 'all'")
}

var policyFileFl string

func AddPolicyFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&policyFileFl, "policyFile", "", "file containing rules (default is ./policy.y[a]ml)")
}

var repositoriesFileFl string

func AddRepositoriesFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&repositoriesFileFl, "repositoriesFile", "", "file containing rules (default is ./repositories.y[a]ml)")
}
