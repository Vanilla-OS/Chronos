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
	"net/http"

	"github.com/vanilla-os/Chronos/structs"
)

// HandleSearch handles requests to /search.
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	lang := getPreferredLanguage(r)

	res, ok := searchArticleCacheByLang(query, lang)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	response := structs.ResultsResponse{
		Query:   query,
		Results: res,
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
