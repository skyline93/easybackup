package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/skyline93/easybackup/internal/mysql"
	"github.com/skyline93/easybackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdRestore = &cobra.Command{
	Use:   "restore -r repo1 -m /usr/local/mysql/8.0.28 -t /data/restore/instance01 -i BACKUPSET_ID",
	Short: "Restore a database from backupset",
	Run: func(cmd *cobra.Command, args []string) {
		r := getRepo(restoreOptions.Repo)
		if r == nil {
			panic(errors.New("repo is not found"))
		}

		repo := repository.Repository{}
		if err := repository.LoadRepository(&repo, restoreOptions.Repo); err != nil {
			fmt.Printf("load repo error: %s", err)
			os.Exit(1)
		}

		restorer := mysql.NewRestorer()

		err := restorer.Restore(&repo, restoreOptions.TargetPath, restoreOptions.MysqlPath, restoreOptions.BackupSetId)
		if err != nil {
			fmt.Printf("restore failed, err: %s", err)
			os.Exit(1)
		}
	},
}

type RestoreOptions struct {
	BackupSetId string
	Repo        string
	TargetPath  string
	MysqlPath   string
}

var restoreOptions RestoreOptions

func init() {
	cmdRoot.AddCommand(cmdRestore)

	f := cmdRestore.Flags()
	f.StringVarP(&restoreOptions.BackupSetId, "backupset_id", "i", "", "backup set id")
	f.StringVarP(&restoreOptions.Repo, "repo", "r", "", "repo name")
	f.StringVarP(&restoreOptions.TargetPath, "target_path", "t", "", "target path")
	f.StringVarP(&restoreOptions.MysqlPath, "mysql_path", "m", "", "mysql path, example: /usr/local/mysql/8.0.28")

	cmdRestore.MarkFlagRequired("backupset_id")
	cmdRestore.MarkFlagRequired("repo")
	cmdRestore.MarkFlagRequired("target_path")
	cmdRestore.MarkFlagRequired("mysql_path")
}
