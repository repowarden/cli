package cmd

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/repowarden/cli/warden/vcsurl"
	"github.com/spf13/cobra"
)

var groupFl string

var (
	reposAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a repository to the repositories.yml file",
		RunE: func(cmd *cobra.Command, args []string) error {

			repository := args[0]

			// ensure the repository URL is valid
			_, err := vcsurl.Parse(repository)
			if err != nil {
				return fmt.Errorf("The repository URL %s isn't valid.", repository)
			}

			repositoriesFile, _, err := loadRepositoriesFile(repositoriesFileFl)
			if err != nil {
				log.Fatal(err)
			}

			group, err := repositoriesFile.Group(groupFl)
			if err != nil {
				return err
			}

			group.Add(RepositoryDefinition{URL: repository})

			filepath, err := repositoriesFile.save(repositoriesFileFl, false)
			if err != nil {
				return err
			}

			fmt.Printf("The repositories file %s has been created/updated.", filepath)

			return nil
		},
	}
)

func init() {

	reposAddCmd.PersistentFlags().StringVar(&groupFl, "group", "", "which group the repository belongs to")
	reposAddCmd.MarkFlagRequired("group")
	AddRepositoriesFileFlag(reposAddCmd)

	reposCmd.AddCommand(reposAddCmd)
}
