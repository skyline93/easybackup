package main

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/skyline93/easybackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   "list backupset -p /data/backup/repo1",
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
		if err := listBackupSets(listOptions.RepoPath); err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
	},
}

type ListOptions struct {
	RepoPath string
}

var listOptions ListOptions

func init() {
	cmdRoot.AddCommand(cmdList)
	cmdList.AddCommand(cmdBackupSet)

	f := cmdBackupSet.Flags()
	f.StringVarP(&listOptions.RepoPath, "repo_path", "p", "", "repo path")

	cmdBackupSet.MarkFlagRequired("repo_path")
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
