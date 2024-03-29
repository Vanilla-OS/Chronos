package settings

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port                  string       `json:"port"`
	GitRepos              []ConfigRepo `json:"gitRepos"`
	LocalRepos            []ConfigRepo `json:"localRepos"`
	BackgroundCacheUpdate bool         `json:"backgroundCacheUpdate"`
	CacheBackend          string       `json:"cacheBackend"`

	// Redis specific settings
	RedisCacheServer   string `json:"redisCacheServer"`
	RedisCachePort     string `json:"redisCachePort"`
	RedisCacheUsername string `json:"redisCacheUsername"`
	RedisCachePassword string `json:"redisCachePassword"`
	RedisCacheDB       int    `json:"redisCacheDB"`
}

type ConfigRepo struct {
	Id           string `json:"id"`
	Url          string `json:"url"`
	RootPath     string `json:"rootPath"`
	FallbackLang string `json:"fallbackLang"`
}

var Cnf *Config

func init() {
	viper.SetDefault("port", "8080")
	viper.SetDefault("gitRepo", "")

	// prod paths
	viper.AddConfigPath("/etc/chronos/")
	viper.AddConfigPath("/usr/share/chronos/")

	// dev paths
	viper.AddConfigPath("config/")

	// tests paths
	viper.AddConfigPath("../config/")

	viper.SetConfigName("chronos")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var gitRepos []ConfigRepo
	err = viper.UnmarshalKey("gitRepos", &gitRepos)
	if err != nil {
		panic(err)
	}

	var localRepos []ConfigRepo
	err = viper.UnmarshalKey("localRepos", &localRepos)
	if err != nil {
		panic(err)
	}

	Cnf = &Config{
		Port:                  viper.GetString("port"),
		GitRepos:              gitRepos,
		LocalRepos:            localRepos,
		BackgroundCacheUpdate: viper.GetBool("backgroundCacheUpdate"),
		CacheBackend:          viper.GetString("cacheBackend"),

		RedisCacheServer:   viper.GetString("redisCacheServer"),
		RedisCachePort:     viper.GetString("redisCachePort"),
		RedisCacheUsername: viper.GetString("redisCacheUsername"),
		RedisCachePassword: viper.GetString("redisCachePassword"),
		RedisCacheDB:       viper.GetInt("redisCacheDB"),
	}
}
