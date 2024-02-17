package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vanilla-os/Chronos/structs"
)

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Chronos is a simple, fast and lightweight documentation server written in Go.
*/

func HandleRepos(w http.ResponseWriter, r *http.Request) {
	reposBytes, err := cacheManager.Get(context.Background(), "Repos")
	if reposBytes == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var repos []structs.Repo
	err = json.Unmarshal(reposBytes, &repos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type repoResponse struct {
		Id              string   `json:"Id"`
		Count           int      `json:"Count"`
		Languages       []string `json:"Languages"`
		FallbackLang    string   `json:"FallbackLang"`
		FallbackEnabled bool     `json:"FallbackEnabled"`
	}
	response := make([]repoResponse, len(repos))
	for i, repo := range repos {
		response[i] = repoResponse{
			Id:              repo.Id,
			Count:           len(repo.Articles),
			Languages:       repo.Languages,
			FallbackLang:    repo.FallbackLang,
			FallbackEnabled: repo.FallbackEnabled,
		}
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
