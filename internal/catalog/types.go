package catalog

import "fmt"

// MetadataFileName identifies the utility metadata file within each utility directory.
const MetadataFileName = "util.yaml"

// Utility describes a registered utility.
type Utility struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Slug        string `yaml:"slug"`
}

// Validate ensures mandatory fields are present and derives fallbacks.
func (u *Utility) Validate(defaultSlug string) error {
	if u.Name == "" {
		return fmt.Errorf("missing name")
	}

	if u.Description == "" {
		return fmt.Errorf("missing description")
	}

	if u.Slug == "" {
		u.Slug = defaultSlug
	}

	return nil
}
