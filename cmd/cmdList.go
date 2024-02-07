package main

import (
	"fmt"
	"os"
	"sort"
	"time"

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
		if err := listBackupSets(listOptions.RepoPath, listOptions.DataType); err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
	},
}

type ListOptions struct {
	RepoPath string
	DataType string
}

var listOptions ListOptions

func init() {
	cmdRoot.AddCommand(cmdList)
	cmdList.AddCommand(cmdBackupSet)

	f := cmdBackupSet.Flags()
	f.StringVarP(&listOptions.RepoPath, "repo_path", "p", "", "repo path")

	cmdBackupSet.MarkFlagRequired("repo_path")
}

type BackupSets []repository.BackupSet

func (b BackupSets) Len() int { return len(b) }

func (b BackupSets) Less(i, j int) bool {
	n, _ := time.ParseInLocation("2006-01-02 15:04:05", b[i].BackupTime, time.Local)
	m, _ := time.ParseInLocation("2006-01-02 15:04:05", b[j].BackupTime, time.Local)
	return n.Unix() > m.Unix()
}

func (b BackupSets) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func listBackupSets(repoPath string, dataType string) error {
	repo := repository.Repository{}
	if err := repository.LoadRepository(&repo, repoPath); err != nil {
		return err
	}

	var backupSets BackupSets

	for _, t := range []string{repository.TypeData, repository.TypeLog} {
		results, err := repo.ListBackupSets(t)
		if err != nil {
			return err
		}

		backupSets = append(backupSets, results...)
	}

	sort.Sort(backupSets)

	items := pterm.TableData{{"BackupTime", "Id", "DataType", "Type", "FromLSN", "ToLSN", "Size(Kb)"}}
	for _, bs := range backupSets {
		item := []string{bs.BackupTime, bs.Id, bs.DataType, bs.Type, bs.FromLSN, bs.ToLSN, fmt.Sprintf("%d", bs.Size/1024)}
		items = append(items, item)
	}

	pterm.DefaultTable.WithHasHeader().WithData(items).Render()

	return nil
}
