package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/skyline93/easybackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   "list backupset -r repo1",
	Short: "List backup sets in repository",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdBackupSet = &cobra.Command{
	Use:   "backupset",
	Short: "backupset",
	Run: func(cmd *cobra.Command, args []string) {
		r := getRepo(listOptions.Repo)
		if r == nil {
			panic(errors.New("repo is not found"))
		}

		if err := listBackupSets(r.Path); err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
	},
}

var cmdRepos = &cobra.Command{
	Use:   "repo",
	Short: "repo",
	Run: func(cmd *cobra.Command, args []string) {
		listRepos()
	},
}

type ListOptions struct {
	Repo string
}

var listOptions ListOptions

func init() {
	cmdRoot.AddCommand(cmdList)
	cmdList.AddCommand(cmdBackupSet)
	cmdList.AddCommand(cmdRepos)

	f := cmdBackupSet.Flags()
	f.StringVarP(&listOptions.Repo, "repo", "r", "", "repo name")

	cmdBackupSet.MarkFlagRequired("repo")
}

func listBackupSets(repoPath string) error {
	repo := repository.Repository{}
	if err := repository.LoadRepository(&repo, repoPath); err != nil {
		return err
	}

	backupSets, err := repo.ListBackupSets()
	if err != nil {
		return err
	}

	items := pterm.TableData{{"BackupTime", "Id", "Type", "FromLSN", "ToLSN", "Size(Kb)"}}
	for _, bs := range backupSets {
		item := []string{bs.BackupTime, bs.Id, bs.Type, bs.FromLSN, bs.ToLSN, fmt.Sprintf("%d", bs.Size/1024)}
		items = append(items, item)
	}

	pterm.DefaultTable.WithHasHeader().WithData(items).Render()

	return nil
}

func listRepos() {
	repos := getRepos()

	items := pterm.TableData{{"Name", "Path"}}
	for _, r := range repos {
		item := []string{r.Name, r.Path}
		items = append(items, item)
	}

	pterm.DefaultTable.WithHasHeader().WithData(items).Render()
}
