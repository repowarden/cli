package cmd

import "fmt"

const (
	ERR_ACCESS_MISSING = "The repository %s doesn't have the user."
	ERR_ACCESS_WRONG   = "The repository %s's user's permission is incorrect."
	ERR_BRANCH_DEFAULT = "The default branch should be '%s', not '%s'."
	ERR_LABEL_EXTRA    = "The label '%s' is present and shouldn't be."
	ERR_LABEL_MISSING  = "The label '%s' is missing."
	ERR_LICENSE        = "The license should be '%s', not '%s'."
	ERR_CO_DIFFERENT   = "The CODEOWNERS file is different from the policy."
	ERR_CO_MISSING     = "The CODEOWNERS file is not present."
	ERR_CO_SYNTAX      = "The CODEOWNERS file has syntax errors:\n%s"
)

// Represents a error when applying a policy, with contextual information.
type PolicyError struct {
	repository RepositoryDefinition
	message    string
	values     []any
}

// Properly print out a policy error.
func (this PolicyError) Error() string {
	return fmt.Sprintf(this.message, this.values...)
}
