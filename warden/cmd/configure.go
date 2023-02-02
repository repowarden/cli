package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	qs = []*survey.Question{
		{
			Name:     "githubToken",
			Prompt:   &survey.Password{Message: "Please enter a GitHub token:"},
			Validate: survey.Required,
		},
	}

	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Store your GitHub token other commands can work",
		RunE: func(cmd *cobra.Command, args []string) error {

			answers := struct {
				GitHubToken string
			}{}

			err := survey.Ask(qs, &answers)
			if err != nil {
				return err
			}

			viper.Set("githubtoken", answers.GitHubToken)

			viper.WriteConfig()

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(configureCmd)
}
