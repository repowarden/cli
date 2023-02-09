package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v50/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Validates that 1 or more repos meet a set of policy",
		RunE: func(cmd *cobra.Command, args []string) error {

			var policyErrors []policyError

			repoFile, _, err := loadRepositoriesFile(repositoriesFileFl)
			if err != nil {
				return err
			}

			policy, _, err := loadPolicyFile(policyFileFl)
			if err != nil {
				return err
			}

			ghToken := viper.GetString("GH_TOKEN")
			if ghToken == "" {
				return errors.New("GitHub credentials were not found. Please run `warden configure`.")
			}

			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: ghToken},
			)
			tc := oauth2.NewClient(context.Background(), ts)
			client := github.NewClient(tc)

			group, err := repoFile.Group(groupFl)
			if err != nil {
				return err
			}

			repos, err := WardenRepos(group.GetRepositories(childrenFl))
			if err != nil {
				return err
			}

			for _, repo := range repos {

				repoResp, _, _ := client.Repositories.Get(context.Background(), repo.Owner, repo.Name)

				if repoResp.GetArchived() != policy.Archived {
					continue
				}

				if repoResp.GetDefaultBranch() != policy.DefaultBranch {
					policyErrors = append(policyErrors, policyError{
						repo,
						ERR_BRANCH_DEFAULT,
						[]any{policy.DefaultBranch, repoResp.GetDefaultBranch()},
					})
				}

				// if license is to be checked...
				if policy.License != nil && policy.License.Scope == repoResp.GetVisibility() || policy.License.Scope == "all" {
					if repoResp.GetLicense().GetKey() == "" {
						policyErrors = append(policyErrors, policyError{
							repo,
							ERR_LICENSE_MISSING,
							nil,
						})
					} else if !slices.Contains(policy.License.Names, repoResp.GetLicense().GetKey()) {
						policyErrors = append(policyErrors, policyError{
							repo,
							ERR_LICENSE_DIFFERENT,
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
								policyErrors = append(policyErrors, policyError{
									repo,
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
								policyErrors = append(policyErrors, policyError{
									repo,
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

					teams, resp, err := client.Repositories.ListTeams(context.Background(), repo.Owner, repo.Name, nil)
					if resp.StatusCode == 404 {

						// considering this repo worked for other audits but not this, this likely
						// means we don't have admin access in order to check teams
						fmt.Fprintf(os.Stderr, "Error: couldn't pull the teams for %s.\nThis is likely a permission issue with the token being used to run Warden. If\nthe user whose token is being used doesn't have admin access\nto the repo, teams can't be pulled.\n\n", repo.ToHTTPS())

						// skip the rest
						policy.Access = nil
					} else if err != nil {
						return err
					}

					for _, accessPolicy := range policy.Access {
						policyErrors = append(policyErrors, auditAccessPolicy(accessPolicy, repo, teams)...)
					}
				}

				// if codeowners are to to be checked...
				if len(policy.Codeowners) > 0 {

					for _, coPolicy := range policy.Codeowners {
						policyErrors = append(policyErrors, auditCodeownersPolicy(coPolicy, repo, client)...)
					}
				}
			}

			if len(policyErrors) > 0 {

				fmt.Fprintf(os.Stderr, "The audit failed.\n\n")

				var curRepo string

				for _, err := range policyErrors {

					if curRepo != err.repository.ToHTTPS() {

						curRepo = err.repository.ToHTTPS()
						fmt.Fprintf(os.Stderr, "%s:\n", curRepo)
					}

					fmt.Fprintf(os.Stderr, "  - %s\n", err.Error())
				}

				fmt.Println("") // intentional

				return fmt.Errorf("The audit failed. Above are the policy failures, by repository.\n")
			}

			fmt.Println("The audit completed successfully.")

			return nil
		},
	}
)

func init() {

	AddChildrenFlag(auditCmd)
	AddGroupFlag(auditCmd)
	AddPolicyFileFlag(auditCmd)
	AddRepositoriesFileFlag(auditCmd)

	rootCmd.AddCommand(auditCmd)
}

func tagsMatched(policyTags, repoTags []string) bool {

	// if policy doesn't specify tags, then all repos are allowed
	if len(policyTags) == 0 {
		return true
	}

	for _, tag := range policyTags {
		if slices.Contains(repoTags, tag) {
			return true
		}
	}

	return false
}
