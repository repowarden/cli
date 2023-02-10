package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
)

// What the codeowners file should look like
type codeownersPolicy struct {
	Content string   `yaml:"content"`
	Tags    []string `yaml:"tags"`
}

// Does the work to check codeowners policy against a repository
func auditCodeownersPolicy(policy codeownersPolicy, repo *wardenRepo, client *github.Client) auditResults {

	var results auditResults

	if !tagsMatched(policy.Tags, repo.Tags()) {
		return nil
	}

	file, _, _, err := client.Repositories.GetContents(context.Background(), repo.Owner, repo.Name, ".github/CODEOWNERS", nil)
	if err != nil {

		switch err.(type) {
		case *github.ErrorResponse:
			if err.(*github.ErrorResponse).Response != nil && err.(*github.ErrorResponse).Response.StatusCode == 404 {
				results.add(
					repo,
					RESULT_ERROR,
					ERR_CO_MISSING,
				)

				return results
			} else {
				fmt.Fprintf(os.Stderr, err.Error())
				return nil
			}
		default:
			fmt.Fprintf(os.Stderr, err.Error())
			return nil
		}
	}

	content, err := file.GetContent()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	// check if the files match
	if policy.Content != content {
		results.add(
			repo,
			RESULT_ERROR,
			ERR_CO_DIFFERENT,
		)
	}

	// check for codeowners syntax errors
	coErrs, _, err := client.Repositories.GetCodeownersErrors(context.Background(), repo.Owner, repo.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}

	if len(coErrs.Errors) > 0 {

		var suggestions []string
		for _, coErr := range coErrs.Errors {
			suggestions = append(suggestions, "    > "+coErr.GetSuggestion())
		}

		results.add(
			repo,
			RESULT_ERROR,
			ERR_CO_SYNTAX,
			suggestions,
		)
	}

	return results
}
