package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v50/github"
	"github.com/repowarden/cli/warden/vcsurl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

type PolicyFile struct {
	DefaultBranch  string           `yaml:"defaultBranch"`
	Archived       bool             `yaml:"archived"` // include archived repos in lookup?
	License        *LicenseRule     `yaml:"license"`
	Labels         []string         `yaml:"labels"`
	LabelStrategy  string           `yaml:"labelStrategy"`
	Access         []UserPermission `yaml:"access"`
	AccessStrategy string           `yaml:"accessStrategy"`
}

var (
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Validates that 1 or more repos meet a set of policy",
		RunE: func(cmd *cobra.Command, args []string) error {

			var policyErrors []PolicyError

			repoFile, _, err := loadRepositoriesFile(repositoriesFileFl)
			if err != nil {
				log.Fatal(err)
			}

			policy, _, err := loadPolicyFile(policyFileFl)
			if err != nil {
				log.Fatal(err)
			}

			ghToken := viper.GetString("githubToken")
			if ghToken == "" {
				return errors.New("GitHub credentials were not found. Please run `warden configure`.")
			}

			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: viper.GetString("githubToken")},
			)
			tc := oauth2.NewClient(context.Background(), ts)
			client := github.NewClient(tc)

			for _, repoDef := range repoFile.RepositoriesByGroup("all") {

				repo, err := vcsurl.Parse(repoDef.URL)
				if err != nil {
					return fmt.Errorf("The following URL cannot be parsed: %s", repoDef.URL)
				}

				repoResp, _, _ := client.Repositories.Get(context.Background(), repo.Owner, repo.Name)

				if repoResp.GetArchived() != policy.Archived {
					continue
				}

				if repoResp.GetDefaultBranch() != policy.DefaultBranch {
					policyErrors = append(policyErrors, PolicyError{
						repoDef,
						ERR_BRANCH_DEFAULT,
						[]any{policy.DefaultBranch, repoResp.GetDefaultBranch()},
					})
				}

				// if license is to be checked...
				if policy.License != nil && policy.License.Scope == repoResp.GetVisibility() || policy.License.Scope == "all" {
					if !slices.Contains(policy.License.Names, repoResp.GetLicense().GetKey()) {
						policyErrors = append(policyErrors, PolicyError{
							repoDef,
							ERR_LICENSE,
							[]any{policy.License.Names, repoResp.GetLicense().GetKey()},
						})
					}
				}

				// if label checks are to happen
				if len(policy.Labels) > 0 {

					labels, _, err := client.Issues.ListLabels(context.Background(), repo.Owner, repo.Name, nil)
					if err != nil {
						return err
					}

					if policy.LabelStrategy == "available" || policy.LabelStrategy == "" {

						// for each labal we're checking for
						for _, label := range policy.Labels {

							found := false

							for _, iLabel := range labels {

								if label == iLabel.GetName() {
									found = true
								}
							}

							if !found {
								policyErrors = append(policyErrors, PolicyError{
									repoDef,
									ERR_LABEL_MISSING,
									[]any{label},
								})
							}
						}
					} else if policy.LabelStrategy == "only" {

						// for each labal we're checking for
						for _, iLabel := range labels {

							found := ""

							for _, label := range policy.Labels {

								if label == iLabel.GetName() {
									found = label
								}
							}

							if found != "" {
								policyErrors = append(policyErrors, PolicyError{
									repoDef,
									ERR_LABEL_EXTRA,
									[]any{found},
								})
							}
						}
					} else {
						return errors.New("The labelStrategy of " + policy.LabelStrategy + " isn't valid.")
					}

				}

				// if access permissions are to be checked...
				if len(policy.Access) > 0 {

					teams, _, err := client.Repositories.ListTeams(context.Background(), repo.Owner, repo.Name, nil)
					if err != nil {
						return err
					}

					if policy.AccessStrategy == "available" || policy.AccessStrategy == "" {

						// for each team we're checking for
						for _, user := range policy.Access {

							found := ""
							matched := ""

							for _, team := range teams {

								if user.Username == team.GetName() {

									found = user.Username

									if user.Permission == team.GetPermission() {
										matched = user.Permission
									}
								}
							}

							if found != "" {
								policyErrors = append(policyErrors, PolicyError{
									repoDef,
									ERR_ACCESS_MISSING,
									[]any{found},
								})
							} else if matched != "" {
								policyErrors = append(policyErrors, PolicyError{
									repoDef,
									ERR_ACCESS_WRONG,
									[]any{matched},
								})
							}
						}
					} else {
						return errors.New("The accessStrategy of " + policy.AccessStrategy + " isn't valid.")
					}
				}
			}

			if len(policyErrors) > 0 {

				fmt.Printf("The audit failed. Here are the errors:\n\n")

				var curRepo string

				for _, err := range policyErrors {

					if curRepo != err.repository.URL {

						curRepo = err.repository.URL
						fmt.Printf("%s:\n", curRepo)
					}

					fmt.Printf("  - %s\n", err.Error())
				}

				return nil
			}

			fmt.Println("The audit completed successfully.")

			return nil
		},
	}
)

func init() {

	AddPolicyFileFlag(auditCmd)
	AddRepositoriesFileFlag(auditCmd)

	rootCmd.AddCommand(auditCmd)
}
