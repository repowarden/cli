package cmd

const (
	ERR_ACCESS_EXTRA      = "The user/team '%s' is present and shouldn't be."
	ERR_ACCESS_MISSING    = "The user/team %s is not defined."
	ERR_ACCESS_DIFFERENT  = "The user/team '%s' should have the permission '%s', not '%s'."
	ERR_ACCESS_STRATEGY   = "'%s' is not a valid access strategy."
	ERR_BRANCH_DEFAULT    = "The default branch should be '%s', not '%s'."
	ERR_LABEL_EXTRA       = "The label '%s' is present and shouldn't be."
	ERR_LABEL_MISSING     = "The label '%s' is missing."
	ERR_LICENSE_DIFFERENT = "The license should be one of '%s', not '%s'."
	ERR_LICENSE_MISSING   = "The license is missing."
	ERR_CO_DIFFERENT      = "The CODEOWNERS file is different from the policy."
	ERR_CO_MISSING        = "The CODEOWNERS file is missing."
	ERR_CO_SYNTAX         = "The CODEOWNERS file has syntax errors:\n%s"
)
