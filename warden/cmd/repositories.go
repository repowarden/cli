package cmd

import "github.com/repowarden/cli/warden/vcsurl"

// There are other Repository types scattered around the codebase, but this should be the main when dealing with the core business logic.
type wardenRepo struct {
	*vcsurl.Repository
	tags []string
}

// Returns a string slice of tags. Generated tags such as org are injected into the response.
func (this *wardenRepo) Tags() []string {
	return append(this.tags, this.Owner)
}

// Create a new WardenRepo
func WardenRepo(repo *vcsurl.Repository, tags []string) *wardenRepo {

	return &wardenRepo{
		repo,
		tags,
	}
}

// Create a slice of new WardenRepos from RepoDefs
func WardenRepos(repoDefs []RepositoryDefinition) ([]*wardenRepo, error) {

	var repos []*wardenRepo

	for _, repoDef := range repoDefs {

		repo, err := vcsurl.Parse(repoDef.URL)
		if err != nil {
			return nil, err
		}

		repos = append(repos, WardenRepo(
			repo,
			repoDef.Tags,
		))
	}

	return repos, nil
}
