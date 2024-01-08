package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/skyline93/easybackup/internal/mysql"
	"github.com/skyline93/easybackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdBackup = &cobra.Command{
	Use:   "backup -r repo1 -t full",
	Short: "Take a backup",
	Run: func(cmd *cobra.Command, args []string) {
		r := getRepo(backupOptions.Repo)
		if r == nil {
			panic(errors.New("repo is not found"))
		}

		repo := repository.Repository{}
		if err := repository.LoadRepository(&repo, r.Path); err != nil {
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
	Repo       string
}

var backupOptions BackupOptions

func init() {
	cmdRoot.AddCommand(cmdBackup)

	f := cmdBackup.Flags()
	f.StringVarP(&backupOptions.BackupType, "backup_type", "t", "full", "backup type")
	f.StringVarP(&backupOptions.Repo, "repo", "r", "", "repo")

	cmdBackup.MarkFlagRequired("backup_type")
	cmdBackup.MarkFlagRequired("repo")
}
