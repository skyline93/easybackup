package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/skyline93/easybackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:   "init -f config.json -p /data/backup -n repo1",
	Short: "Init a repository",
	Run: func(cmd *cobra.Command, args []string) {
		configData, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		var config repository.Config
		if err = json.Unmarshal(configData, &config); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		r := repository.NewRepository(repoName, &config)
		if err := r.Init(rootPath); err != nil {
			fmt.Printf("err: %s\n", err)
			os.Exit(1)
		}
	},
}

var (
	configFile string
	repoName   string
	rootPath   string
)

func init() {
	cmdRoot.AddCommand(cmdInit)

	cmdInit.Flags().StringVarP(&configFile, "config", "f", "config.yaml", "config file path")
	cmdInit.Flags().StringVarP(&repoName, "repo_name", "n", "", "repository name")
	cmdInit.Flags().StringVarP(&rootPath, "path", "p", "", "repository path")

	cmdInit.MarkFlagRequired("config")
	cmdInit.MarkFlagRequired("repo_name")
	cmdInit.MarkFlagRequired("path")
}
