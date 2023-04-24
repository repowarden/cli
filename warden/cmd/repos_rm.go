package cmd

import (
	"fmt"
	"log"

	"github.com/repowarden/cli/warden/vcsurl"
	"github.com/spf13/cobra"
)

var (
	reposRMCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a repository from the repositories.yml file",
		RunE: func(cmd *cobra.Command, args []string) error {

			repository := args[0]

			// ensure the repository URL is valid
			_, err := vcsurl.Parse(repository)
			if err != nil {
				return fmt.Errorf("The repository URL %s is invalid: %s", repository, err)
			}

			repositoriesFile, _, err := loadRepositoriesFile(repositoriesFileFl)
			if err != nil {
				log.Fatal(err)
			}

			group, err := repositoriesFile.Group("all")
			if err != nil {
				return err
			}

			removed := group.Remove(RepositoryDefinition{URL: repository})
			if !removed {
				fmt.Println("No change was made. Couldn't find repository.")
				return nil
			}

			filepath, err := repositoriesFile.save(repositoriesFileFl, false)
			if err != nil {
				return err
			}

			fmt.Printf("The repositories file %s has been updated.", filepath)

			return nil
		},
	}
)

func init() {

	AddRepositoriesFileFlag(reposRMCmd)

	reposCmd.AddCommand(reposRMCmd)
}
