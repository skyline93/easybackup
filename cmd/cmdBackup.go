package main

import (
	"fmt"
	"os"

	"github.com/skyline93/easybackup/internal/mysql"
	"github.com/skyline93/easybackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdBackup = &cobra.Command{
	Use:   "backup -p /data/backup/repo1 -t full",
	Short: "Take a backup",
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.Repository{}
		if err := repository.LoadRepository(&repo, backupOptions.RepoPath); err != nil {
			fmt.Printf("load repo error: %s", err)
			os.Exit(1)
		}

		backuper := mysql.NewBackuper()
		if err := backuper.Backup(&repo, backupOptions.BackupType); err != nil {
			fmt.Printf("backup failed error: %s", err)
			os.Exit(1)
		}
	},
}

type BackupOptions struct {
	BackupType string
	RepoPath   string
}

var backupOptions BackupOptions

func init() {
	cmdRoot.AddCommand(cmdBackup)

	f := cmdBackup.Flags()
	f.StringVarP(&backupOptions.BackupType, "backup_type", "t", "full", "backup type")
	f.StringVarP(&backupOptions.RepoPath, "repo_path", "p", "", "repo path")

	cmdBackup.MarkFlagRequired("backup_type")
	cmdBackup.MarkFlagRequired("repo_path")
}
