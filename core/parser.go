package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import (
	"github.com/vanilla-os/Chronos/structs"
	"gopkg.in/yaml.v3"
)

// parseArticleHeader parses the YAML header of an article and returns a
// structs.ArticleHeader object.
func parseArticleHeader(header string) (structs.ArticleHeader, error) {
	var articleHeader structs.ArticleHeader
	err := yaml.Unmarshal([]byte(header), &articleHeader)
	if err != nil {
		return structs.ArticleHeader{}, err
	}

	return articleHeader, nil
}
