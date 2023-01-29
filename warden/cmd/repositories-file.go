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

//=============================================================================
// A list of repositories belonging to a group, the top-level object in a
// repositories.yml file.
//=============================================================================
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

//=============================================================================
// A repositories.yml file
//=============================================================================
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

	content, err := yaml.Marshal(this)
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
