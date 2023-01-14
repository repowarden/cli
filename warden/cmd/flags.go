package cmd

import "github.com/spf13/cobra"

var wardenFileFl string

func AddWardenFileFlag(cmd *cobra.Command) {

	cmd.PersistentFlags().StringVar(&wardenFileFl, "wardenFile", "", "file containing rules (default is ./warden.y[a]ml)")
}
