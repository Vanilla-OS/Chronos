package settings

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port    string `json:"port"`
	GitRepo string `json:"gitRepo"`
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

	Cnf = &Config{
		Port:    viper.GetString("port"),
		GitRepo: viper.GetString("gitRepo"),
	}
}
