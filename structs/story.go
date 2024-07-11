package structs

// Story is the struct that represents a story, a sequence of articles.
type Story struct {
	Id          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	StartSlug   string `yaml:"startSlug"`
}
