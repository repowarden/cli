package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var (
	reposListCmd = &cobra.Command{
		Use:   "list",
		Short: "List repositories found in the repositories.yml file",
		RunE: func(cmd *cobra.Command, args []string) error {

			repositoriesFile, _, err := loadRepositoriesFile(repositoriesFileFl)
			if err != nil {
				log.Fatal(err)
			}

			group, err := repositoriesFile.Group(groupFl)
			if err != nil {
				return err
			}

			repos := group.ListRepositories(childrenFl)

			if len(repos) == 0 {
				return errors.New("No repositories matched.")
			}

			fmt.Printf("%s\n", strings.Join(repos, "\n"))

			return nil
		},
	}
)

func init() {

	AddChildrenFlag(reposListCmd)
	AddGroupFlag(reposListCmd)
	AddRepositoriesFileFlag(reposListCmd)

	reposCmd.AddCommand(reposListCmd)
}
