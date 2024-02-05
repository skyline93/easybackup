package file

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/skyline93/easybackup/internal/repository"
	"github.com/skyline93/easybackup/internal/snapshot"
)

type Backuper struct {
}

func NewBackuper() *Backuper {
	return &Backuper{}
}

func (b *Backuper) Backup(repo *repository.Repository, sourcePath string, backupType string) (err error) {
	bs := repository.NewBackupSet(backupType)
	targetPath, err := filepath.Abs(filepath.Join(repo.DataPath(), bs.Id))
	if err != nil {
		return err
	}

	if _, err = os.Stat(targetPath); os.IsNotExist(err) {
		err = os.MkdirAll(targetPath, 0775)
		if err != nil {
			return err
		}

		log.Printf("create path: %s", targetPath)
	}

	defer func() {
		if err != nil {
			log.Printf("backup failed, err: %s", err)
			os.RemoveAll(targetPath)
		}
	}()

	var lastSnap *snapshot.Snapshot
	if backupType == repository.TypeBackupSetIncr {
		lastBackupSet, err := repo.GetLastBackupSet(repository.TypeLog)
		if err != nil {
			return err
		}

		lastSnap, err = snapshot.LoadSnapshot(filepath.Join(repo.SnapshotPath(), lastBackupSet.Snapshot))
		if err != nil {
			return err
		}
	}

	backupTime := time.Now().Format("2006-01-02 15:04:05")
	snap := snapshot.New(repo.SnapshotPath())
	if err := snap.Snap(sourcePath, lastSnap); err != nil {
		return err
	}

	for _, node := range snap.Nodes {
		if !node.IsDir && (node.Type == snapshot.NodeNew || node.Type == snapshot.NodeModified) {
			cmd := exec.Command("cp", node.Path, targetPath)
			if err = cmd.Run(); err != nil {
				return err
			}
		}
	}

	size, err := b.getBackupSize(targetPath)
	if err != nil {
		return err
	}

	bs.Snapshot = fmt.Sprintf("%d", snap.Time)
	bs.DataType = repository.TypeLog
	bs.Size = int64(size)
	bs.BackupTime = backupTime

	if err = repo.AddBackupSet(bs); err != nil {
		return err
	}

	log.Printf("backup completed. backupset: %s\nsnapshot: %s\nsize: %dbyte", bs.Id, bs.Snapshot, bs.Size)
	return nil
}

func (b *Backuper) getBackupSize(targetPath string) (uint64, error) {
	cmd := exec.Command("du", "-sb", targetPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// 解析 du 输出获取备份数据量
	fields := strings.Fields(string(output))
	size, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}
