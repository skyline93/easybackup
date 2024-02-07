package file

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/skyline93/easybackup/internal/log"

	"github.com/skyline93/easybackup/internal/repository"
)

type Restorer struct {
}

func NewRestorer() *Restorer {
	return &Restorer{}
}

func (r *Restorer) Restore(repo *repository.Repository, targetPath string, backupSetId string) (err error) {
	if _, err = os.Stat(targetPath); os.IsNotExist(err) {
		if err = os.MkdirAll(targetPath, 0775); err != nil {
			return
		}
	}

	defer func() {
		if err != nil {
			cmd := exec.Command("sudo", "rm", "-rf", targetPath)
			log.Infof("run cmd: %s", cmd)
			if err = cmd.Run(); err != nil {
				return
			}
		}
	}()

	var backupSets []repository.BackupSet

	backupSet, err := repo.GetBackupSet(repository.TypeData, backupSetId)
	if err != nil {
		return err
	}

	if backupSet.Type == repository.TypeBackupSetFull {
		backupSets = append(backupSets, *backupSet)
	} else if backupSet.Type == repository.TypeBackupSetIncr {
		backupSets, err = repo.GetBeforeBackupSet(repository.TypeData, backupSetId)
		if err != nil {
			return err
		}
	}

	var dataPath string
	for _, bs := range backupSets {
		if bs.Type == repository.TypeBackupSetFull {
			dataPath = targetPath
		}

		files, err := filepath.Glob(fmt.Sprintf("%s/%s", filepath.Join(repo.DataPath(), bs.Id), "*"))
		if err != nil {
			return err
		}

		for _, file := range files {
			cmd := exec.Command("cp", "-r", file, dataPath)
			log.Infof("run cmd: %s", cmd)
			if err = cmd.Run(); err != nil {
				return err
			}
		}
	}
	return nil
}
