package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	reposCountCmd = &cobra.Command{
		Use:   "count",
		Short: "Returns the number of repositories found in the repositories.yml file",
		RunE: func(cmd *cobra.Command, args []string) error {

			repositoriesFile, _, err := loadRepositoriesFile(repositoriesFileFl)
			if err != nil {
				log.Fatal(err)
			}

			group, err := repositoriesFile.Group(groupFl)
			if err != nil {
				return err
			}

			fmt.Printf("%d\n", len(group.ListRepositories(childrenFl)))

			return nil
		},
	}
)

func init() {

	AddChildrenFlag(reposCountCmd)
	AddGroupFlag(reposCountCmd)
	AddRepositoriesFileFlag(reposCountCmd)

	reposCmd.AddCommand(reposCountCmd)
}
