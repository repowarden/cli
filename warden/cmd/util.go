package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// loadWardenFile tries to intelligently choose a filepath for the Wardenfile
// and then returns a []byte with the contents of it. If customPath is not
// empty, it will try to use that before the default filenames.
func loadWardenFile(customPath string) ([]byte, error) {

	var wardenFile []byte

	possibleFilepaths := []string{"warden.yml", "warden.yaml"}

	if customPath != "" {
		possibleFilepaths = append([]string{customPath}, possibleFilepaths...)
	}

	for _, path := range possibleFilepaths {

		wardenFileTmp, err := os.ReadFile(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if errors.Is(err, fs.ErrPermission) {
			return nil, fmt.Errorf("Warden doesn't have permission to open %s", path)
		} else if err != nil {
			return nil, err
		}

		wardenFile = wardenFileTmp // this is stupid but needed to shut Go up about declared but not used

		break
	}

	if len(wardenFile) == 0 {
		return nil, fmt.Errorf("A Warden File was not found. Either './warden.yml' needs to be used or the '--wardenFile' flag set.")
	}

	return wardenFile, nil
}
