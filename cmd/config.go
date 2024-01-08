package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

type Repo struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

var (
	repos      []Repo
	configPath string
)

func initConfig() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	userDir := currentUser.HomeDir

	viper.AddConfigPath(userDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	configPath = filepath.Join(userDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		viper.SetDefault("repos", []Repo{})

		if err = viper.WriteConfigAs(configPath); err != nil {
			panic(err)
		}
	}

	if err = viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.UnmarshalKey("repos", &repos)
	if err != nil {
		panic(err)
	}
}

func addRepo(repo *Repo) error {
	for _, r := range repos {
		if r.Name == repo.Name {
			return errors.New("repo is exist")
		}
	}

	repos = append(repos, *repo)
	viper.Set("repos", repos)
	if err := viper.WriteConfigAs(configPath); err != nil {
		return err
	}

	return nil
}

func getRepo(repoName string) *Repo {
	for _, r := range repos {
		if r.Name == repoName {
			return &r
		}
	}

	return nil
}

func getRepos() []Repo {
	return repos
}
