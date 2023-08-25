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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleSearch handles requests to /search.
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId := vars["repoId"]
	lang := vars["lang"]

	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if repoId == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if lang == "" || !isValidLocale(lang) {
		fmt.Println("Lang not found, redirecting to en")
		http.Redirect(w, r, "/articles/"+repoId+"en/", http.StatusFound)
	}

	repo, err := getRepo(repoId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	results := searchArticles(repo.Id, lang, query)
	jsonData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)

}
