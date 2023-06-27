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

// HandleArticles handles requests to /articles.
func HandleArticles(w http.ResponseWriter, r *http.Request) {
	response := structs.ArticlesResponse{
		Title:           "Chronos",
		SupportedLang:   SupportedLang,
		GroupedArticles: ArticleCacheGrouped,
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
