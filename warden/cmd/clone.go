package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cloneCmd = &cobra.Command{
		Use:   "clone",
		Short: "Clones a git repo to the ~/Repos/org/name structure",
		RunE: func(cmd *cobra.Command, args []string) error {

			repo := args[0]
			urlPieces := strings.Split(repo, "/")
			org := urlPieces[3]
			name := urlPieces[4]

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			// create directories if needed
			filepath := homeDir + "/Repos/" + org
			aCmd := exec.Command("mkdir", "-p", filepath)
			stdout, err := aCmd.Output()
			if err != nil {
				fmt.Printf("%s/n", stdout)
				return err
			}

			aCmd = exec.Command("git", "clone", repo, filepath+"/"+name)
			stdout, err = aCmd.Output()
			if err != nil {
				fmt.Printf("%s/n", stdout)
				return err
			}

			fmt.Printf("%s/n", stdout)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(cloneCmd)
}
