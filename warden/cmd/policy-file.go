package cmd

//===============================================================
// Custom types, methods, and functions needed to load and use a policy.yml file.
//===============================================================

// The top-level structure representing a policy.yml file.
type PolicyFile struct {
	DefaultBranch string         `yaml:"defaultBranch"`
	Archived      bool           `yaml:"archived"` // include archived repos in lookup?
	License       *licensePolicy `yaml:"license"`
	Labels        []string       `yaml:"labels"`
	LabelStrategy string         `yaml:"labelStrategy"`
	Access        []accessPolicy `yaml:"access"`
	CodeOwners    string         `yaml:"codeowners"`
}

// Which code licenses to allow and for which scope
type licensePolicy struct {
	Scope string   `yaml:"scope"`
	Names []string `yaml:"names"`
}
