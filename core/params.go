package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

import (
	"github.com/vanilla-os/Chronos/structs"
)

var (
	articlesDir = "articles" // Directory containing the articles

	// populated at runtime
	SupportedLang       []string                              // List of supported languages
	ArticleCache        map[string]structs.Article            // Cache of articles
	ArticleCacheGrouped map[string]map[string]structs.Article // Cache of articles grouped by language
)
