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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vanilla-os/Chronos/settings"
	"github.com/vanilla-os/Chronos/structs"
)

var (
	reposDir = "repos/"
)

func LoadChronos() {
	if _, err := os.Stat(reposDir); os.IsNotExist(err) {
		os.MkdirAll(reposDir, 0755)
	}

	prepareCache()

	err := prepareRepos()
	if err != nil {
		panic(err)
	}
}

// prepareRepos prepares both local and Git repositories.
func prepareRepos() error {
	var repos []structs.Repo

	cRepos, err := cacheManager.Get(context.Background(), "Repos")
	if err == nil {
		fmt.Println("Repos cache found")
		err = json.Unmarshal(cRepos, &repos)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Repos cache not found, resetting")
	}

	for _, repo := range settings.Cnf.GitRepos {
		fmt.Printf("Synchronizing Git repository: %s\n", repo.Url)
		err := synGitRepo(repo.Url)
		if err != nil {
			return fmt.Errorf("failed to synchronize Git repository: %v", err)
		}

		_repo := structs.Repo{
			Id:   repo.Id,
			Path: reposDir + strings.ReplaceAll(repo.Url, "/", "_"),
		}

		_repo.Languages, err = getRepoLanguages(_repo)
		if err != nil {
			return err
		}

		_repo.Articles, err = getRepoArticles(_repo)
		if err != nil {
			return err
		}

		repos = append(repos, _repo)
	}

	for _, repo := range settings.Cnf.LocalRepos {
		_repo := structs.Repo{
			Id:   repo.Id,
			Path: repo.Path,
		}

		_repo.Languages, err = getRepoLanguages(_repo)
		if err != nil {
			return err
		}

		_repo.Articles, err = getRepoArticles(_repo)
		if err != nil {
			return err
		}

		repos = append(repos, _repo)
	}

	reposBytes, err := json.Marshal(repos)
	if err != nil {
		fmt.Printf("Failed to marshal repos: %v\n", err)
	}

	cacheManager.Set(context.Background(), "Repos", reposBytes)

	return nil
}

// getRepoLanguages populates the language cache.
func getRepoLanguages(repo structs.Repo) ([]string, error) {
	langs, err := loadLanguagesFromRepo(repo)
	if err != nil {
		return nil, err
	}

	tmpLangCache := make([]string, 0)
	tmpLangCache = append(tmpLangCache, langs...)

	return tmpLangCache, nil
}

// getRepoArticles populates the article cache.
func getRepoArticles(repo structs.Repo) (map[string]structs.Article, error) {
	articlePaths, err := loadArticlesFromRepo(repo)
	if err != nil {
		return nil, err
	}

	tmpArticleCache := make(map[string]structs.Article)
	tmpArticleCacheGrouped := make(map[string]map[string]structs.Article)
	for _, articlePath := range articlePaths {
		article, err := loadArticle(articlePath)
		if err != nil {
			return nil, err
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

	return tmpArticleCache, nil
}

// loadLanguagesFromRepo returns a list of languages from the repo folder.
func loadLanguagesFromRepo(repo structs.Repo) ([]string, error) {
	langs := make([]string, 0)
	dirEntries, err := os.ReadDir(filepath.Join(repo.Path, "articles"))
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() || !isValidLocale(entry.Name()) {
			continue
		}

		langs = append(langs, entry.Name())
	}

	if len(langs) == 0 {
		return nil, errors.New("no languages found")
	}

	return langs, nil
}

// loadArticlesFromRepo returns a list of articles from the repo folder.
func loadArticlesFromRepo(repo structs.Repo) ([]string, error) {
	articles := make([]string, 0)

	for _, lang := range repo.Languages {
		langDir := filepath.Join(repo.Path, "articles", lang)
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
	slug := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	lang := strings.Split(path, string(filepath.Separator))[3]

	article := structs.Article{
		Title:           title,
		Description:     description,
		PublicationDate: publicationDate,
		Authors:         authors,
		Body:            body,
		Path:            path,
		Url:             strings.TrimSuffix(path, filepath.Ext(path)),
		Slug:            slug,
		Language:        lang,
	}

	return article, nil
}
