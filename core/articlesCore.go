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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/vanilla-os/Chronos/settings"
	"github.com/vanilla-os/Chronos/structs"
)

// loadArticle loads an article from the specified path.
func loadArticle(path string) (structs.Article, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return structs.Article{}, err
	}

	parts := strings.SplitN(string(content), "\n\n", 2)
	if len(parts) != 2 {
		return structs.Article{}, fmt.Errorf("invalid article format: %s", path)
	}

	header := parts[0]
	body := parts[1]

	title, description, publicationDate, authors := parseArticleHeader(header)

	article := structs.Article{
		Title:           title,
		Description:     description,
		PublicationDate: publicationDate,
		Authors:         authors,
		Body:            body,
		Path:            path,
		Url:             strings.TrimSuffix(path, filepath.Ext(path)),
	}

	return article, nil
}

// parseArticleHeader parses the header of an article and extracts the information.
func parseArticleHeader(header string) (string, string, string, []string) {
	title := ""
	description := ""
	publicationDate := ""
	authors := make([]string, 0)

	lines := strings.Split(header, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Title":
			title = value
		case "Description":
			description = value
		case "PublicationDate":
			publicationDate = value
		case "Authors":
			authors = strings.Split(value, ",")
		}
	}

	return title, description, publicationDate, authors
}

// PopulateArticleCache populates the article cache.
func PopulateArticleCache() error {
	articlePaths, err := getArticleList()
	if err != nil {
		return err
	}

	tmpArticleCache := make(map[string]structs.Article)
	tmpArticleCacheGrouped := make(map[string]map[string]structs.Article)
	for _, articlePath := range articlePaths {
		article, err := loadArticle(articlePath)
		if err != nil {
			return err
		}

		tmpArticleCache[articlePath] = article

		lang := strings.Split(articlePath, string(filepath.Separator))[1]
		if lang == "articles" {
			lang = strings.Split(articlePath, string(filepath.Separator))[2]
		}
		if _, ok := tmpArticleCacheGrouped[lang]; !ok {
			tmpArticleCacheGrouped[lang] = make(map[string]structs.Article)
		}

		tmpArticleCacheGrouped[lang][articlePath] = article
	}

	ArticleCache = tmpArticleCache
	ArticleCacheGrouped = tmpArticleCacheGrouped

	return nil
}

// getArticleList retrieves the list of available articles.
func getArticleList() ([]string, error) {
	articles := make([]string, 0)

	if settings.Cnf.GitRepo != "" {
		repoDir := "articles_repo"

		if _, err := os.Stat(repoDir); os.IsNotExist(err) {
			os.Mkdir(repoDir, 0755)

			_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
				URL: settings.Cnf.GitRepo,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to clone Git repository: %v", err)
			}
		} else {
			r, err := git.PlainOpen(repoDir)
			if err != nil {
				return nil, fmt.Errorf("failed to open Git repository: %v", err)
			}

			remotes, err := r.Remotes()
			if err != nil {
        return nil, fmt.Errorf("failed to find Git remote settings: %v", err)
			}

      origin := remotes[0].Config().URLs[0]
      
      if origin != settings.Cnf.GitRepo {
      //not the same origin
      //overwrite 
        return nil, nil
      }

			w, err := r.Worktree()
			if err != nil {
				return nil, fmt.Errorf("failed to open Git worktree: %v", err)
			}

			err = w.Pull(&git.PullOptions{})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				return nil, fmt.Errorf("failed to pull Git repository: %v", err)
			}
		}

		articlesDir = filepath.Join(repoDir, "articles")
	}

	for _, lang := range SupportedLang {
		langDir := filepath.Join(articlesDir, lang)
		dirEntries, err := os.ReadDir(langDir)
		if err != nil {
			return nil, err
		}

		for _, entry := range dirEntries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
				continue
			}

			articlePath := filepath.Join(langDir, entry.Name())
			articles = append(articles, articlePath)
		}
	}

	if len(articles) == 0 {
		return nil, errors.New("no articles found")
	}

	return articles, nil
}
