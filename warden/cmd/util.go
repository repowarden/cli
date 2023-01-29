package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// loadPolicyFile tries to intelligently choose a filepath for the
// policy file and then return the unmarshalled struct. If customPath is not
// empty, it will try to use that before the default filenames.
func loadPolicyFile(customPath string) (*PolicyFile, []byte, error) {

	var file PolicyFile
	var yamlContent []byte
	var err error

	yamlContent, err = loadYAMLFile("policy.yml", customPath)
	if err != nil {
		return nil, nil, err
	}

	err = yaml.Unmarshal(yamlContent, &file)
	if err != nil {
		return nil, nil, err
	}

	return &file, yamlContent, nil
}

// loadRepositoriesFile tries to intelligently choose a filepath for the
// Wardenfile and then return the unmarshalled struct. If customPath is not
// empty, it will try to use that before the default filenames.
func loadRepositoriesFile(customPath string) (RepositoriesFile, []byte, error) {

	var repositoriesFile []RepositoryGroup
	var yamlContent []byte
	var err error

	yamlContent, err = loadYAMLFile("repositories.yml", customPath)
	if err != nil {
		return nil, nil, fmt.Errorf("A Repositories file was not found. Either './repositories.yml' needs to be used or the '--repositoriesFile' flag set.")
	}

	err = yaml.Unmarshal(yamlContent, &repositoriesFile)
	if err != nil {
		return nil, nil, fmt.Errorf("The repositories file couldn't be parsed. Something is wrong.")
	}

	return repositoriesFile, yamlContent, nil
}

// loadYAMLFile attempts to load a YAML file based on one or more possible file
// names. Both .yml and .yaml will be attempted and in that order.
func loadYAMLFile(filepaths ...string) ([]byte, error) {

	var yamlContent []byte
	var possiblePaths []string
	var err error

	if len(filepaths) == 0 {
		return nil, fmt.Errorf("At least one filepath needs to be provided.")
	}

	for _, path := range filepaths {

		if strings.HasSuffix(path, ".yml") {
			path = path[0 : len(path)-3]
		} else if strings.HasSuffix(path, ".yaml") {
			path = path[0 : len(path)-4]
		} else if path == "" {
			continue
		} else {
			return nil, fmt.Errorf("Only YAML files are supported.")
		}

		possiblePaths = append(possiblePaths, path+"yml")
		possiblePaths = append(possiblePaths, path+"yaml")
	}

	for _, path := range filepaths {

		yamlContent, err = os.ReadFile(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if errors.Is(err, fs.ErrPermission) {
			return nil, fmt.Errorf("Warden doesn't have permission to open %s", path)
		} else if err != nil {
			return nil, err
		}

		break
	}

	if len(yamlContent) == 0 {
		return nil, fmt.Errorf("The YAML file was not found.")
	}

	return yamlContent, nil
}
