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
	"strings"

	"github.com/vanilla-os/Chronos/structs"
)

func searchArticles(repoId string, lang string, query string) []structs.Article {
	var results []structs.Article
	repo, err := getRepo(repoId)
	if err != nil {
		return results
	}

	for _, article := range repo.Articles {
		if article.Language == lang {
			results = append(results, article)
		}
	}

	return filterByMatch(query, results)
}

func searchArticle(repoId string, lang string, query string) (structs.Article, bool) {
	articles := searchArticles(repoId, lang, query)
	if len(articles) > 0 {
		return articles[0], true
	}
	return structs.Article{}, false
}

func filterByMatch(query string, articles []structs.Article) []structs.Article {
	var exactMatch []structs.Article
	var partialMatch []structs.Article

	for _, article := range articles {
		if article.Slug == query || article.Title == query {
			exactMatch = append(exactMatch, article)
		} else if strings.Contains(article.Title, query) || strings.Contains(article.Slug, query) {
			partialMatch = append(partialMatch, article)
		}
	}

	orderedResults := append(exactMatch, partialMatch...)
	return orderedResults
}
