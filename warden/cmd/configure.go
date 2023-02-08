package cmd

import (
	"os"

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
		Long: `Store your GitHub token other commands can work,
Optionally, the environment variable 'RW_GH_TOKEN' can be set. This is useful in CI environments.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			answers := struct {
				GitHubToken string
			}{}

			err := survey.Ask(qs, &answers)
			if err != nil {
				return err
			}

			// makes sure the path to the config file exists
			err = os.MkdirAll(os.ExpandEnv("$HOME/.config/warden"), os.ModePerm)
			if err != nil {
				return err
			}

			viper.Set("GH_TOKEN", answers.GitHubToken)

			err = viper.WriteConfig()
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(configureCmd)
}
