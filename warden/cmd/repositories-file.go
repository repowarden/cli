package cmd

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type RepositoryDefinition struct {
	URL  string   `yaml:"url"`
	Tags []string `yaml:"tags,omitempty"`
}

// =============================================================================
// A list of repositories belonging to a group, the top-level object in a
// repositories.yml file.
// =============================================================================
type RepositoryGroup struct {
	Group        string                 `yaml:"group"`
	Repositories []RepositoryDefinition `yaml:"repositories"`
	Children     RepositoriesFile       `yaml:"children,omitempty"`
}

// Add a repository to the group
func (this *RepositoryGroup) Add(repoDef RepositoryDefinition) {
	this.Repositories = append(this.Repositories, repoDef)
}

func (this *RepositoryGroup) HasChildren() bool {
	return len(this.Children) > 0
}

// Returns the list of URLs for each repository in this group. The parameter
// decides if to include children or not.
func (this *RepositoryGroup) ListRepositories(listChildren bool) []string {

	var repos []string

	for _, repo := range this.Repositories {
		repos = append(repos, repo.URL)
	}

	if listChildren {
		for _, group := range this.Children {
			repos = append(repos, group.ListRepositories(true)...)
		}
	}

	return repos
}

// Returns the repositories belong to the group. The parameter decides if to
// include children or not.
func (this *RepositoryGroup) GetRepositories(listChildren bool) []RepositoryDefinition {

	repos := this.Repositories

	if listChildren {
		for _, group := range this.Children {
			repos = append(repos, group.GetRepositories(true)...)
		}
	}

	return repos
}

// Remove a repository from the group
func (this *RepositoryGroup) Remove(repoDef RepositoryDefinition) bool {

	repos := this.Repositories

	for i, repo := range repos {

		if repoDef.URL == repo.URL {

			repos[i] = repos[len(repos)-1]
			repos = repos[:len(repos)-1]
			return true
		}
	}

	for _, group := range this.Children {
		if group.Remove(repoDef) {
			return true
		}
	}

	return false
}

// =============================================================================
// A repositories.yml file
// =============================================================================
type RepositoriesFile []*RepositoryGroup

// Get a group from a repositories file.
func (this RepositoriesFile) Group(groupName string) (*RepositoryGroup, error) {

	// check at the current level
	for _, group := range this {
		if group.Group == groupName {
			return group, nil
		}
	}

	// otherwise check the children
	for _, group := range this {
		if group.HasChildren() {
			return group.Children.Group(groupName)
		}
	}

	return nil, errors.New("Group doesn't exist.")
}

// Gets repositories from a repositories.yml file by group. The only group
// supported at this time is 'all'.
func (this *RepositoriesFile) RepositoriesByGroup(group string) []RepositoryDefinition {

	var repos []RepositoryDefinition

	if group != "all" {
		return repos
	}

	repos = recurseGroups(this)

	return repos
}

// saveRepositoriesFile tries to intelligently choose the filepath for the
// repositories file to be saved to. If customPath is not empty, that will be
// the filepath choosen. Unless 'create' is true, this will only try to
// override an existing file.
func (this RepositoriesFile) save(customPath string, create bool) (string, error) {

	// we never actually want the all group so remove it
	all, err := this.Group("all")
	if err != nil {
		return "", err
	}

	content, err := yaml.Marshal(all.Children)
	if err != nil {
		return "", fmt.Errorf("Unable to create YAML from repositories data. Something is wrong.")
	}

	return saveYAMLFile(content, create, customPath, "repositories.yml")
}

//=============================================================================
// Helper functions
//=============================================================================

func recurseGroups(groups *RepositoriesFile) []RepositoryDefinition {

	var repos []RepositoryDefinition

	// recursively get repos
	for _, group := range *groups {

		if len(group.Repositories) > 0 {
			repos = append(repos, group.Repositories...)
		}

		if group.HasChildren() {
			repos = recurseGroups(&group.Children)
		}
	}

	return repos
}
