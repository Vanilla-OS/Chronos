package settings

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port       string            `json:"port"`
	GitRepos   []ConfigGitRepo   `json:"gitRepos"`
	LocalRepos []ConfigLocalRepo `json:"localRepos"`
}

type ConfigGitRepo struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type ConfigLocalRepo struct {
	Id   string `json:"id"`
	Path string `json:"path"`
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

	var gitRepos []ConfigGitRepo
	err = viper.UnmarshalKey("gitRepos", &gitRepos)
	if err != nil {
		panic(err)
	}

	var localRepos []ConfigLocalRepo
	err = viper.UnmarshalKey("localRepos", &localRepos)
	if err != nil {
		panic(err)
	}

	Cnf = &Config{
		Port:       viper.GetString("port"),
		GitRepos:   gitRepos,
		LocalRepos: localRepos,
	}
}
