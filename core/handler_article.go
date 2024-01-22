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
	"github.com/russross/blackfriday/v2"
)

func HandleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId := vars["repoId"]
	lang := vars["lang"]
	slug := vars["slug"]

	if repoId == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if lang == "" || !isValidLocale(lang) {
		http.Redirect(w, r, fmt.Sprintf("/%s/articles/en/%s", repoId, slug), http.StatusFound)
	}

	result, ok := searchArticle(repoId, lang, slug)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	parsedBody := blackfriday.Run([]byte(result.Body))
	result.Body = string(parsedBody)

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
