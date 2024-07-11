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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanilla-os/Chronos/structs"
)

// HandleArticles handles requests to /articles.
func HandleArticles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId := vars["repoId"]
	lang := vars["lang"]

	if repoId == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if lang == "" || !isValidLocale(lang) {
		http.Redirect(w, r, fmt.Sprintf("/%s/articles/en", repoId), http.StatusFound)
	}

	repo, err := getRepo(repoId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !repo.IsLangSupported(lang) {
		http.Redirect(w, r, fmt.Sprintf("/%s/articles/en", repoId), http.StatusFound)
	}

	articles := repo.ArticlesGrouped[lang]
	tags := getTags(articles)
	response := structs.ArticlesResponse{
		Title:         repo.Id,
		SupportedLang: repo.Languages,
		Tags:          tags,
		Articles:      articles,
		Stories:       repo.Stories,
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func getTags(articles []structs.Article) []string {
	tags := make(map[string]bool)
	for _, article := range articles {
		for _, tag := range article.Tags {
			tags[tag] = true
		}
	}

	var tagsList []string
	for tag := range tags {
		tagsList = append(tagsList, tag)
	}

	return tagsList
}
