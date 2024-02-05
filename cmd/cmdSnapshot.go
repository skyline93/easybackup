package main

import (
	"github.com/skyline93/easybackup/internal/snapshot"
	"github.com/spf13/cobra"
)

var cmdSnapshot = &cobra.Command{
	Use:   "snap",
	Short: "Take a snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err      error
			lastSnap *snapshot.Snapshot
		)

		if snapOptions.lastSnapPath != "" {
			lastSnap, err = snapshot.LoadSnapshot(snapOptions.lastSnapPath)
			if err != nil {
				panic(err)
			}
		}

		snap := snapshot.New("./")
		if err := snap.Snap(snapOptions.SourcePath, lastSnap); err != nil {
			panic(err)
		}
	},
}

type SnapOptions struct {
	SourcePath   string
	lastSnapPath string
	RepoPath     string
}

var snapOptions SnapOptions

func init() {
	cmdRoot.AddCommand(cmdSnapshot)

	f := cmdSnapshot.Flags()
	f.StringVarP(&snapOptions.SourcePath, "source_path", "s", "", "source path")
	f.StringVarP(&snapOptions.lastSnapPath, "last_snap_path", "l", "", "last snap path")
	f.StringVarP(&snapOptions.RepoPath, "repo_path", "p", "", "repo path")

	cmdSnapshot.MarkFlagRequired("source_path")
}
