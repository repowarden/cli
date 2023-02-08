package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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

func (this *UserPermission) GetUsername() string {

	if this.IsUser() {
		return this.Username
	}

	return this.Username[this.SlashPos()+1 : len(this.Username)]
}

func (this *UserPermission) IsTeam() bool {
	return !this.IsUser()
}

func (this *UserPermission) IsUser() bool {
	return this.SlashPos() == -1
}

func (this *UserPermission) Org() string {

	if this.IsUser() {
		return this.Username
	}

	return this.Username[0:this.SlashPos()]
}

func (this *UserPermission) SlashPos() int {
	return strings.Index(this.Username, "/")
}

type PolicyFile struct {
	DefaultBranch  string           `yaml:"defaultBranch"`
	Archived       bool             `yaml:"archived"` // include archived repos in lookup?
	License        *LicenseRule     `yaml:"license"`
	Labels         []string         `yaml:"labels"`
	LabelStrategy  string           `yaml:"labelStrategy"`
	Access         []UserPermission `yaml:"access"`
	AccessStrategy string           `yaml:"accessStrategy"`
	CodeOwners     string           `yaml:"codeowners"`
}

var (
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Validates that 1 or more repos meet a set of policy",
		RunE: func(cmd *cobra.Command, args []string) error {

			var policyErrors []PolicyError

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

			for _, repoDef := range group.GetRepositories(childrenFl) {

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
					if repoResp.GetLicense().GetKey() == "" {
						policyErrors = append(policyErrors, PolicyError{
							repoDef,
							ERR_LICENSE_MISSING,
							nil,
						})
					} else if !slices.Contains(policy.License.Names, repoResp.GetLicense().GetKey()) {
						policyErrors = append(policyErrors, PolicyError{
							repoDef,
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

					if !slices.Contains([]string{"available", "only", ""}, policy.AccessStrategy) {
						return errors.New("The accessStrategy of " + policy.AccessStrategy + " isn't valid.")
					}

					// for each user/team we're checking for
					for _, user := range policy.Access {

						found := ""
						matched := ""
						onlyMatches := make(map[string]bool)

						// only checking teams for now
						if user.IsUser() {
							continue
						}

						// for teams, the team check only matters if we're in the same org
						if user.Org() != repo.Owner {
							continue
						}

						for _, team := range teams {

							fullTeamName := repo.Owner + "/" + team.GetSlug()

							if user.GetUsername() == team.GetSlug() {

								found = user.Username
								onlyMatches[fullTeamName] = true

								if user.Permission != team.GetPermission() {
									matched = team.GetPermission()
								}
							} else {
								foundAlready, err := onlyMatches[fullTeamName]
								if !foundAlready || err {
									onlyMatches[fullTeamName] = false
								}
							}
						}

						if found == "" {
							policyErrors = append(policyErrors, PolicyError{
								repoDef,
								ERR_ACCESS_MISSING,
								[]any{
									"team",
									user.Username,
								},
							})
						} else if matched != "" {
							policyErrors = append(policyErrors, PolicyError{
								repoDef,
								ERR_ACCESS_DIFFERENT,
								[]any{
									found,
									user.Permission,
									matched,
								},
							})
						}

						if policy.AccessStrategy == "only" {

							for team, _ := range onlyMatches {

								fmt.Printf("The team is: %s\n", team) //DEBUG
								if onlyMatches[team] == false {
									policyErrors = append(policyErrors, PolicyError{
										repoDef,
										ERR_ACCESS_EXTRA,
										[]any{
											team,
										},
									})
								}
							}
						}
					}
				}

				// if the CODEOWNERS file is to be checked...
				if policy.CodeOwners != "" {

					file, _, _, err := client.Repositories.GetContents(context.Background(), repo.Owner, repo.Name, ".github/CODEOWNERS", nil)
					if err != nil {

						switch err.(type) {
						case *github.ErrorResponse:
							if err.(*github.ErrorResponse).Response != nil && err.(*github.ErrorResponse).Response.StatusCode == 404 {
								policyErrors = append(policyErrors, PolicyError{
									repoDef,
									ERR_CO_MISSING,
									nil,
								})

								continue
							} else {
								return err
							}
						default:
							return err
						}
					}

					content, err := file.GetContent()
					if err != nil {
						return err
					}
					// handle manual tabs
					policyContent := fmt.Sprintf(policy.CodeOwners)

					// check if the files match
					if policyContent != content {
						policyErrors = append(policyErrors, PolicyError{
							repoDef,
							ERR_CO_DIFFERENT,
							nil,
						})
					}

					// check for codeowners syntax errors
					coErrs, _, err := client.Repositories.GetCodeownersErrors(context.Background(), repo.Owner, repo.Name)
					if err != nil {
						return err
					}

					if len(coErrs.Errors) > 0 {

						var suggestions []string
						for _, coErr := range coErrs.Errors {
							suggestions = append(suggestions, "    > "+coErr.GetSuggestion())
						}

						policyErrors = append(policyErrors, PolicyError{
							repoDef,
							ERR_CO_SYNTAX,
							[]any{strings.Join(suggestions, "\n")},
						})
					}
				}
			}

			if len(policyErrors) > 0 {

				fmt.Fprintf(os.Stderr, "The audit failed.\n\n")

				var curRepo string

				for _, err := range policyErrors {

					if curRepo != err.repository.URL {

						curRepo = err.repository.URL
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
