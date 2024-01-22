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
	"context"
	"encoding/json"
	"errors"

	"github.com/vanilla-os/Chronos/structs"
)

func getRepo(repoId string) (*structs.Repo, error) {
	var repos []structs.Repo

	cRepos, err := cacheManager.Get(context.Background(), "Repos")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cRepos, &repos)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.Id == repoId {
			return &repo, nil
		}
	}

	return nil, errors.New("repo not found")
}
