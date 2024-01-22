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
	"net/http"

	"github.com/gorilla/mux"
)

func HandleLangs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId := vars["repoId"]

	repo, err := getRepo(repoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(repo.Languages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
