package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v48/github"
	"github.com/repowarden/cli/warden/vcsurl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

type LicenseRule struct {
	Scope string   `yaml:"scope"`
	Names []string `yaml:"names"`
}

type UserPermission struct {
	Username   string `yaml:"user"`
	Permission string `yaml:"permission"`
}

type Rule struct {
	Repos          []string         `yaml:"repositories"`
	DefaultBranch  string           `yaml:"defaultBranch"`
	Archived       bool             `yaml:"archived"` // include archived repos in lookup?
	License        *LicenseRule     `yaml:"license"`
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

			rawRule, err := loadWardenFile(wardenFileFl)
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

			for _, repoURL := range rule.Repos {

				repo, err := vcsurl.Parse(repoURL)
				if err != nil {
					return fmt.Errorf("The following URL cannot be parsed: %s", repoURL)
				}

				repoResp, _, _ := client.Repositories.Get(context.Background(), repo.Owner, repo.Name)

				if repoResp.GetArchived() != rule.Archived {
					continue
				}

				if repoResp.GetDefaultBranch() != rule.DefaultBranch {
					res = append(res, RuleError{
						Repo{org: repo.Owner, repo: repo.Name},
						ERR_DEFAULT_BRANCH,
					})

					fmt.Printf("Error: The default branch should be %s, not %s.\n", rule.DefaultBranch, repoResp.GetDefaultBranch())
				}

				// if license is to be checked...
				if rule.License != nil && rule.License.Scope == repoResp.GetVisibility() || rule.License.Scope == "all" {
					if !slices.Contains(rule.License.Names, repoResp.GetLicense().GetKey()) {
						res = append(res, RuleError{
							Repo{org: repo.Owner, repo: repo.Name},
							ERR_LICENSE,
						})
					}
				}

				// if label checks are to happen
				if len(rule.Labels) > 0 {

					labels, _, err := client.Issues.ListLabels(context.Background(), repo.Owner, repo.Name, nil)
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
									Repo{org: repo.Owner, repo: repo.Name + " label:" + label},
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
									Repo{org: repo.Owner, repo: repo.Name + " label:" + iLabel.GetName()},
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

	AddWardenFileFlag(auditCmd)

	rootCmd.AddCommand(auditCmd)
}
