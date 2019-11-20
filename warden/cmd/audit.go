package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

const (
	ERR_DEFAULT_BRANCH = "The repository %s has an incorrect default branch."
	ERR_LICENSE        = "The repository %s has an incorrect license."
	ERR_LABEL_MISSING  = "The repository %s is missing the label %s."
	ERR_LABEL_EXTRA    = "The repository %s has an extra label."
	ERR_ACCESS_MISSING = "The repository %s doesn't have the user."
	ERR_ACCESS_WRONG   = "The repository %s's user's permission is incorrect."
)

type Repo struct {
	org  string
	repo string
}

type UserPermission struct {
	Username   string `yaml:"user"`
	Permission string `yaml:"permission"`
}

type Rule struct {
	Repos          []string         `yaml:"repositories"`
	DefaultBranch  string           `yaml:"defaultBranch"`
	Archived       bool             `yaml:"archived"` // include archived repos in lookup?
	License        string           `yaml:"license"`
	Labels         []string         `yaml:"labels"`
	LabelStrategy  string           `yaml:"labelStrategy"`
	Access         []UserPermission `yaml:"access"`
	AccessStrategy string           `yaml:"accessStrategy"`
}

type RuleError struct {
	repo  Repo
	error string
}

func (re *RuleError) Error() string {
	return fmt.Sprintf(re.error, re.repo.org+"/"+re.repo.repo)
}

var (
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Validates that 1 or more repos meet a set of rules",
		RunE: func(cmd *cobra.Command, args []string) error {

			rawRule, err := ioutil.ReadFile("warden.yml")
			if err != nil {
				log.Fatal(err)
			}

			var rule Rule
			var res []RuleError

			err = yaml.Unmarshal(rawRule, &rule)
			if err != nil {
				log.Fatal(err)
			}

			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: viper.GetString("githubToken")},
			)
			tc := oauth2.NewClient(context.Background(), ts)
			client := github.NewClient(tc)

			for _, repo := range rule.Repos {

				urlPieces := strings.Split(repo, "/")
				org := urlPieces[3]
				name := urlPieces[4]

				repoResp, _, _ := client.Repositories.Get(context.Background(), org, name)

				if *repoResp.Archived != rule.Archived {
					continue
				}

				if *repoResp.DefaultBranch != rule.DefaultBranch {
					res = append(res, RuleError{
						Repo{org: org, repo: name},
						ERR_DEFAULT_BRANCH,
					})

					fmt.Printf("Error: The default branch should be %s, not %s.\n", rule.DefaultBranch, *repoResp.DefaultBranch)
				}

				// if license is to be checked...
				if rule.License != "" && *repoResp.License.Key != rule.License {
					res = append(res, RuleError{
						Repo{org: org, repo: name},
						ERR_LICENSE,
					})
				}

				// if label checks are to happen
				if len(rule.Labels) > 0 {

					labels, _, err := client.Issues.ListLabels(context.Background(), org, name, nil)
					if err != nil {
						return err
					}

					if rule.LabelStrategy == "available" || rule.LabelStrategy == "" {

						// for each labal we're checking for
						for _, label := range rule.Labels {

							found := false

							for _, iLabel := range labels {

								if label == iLabel.GetName() {
									found = true
								}
							}

							if !found {
								res = append(res, RuleError{
									Repo{org: org, repo: name + " label:" + label},
									ERR_LABEL_MISSING,
								})
							}
						}
					} else if rule.LabelStrategy == "only" {

						// for each labal we're checking for
						for _, iLabel := range labels {

							found := false

							for _, label := range rule.Labels {

								if label == iLabel.GetName() {
									found = true
								}
							}

							if !found {
								res = append(res, RuleError{
									Repo{org: org, repo: name + " label:" + iLabel.GetName()},
									ERR_LABEL_EXTRA,
								})
							}
						}
					} else {
						return errors.New("The labelStrategy of " + rule.LabelStrategy + " isn't valid.")
					}

				}

			}

			if len(res) > 0 {

				fmt.Println("The audit failed. Here are the errors:")

				for _, err := range res {
					fmt.Printf("- %s\n", err.Error())
				}
			}

			fmt.Println("The audit completed successfully.")

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(auditCmd)
}
