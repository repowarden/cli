package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

// loadPolicyFile tries to intelligently choose a filepath for the
// policy file and then return the unmarshalled struct. If customPath is not
// empty, it will try to use that before the default filenames.
func loadPolicyFile(customPath string) (*PolicyFile, []byte, error) {

	var file PolicyFile
	var foundFile bool
	var yamlContent []byte
	var err error

	possibleFilepaths := []string{"policy."}

	if customPath != "" {
		possibleFilepaths = append([]string{customPath}, possibleFilepaths...)
	}

	for _, path := range possibleFilepaths {

		yamlContent, err = loadYAMLFile(path)
		if err == nil {
			foundFile = true
			break
		}
	}

	err = yaml.Unmarshal(yamlContent, &file)
	if err != nil {
		return nil, nil, err
	}

	if !foundFile {
		return nil, nil, fmt.Errorf("A policy file was not found. Either './policy.yml' needs to be used or the '--policyFile' flag set.")
	}

	return &file, yamlContent, nil
}

// loadRepositoriesFile tries to intelligently choose a filepath for the
// Wardenfile and then return the unmarshalled struct. If customPath is not
// empty, it will try to use that before the default filenames.
func loadRepositoriesFile(customPath string) ([]RepositoryDefinition, []byte, error) {

	var repositoriesFile []RepositoryDefinition
	var yamlContent []byte
	var foundFile bool
	var err error

	possibleFilepaths := []string{"repositories."}

	if customPath != "" {
		possibleFilepaths = append([]string{customPath}, possibleFilepaths...)
	}

	for _, path := range possibleFilepaths {

		yamlContent, err = loadYAMLFile(path)
		if err == nil {
			foundFile = true
			break
		}
	}

	err = yaml.Unmarshal(yamlContent, &repositoriesFile)
	if err != nil {
		return nil, nil, err
	}

	if !foundFile {
		return nil, nil, fmt.Errorf("A Repositories file was not found. Either './repositories.yml' needs to be used or the '--repositoriesFile' flag set.")
	}

	return repositoriesFile, yamlContent, nil
}

// loadYAMLFile loads and unmarshals a YAML file. Both .yml and .yaml will be
// attempted and in that order, so end the filename with a period. For example:
// `my-file.` or `repositories.`.
func loadYAMLFile(filepath string) ([]byte, error) {

	var yamlContent []byte
	var err error

	possibleFiles := []string{filepath + "yml", filepath + "yaml"}

	for _, path := range possibleFiles {

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
