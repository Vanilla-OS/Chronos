package structs

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

type Repo struct {
	Id              string
	Path            string
	Stories         map[string]Story
	Articles        map[string]Article
	ArticlesGrouped map[string][]Article
	Languages       []string
	RootPath        string
	FallbackLang    string
	FallbackEnabled bool
}

func (r *Repo) IsLangSupported(lang string) bool {
	for _, l := range r.Languages {
		if l == lang {
			return true
		}
	}
	return false
}
