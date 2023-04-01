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
	branchFl string

	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Validates that 1 or more repos meet a set of policy",
		RunE: func(cmd *cobra.Command, args []string) error {

			var results auditResults

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

				var currentBranch string

				repoResp, _, _ := client.Repositories.Get(context.Background(), repo.Owner, repo.Name)

				if repoResp.GetArchived() != policy.Archived {
					continue
				}

				// check if we're working with the default branch or a specific one
				if branchFl != "" {
					currentBranch = branchFl
				} else {
					currentBranch = repoResp.GetDefaultBranch()
				}

				if repoResp.GetDefaultBranch() != policy.DefaultBranch {
					results.add(
						repo,
						RESULT_ERROR,
						ERR_BRANCH_DEFAULT,
						policy.DefaultBranch,
						repoResp.GetDefaultBranch(),
					)
				}

				// if license is to be checked...
				if policy.License != nil && policy.License.Scope == repoResp.GetVisibility() || policy.License.Scope == "all" {
					if repoResp.GetLicense().GetKey() == "" {
						results.add(
							repo,
							RESULT_ERROR,
							ERR_LICENSE_MISSING,
						)
					} else if !slices.Contains(policy.License.Names, repoResp.GetLicense().GetKey()) {
						results.add(
							repo,
							RESULT_ERROR,
							ERR_LICENSE_DIFFERENT,
							policy.License.Names,
							repoResp.GetLicense().GetKey(),
						)
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
								results.add(
									repo,
									RESULT_ERROR,
									ERR_LABEL_MISSING,
									label,
								)
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
								results.add(
									repo,
									RESULT_ERROR,
									ERR_LABEL_EXTRA,
									found,
								)
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
						results.add(repo, RESULT_WARNING, fmt.Sprintf("Couldn't pull teams. There's a visibility issue here."))

						// skip the rest
						policy.Access = nil
					} else if err != nil {
						return err
					}

					for _, accessPolicy := range policy.Access {
						results.merge(auditAccessPolicy(accessPolicy, repo, teams))
					}
				}

				// if codeowners are to to be checked...
				if len(policy.Codeowners) > 0 {

					for _, coPolicy := range policy.Codeowners {
						results.merge(auditCodeownersPolicy(coPolicy, repo, client, currentBranch))
					}
				}
			}

			fmt.Printf(
				`======================================================================
                         Warden Audit Results

     errors: %d     warnings: %d     repos: %d     group: %s
======================================================================

`,
				len(results.ByType(RESULT_ERROR)),
				len(results.ByType(RESULT_WARNING)),
				len(repos),
				groupFl,
			)

			if len(results) > 0 {

				var curRepo string

				for _, result := range results {

					// print repo URL whenever we move to the next one
					if curRepo != result.repository.ToHTTPS() {

						curRepo = result.repository.ToHTTPS()
						fmt.Fprintf(os.Stderr, "%s:\n", curRepo)
					}

					switch result.resultType {
					case RESULT_ERROR:
						fmt.Fprintf(os.Stderr, "  \033[31mx\033[0m %s\n", result)
					case RESULT_WARNING:
						fmt.Printf("  \033[33mo\033[0m %s\n", result)
					}
				}

				fmt.Println("") // intentional
			}

			if len(results.ByType(RESULT_ERROR)) > 0 {
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

	auditCmd.PersistentFlags().StringVar(&branchFl, "branch", "", "git branch to audit (for applicable polcies")

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
