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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vanilla-os/Chronos/settings"
	"github.com/vanilla-os/Chronos/structs"
	"gopkg.in/yaml.v3"
)

var (
	reposDir = "repos/"
)

// LoadChronos loads the Chronos server, preparing the cache and start the
// background cache update if configured.
func LoadChronos() (err error) {
	_, err = os.Stat(reposDir)
	if os.IsNotExist(err) {
		os.MkdirAll(reposDir, 0755)
	}

	err = InitCacheManager()
	if err != nil {
		return err
	}

	err = prepareRepos(true)
	if err != nil {
		return err
	}

	if settings.Cnf.BackgroundCacheUpdate {
		var wg sync.WaitGroup
		wg.Add(1)
		go backgroundCacheUpdate(15*time.Minute, &wg)
		wg.Wait()
	}

	return nil
}

// backgroundCacheUpdate updates the cache in the background.
func backgroundCacheUpdate(interval time.Duration, wg *sync.WaitGroup) {
	for {
		log.Println("(loader): Starting background cache update...")

		for _, repo := range settings.Cnf.GitRepos {
			changed, err := detectGitChanges(repo.Url)
			if err != nil {
				log.Printf("(loader): Failed to detect Git changes: %v\n", err)
			}

			if changed {
				err := synGitRepo(repo.Url, true)
				if err != nil {
					log.Printf("(loader): Failed to synchronize Git repository: %v\n", err)
				}
			}
		}

		err := prepareRepos(false)
		if err != nil {
			log.Printf("(loader): Failed to prepare repos: %v\n", err)
			continue
		}

		log.Println("(loader): Finished background cache update")

		if wg != nil {
			wg.Done()
			wg = nil
		}

		log.Printf("(loader): Sleeping for %s\n", interval)
		time.Sleep(interval)

		log.Println("(loader): Waking up for next background cache update")
	}
}

// prepareRepos prepares both local and Git repositories.
func prepareRepos(needSyncGit bool) error {
	var repos []structs.Repo
	var err error

	log.Println("(loader): Preparing Git repositories cache")

	for _, repo := range settings.Cnf.GitRepos {
		rootPath := "articles"
		if repo.RootPath != "" {
			rootPath = repo.RootPath
		}

		if needSyncGit {
			log.Printf("(loader): Synchronizing Git repository: %s\n", repo.Url)
			err := synGitRepo(repo.Url, false)
			if err != nil {
				return fmt.Errorf("failed to synchronize Git repository: %v", err)
			}
		}

		_repo := structs.Repo{
			Id:           repo.Id,
			Path:         reposDir + strings.ReplaceAll(repo.Url, "/", "_"),
			RootPath:     rootPath,
			FallbackLang: repo.FallbackLang,
		}

		log.Printf("(loader): Loading languages for Git repository: %s\n", repo.Url)
		_repo.Languages, err = getRepoLanguages(&_repo)
		if err != nil {
			return err
		}

		log.Printf("(loader): Loading stories for local repository: %s\n", repo.Url)
		_repo.Stories, err = loadStories(&_repo)
		if err != nil {
			return err
		}

		log.Printf("(loader): Loading articles for Git repository: %s\n", repo.Url)
		_repo.Articles, err = getRepoArticles(_repo)
		if err != nil {
			return err
		}

		log.Printf("(loader): Grouping articles for Git repository: %s\n", repo.Url)
		_repo.ArticlesGrouped, err = groupArticles(_repo)
		if err != nil {
			return err
		}

		repos = append(repos, _repo)
	}

	log.Println("(loader): Preparing local repositories cache")

	for _, repo := range settings.Cnf.LocalRepos {
		rootPath := "articles"
		if repo.Url != "" {
			rootPath = repo.RootPath
		}

		_repo := structs.Repo{
			Id:           repo.Id,
			Path:         repo.Url,
			RootPath:     rootPath,
			FallbackLang: repo.FallbackLang,
		}

		log.Printf("(loader): Loading languages for local repository: %s\n", repo.Url)
		_repo.Languages, err = getRepoLanguages(&_repo)
		if err != nil {
			return err
		}

		log.Printf("(loader): Loading stories for local repository: %s\n", repo.Url)
		_repo.Stories, err = loadStories(&_repo)
		if err != nil {
			return err
		}

		log.Printf("(loader): Loading articles for local repository: %s\n", repo.Url)
		_repo.Articles, err = getRepoArticles(_repo)
		if err != nil {
			return err
		}

		log.Printf("(loader): Grouping articles for local repository: %s\n", repo.Url)
		_repo.ArticlesGrouped, err = groupArticles(_repo)
		if err != nil {
			return err
		}

		repos = append(repos, _repo)
	}

	reposBytes, err := json.Marshal(repos)
	if err != nil {
		log.Printf("(loader): Failed to marshal repos: %v\n", err)
	}

	cacheManager.Set(context.Background(), "Repos", reposBytes)

	log.Printf("(loader): Finished preparing repositories cache: %d repos\n", len(repos))

	return nil
}

// getRepoLanguages populates the language cache.
func getRepoLanguages(repo *structs.Repo) ([]string, error) {
	langs, err := loadLanguagesFromRepo(repo)
	if err != nil {
		return nil, err
	}

	tmpLangCache := make([]string, len(langs))
	copy(tmpLangCache, langs)

	return tmpLangCache, nil
}

// getRepoArticles populates the article cache.
func getRepoArticles(repo structs.Repo) (map[string]structs.Article, error) {
	articlePaths, err := loadArticlesFromRepo(repo)
	if err != nil {
		return nil, err
	}

	tmpArticleCache := make(map[string]structs.Article, len(articlePaths))
	tmpArticleCacheGrouped := make(map[string]map[string]structs.Article)
	for _, articlePath := range articlePaths {
		article, err := loadArticle(repo, articlePath)
		if err != nil {
			return nil, err
		}

		tmpArticleCache[articlePath] = article

		var lang string
		if repo.FallbackEnabled {
			lang = repo.FallbackLang
			if lang == "" {
				lang = "en"
			}
		} else {
			furtherPath := strings.TrimPrefix(articlePath, filepath.Join(repo.Path, repo.RootPath))
			lang = strings.Split(furtherPath, string(filepath.Separator))[1]
		}

		if _, ok := tmpArticleCacheGrouped[lang]; !ok {
			tmpArticleCacheGrouped[lang] = make(map[string]structs.Article)
		}

		tmpArticleCacheGrouped[lang][articlePath] = article
	}

	return tmpArticleCache, nil
}

// groupArticles groups articles by language.
func groupArticles(repo structs.Repo) (map[string][]structs.Article, error) {
	tmpArticleCacheGrouped := make(map[string][]structs.Article)
	for _, article := range repo.Articles {
		tmpArticleCacheGrouped[article.Language] = append(tmpArticleCacheGrouped[article.Language], article)
	}

	return tmpArticleCacheGrouped, nil
}

// loadLanguagesFromRepo returns a list of languages from the repo folder.
func loadLanguagesFromRepo(repo *structs.Repo) ([]string, error) {
	langs := make([]string, 0)
	dirEntries, err := os.ReadDir(filepath.Join(repo.Path, repo.RootPath))
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
		if repo.FallbackLang == "" {
			log.Printf("(loader): No languages found for repo: %s, assuming English\n", repo.Path)
			langs = append(langs, "en")
			repo.FallbackEnabled = true
		}
	}

	return langs, nil
}

// loadArticlesFromRepo returns a list of articles from the repo folder.
func loadArticlesFromRepo(repo structs.Repo) ([]string, error) {
	articles := make([]string, 0)

	for _, lang := range repo.Languages {
		langDir := filepath.Join(repo.Path, repo.RootPath, lang)
		// if langDir does not exist, assuming we are in a fallback language
		// situation and we should use the root path instead
		if _, err := os.Stat(langDir); os.IsNotExist(err) {
			langDir = filepath.Join(repo.Path, repo.RootPath)
		}
		// if still langDir does not exist, we have a problem, well, the user
		// has a problem, we just panic
		if _, err := os.Stat(langDir); os.IsNotExist(err) {
			return nil, errors.New("no articles found")
		}

		dirEntries, err := os.ReadDir(langDir)
		if err != nil {
			return nil, err
		}

		for _, entry := range dirEntries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
				continue
			}

			articles = append(articles, filepath.Join(langDir, entry.Name()))
		}
	}

	if len(articles) == 0 {
		return nil, errors.New("no articles found")
	}

	return articles, nil
}

// loadArticle loads an article from the specified path.
func loadArticle(repo structs.Repo, path string) (structs.Article, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return structs.Article{}, err
	}

	parts := strings.SplitN(string(content), "\n\n", 2)
	if len(parts) != 2 {
		return structs.Article{}, fmt.Errorf("invalid article format: %s", path)
	}

	rawHeader := parts[0]
	body := parts[1]

	header, err := parseArticleHeader(rawHeader)
	if err != nil {
		return structs.Article{}, err
	}

	slug := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	lang := repo.FallbackLang
	if repo.FallbackEnabled && lang == "" {
		lang = "en"
	} else {
		furtherPath := strings.TrimPrefix(path, filepath.Join(repo.Path, repo.RootPath))
		lang = strings.Split(furtherPath, string(filepath.Separator))[1]
	}

	story, err := loadStory(repo, header.StoryId)
	if err != nil {
		return structs.Article{}, fmt.Errorf("failed to load story: %v", err)
	}

	article := structs.Article{
		StoryId:         header.StoryId,
		Story:           story,
		Previous:        header.Previous,
		Next:            header.Next,
		Listed:          header.Listed,
		Title:           header.Title,
		Description:     header.Description,
		PublicationDate: header.PublicationDate,
		Authors:         header.Authors,
		Tags:            header.Tags,
		Body:            body,
		Path:            path,
		Url:             strings.TrimSuffix(path, filepath.Ext(path)),
		Slug:            slug,
		Language:        lang,
	}

	return article, nil
}

// loadStories loads all the stories from the stories.yml file in the repository.
func loadStories(repo *structs.Repo) (map[string]structs.Story, error) {
	storiesPath := filepath.Join(repo.Path, repo.RootPath, "stories.yml")
	storiesFile, err := os.ReadFile(storiesPath)
	if err != nil {
		fmt.Printf("(loader): No stories file found for repo: %s\n", repo.Path)
		return nil, nil // safe to ignore, stories file is optional
	}

	var stories []structs.Story
	err = yaml.Unmarshal(storiesFile, &stories)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stories file: %v", err)
	}

	storiesMap := make(map[string]structs.Story)
	for _, story := range stories {
		storiesMap[story.Id] = story
	}

	return storiesMap, nil
}

// loadStory loads a story from the repository's stories map using its ID.
func loadStory(repo structs.Repo, storyId string) (*structs.Story, error) {
	if storyId == "" {
		return &structs.Story{}, nil // safe to ignore, stories are optional
	}

	if len(repo.Stories) == 0 || repo.Stories == nil {
		return &structs.Story{}, fmt.Errorf("no stories found but requested story with ID %s", storyId)
	}

	story, exists := repo.Stories[storyId]
	if !exists {
		return &structs.Story{}, fmt.Errorf("story with ID %s not found", storyId)
	}

	return &story, nil
}
