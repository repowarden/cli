package cmd

import "github.com/spf13/cobra"

var repositoriesFileFl string
var wardenFileFl string

func AddRepositoriesFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&repositoriesFileFl, "repositoriesFile", "", "file containing rules (default is ./repositories.y[a]ml)")
}

func AddWardenFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&wardenFileFl, "wardenFile", "", "file containing rules (default is ./warden.y[a]ml)")
}
