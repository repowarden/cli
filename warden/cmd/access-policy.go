package cmd

import (
	"strings"

	"golang.org/x/exp/slices"

	"github.com/google/go-github/v53/github"
)

// The list of users/teams, their permissions, and a strategy that should be applied.
type accessPolicy struct {
	Strategy    string           `yaml:"strategy"`
	Permissions []userPermission `yaml:"permissions"`
	Tags        []string         `yaml:"tags"`
}

// A user/team & permission pairing
type userPermission struct {
	User       string `yaml:"user"`
	Permission string `yaml:"permission"`
}

// Whether or not this userPermission is for a team
func (this *userPermission) IsTeam() bool {
	return !this.IsUser()
}

// Whether or not this userPermission is for a user
func (this *userPermission) IsUser() bool {
	return this.SlashPos() == -1
}

// Return the user's own org or the teams org
func (this *userPermission) Owner() string {

	if this.IsUser() {
		return this.User
	}

	return this.User[0:this.SlashPos()]
}

// Returns the string index of the slash, which will be -1 when this is a user and not a team
func (this *userPermission) SlashPos() int {
	return strings.Index(this.User, "/")
}

// Returns the user's username or full team name
func (this *userPermission) UserSlug() string {

	if this.IsUser() {
		return this.User
	}

	return this.User[this.SlashPos()+1 : len(this.User)]
}

// Does the work to check an access policy against a repository
func auditAccessPolicy(policy accessPolicy, repo *wardenRepo, teams []*github.Team) auditResults {

	var results auditResults

	if !tagsMatched(policy.Tags, repo.Tags()) {
		return nil
	}

	if !slices.Contains([]string{"available", "only", ""}, policy.Strategy) {
		results.add(
			repo,
			RESULT_ERROR,
			ERR_ACCESS_STRATEGY,
			policy.Strategy,
		)

		return results
	}

	onlyMatches := make(map[string]bool)

	// for each user/team we're checking for
	for _, user := range policy.Permissions {

		found := ""
		matched := ""

		// only checking teams for now
		if user.IsUser() {
			continue
		}

		// for teams, the team check only matters if we're in the same org
		if user.Owner() != repo.Owner {
			continue
		}

		for _, team := range teams {

			fullTeamName := strings.TrimSpace(repo.Owner + "/" + team.GetSlug())

			if user.UserSlug() == team.GetSlug() {

				found = user.UserSlug()
				onlyMatches[fullTeamName] = true

				if user.Permission != team.GetPermission() {
					matched = team.GetPermission()
				}
			} else {

				if onlyMatches[fullTeamName] != true {
					onlyMatches[fullTeamName] = false
				}
			}
		}

		if found == "" {
			results.add(
				repo,
				RESULT_ERROR,
				ERR_ACCESS_MISSING,
				user.UserSlug(),
			)
		} else if matched != "" {
			results.add(
				repo,
				RESULT_ERROR,
				ERR_ACCESS_DIFFERENT,
				found,
				user.Permission,
				matched,
			)
		}

	}

	if policy.Strategy == "only" {

		for team, _ := range onlyMatches {

			if onlyMatches[team] == false {
				results.add(
					repo,
					RESULT_ERROR,
					ERR_ACCESS_EXTRA,
					team,
				)
			}
		}

	}

	return results
}
