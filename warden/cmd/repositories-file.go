package cmd

type RepositoryDefinition struct {
	URL  string   `yaml:"url"`
	Tags []string `yaml:"tags"`
}

type RepositoryGroup struct {
	Group        string                 `yaml:"group"`
	Repositories []RepositoryDefinition `yaml:"repositories"`
	Children     []RepositoryGroup      `yaml:"children"`
}

// repositories.yml
type RepositoriesFile []RepositoryGroup

// Gets repositories from a repositories.yml file by group. The only group
// supported at this time is 'all'.
func (this RepositoriesFile) RepositoriesByGroup(group string) []RepositoryDefinition {

	var repos []RepositoryDefinition

	if group != "all" {
		return repos
	}

	repos = recurseGroups(this)

	return repos
}

func recurseGroups(groups RepositoriesFile) []RepositoryDefinition {

	var repos []RepositoryDefinition

	// recursively get repos
	for _, group := range groups {

		if len(group.Repositories) > 0 {
			repos = append(repos, group.Repositories...)
		}

		if len(group.Children) > 0 {
			repos = recurseGroups(group.Children)
		}
	}

	return repos
}
