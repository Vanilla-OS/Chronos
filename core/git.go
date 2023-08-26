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
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

func synGitRepo(repo string, force bool) error {
	repoDir := reposDir + strings.ReplaceAll(repo, "/", "_")

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		os.Mkdir(repoDir, 0755)

		_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
			URL: repo,
		})
		if err != nil {
			return fmt.Errorf("failed to clone Git repository: %v", err)
		}
	} else {
		r, err := git.PlainOpen(repoDir)
		if err != nil {
			return fmt.Errorf("failed to open Git repository: %v", err)
		}

		remotes, err := r.Remotes()
		if err != nil {
			return fmt.Errorf("failed to find Git remote settings: %v", err)
		}

		origin := remotes[0].Config().URLs[0]

		if origin != repo {
			var confirmation bool

			if !force {
				confirmation = askForConfirmation("The Git repository has been modified. Do you want to overwrite the current one?")
			} else {
				confirmation = true
			}

			if confirmation {
				err := os.RemoveAll(repoDir)
				if err != nil {
					return fmt.Errorf("failed to remove old Git repository: %v", err)
				}

				_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
					URL: repo,
				})

				if err != nil {
					return fmt.Errorf("failed to clone Git repository: %v", err)
				}
			}
		}

		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("failed to open Git worktree: %v", err)
		}

		err = w.Pull(&git.PullOptions{})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to pull Git repository: %v", err)
		}
	}

	return nil
}

func detectGitChanges(repo string) (bool, error) {
	repoDir := reposDir + strings.ReplaceAll(repo, "/", "_")

	r, err := git.PlainOpen(repoDir)
	if err != nil {
		return false, fmt.Errorf("failed to open Git repository: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to open Git worktree: %v", err)
	}

	status, err := w.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get Git status: %v", err)
	}

	if status.IsClean() {
		return false, nil
	}

	return true, nil
}
