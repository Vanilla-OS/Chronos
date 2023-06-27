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

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

func HandleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lang := vars["lang"]
	article := vars["article"]

	if lang == "" {
		// Redirect to the user's preferred language based on the "Accept-Language" header
		lang = getPreferredLanguage(r)
		if lang == "" || !isLanguageSupported(lang) {
			// Redirect to /en if no preferred language or unsupported language
			http.Redirect(w, r, "/articles/en/"+article, http.StatusFound)
			return
		}
	}

	if !isLanguageSupported(lang) {
		// Redirect to /en if the specified language is not supported
		http.Redirect(w, r, "/articles/en/"+article, http.StatusFound)
		return
	}

	results, ok := searchArticleCacheByLang(article, lang)
	if !ok {
		// Redirect to /en if the article doesn't exist in the specified language
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if len(results) == 0 {
		// Redirect to /en if the article doesn't exist in the specified language
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := results[0]

	parsedBody := blackfriday.Run([]byte(data.Body))
	data.Body = string(parsedBody)

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
