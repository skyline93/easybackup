package mysql

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/skyline93/easybackup/internal/log"

	"github.com/skyline93/easybackup/internal/repository"
)

type Backuper struct {
}

func NewBackuper() *Backuper {
	return &Backuper{}
}

func (b *Backuper) Backup(repo *repository.Repository, backupType string) (err error) {
	var (
		bakSet     *repository.BackupSet
		lastBakSet *repository.BackupSet
	)

	bakSet, err = repo.CreateBackupSet(backupType)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			log.Infof("backup failed, err: %s", err)
			err = repo.DeleteBackupSet(bakSet)
		}
	}()

	if backupType == repository.TypeBackupSetIncr {
		lastBakSet, err = repo.GetLastBackupSet(repository.TypeData)
		if err != nil {
			return err
		}
	}

	backupTime := time.Now().UnixNano()

	result, err := b.runXtrabackup(*repo.Config, lastBakSet.FromLSN, bakSet.Path())
	if err != nil {
		return err
	}

	size, err := bakSet.GetSize()
	if err != nil {
		return err
	}

	if err = bakSet.Update(map[string]interface{}{
		"FromLSN":    result.FromLSN,
		"ToLSN":      result.ToLSN,
		"Size":       int64(size),
		"BackupTime": backupTime,
		"DataType":   repository.TypeData,
	}); err != nil {
		return err
	}

	log.Infof("backup completed.\nbackupset: %s", bakSet)
	return nil
}

type xtarbackupResult struct {
	FromLSN string
	ToLSN   string
}

func (b *Backuper) runXtrabackup(config repository.Config, startLsn string, targetPath string) (*xtarbackupResult, error) {
	sourceArgs := []string{
		filepath.Join(config.BinPath, "xtrabackup"),
		"--backup",
		fmt.Sprintf("--throttle=%d", config.Throttle),
		fmt.Sprintf("--login-path=%s", config.LoginPath),
		fmt.Sprintf("--datadir=%s", config.DataPath),
		"--stream=xbstream",
	}

	if config.TryCompress {
		sourceArgs = append(sourceArgs, "--compress")
	}

	if startLsn != "" {
		sourceArgs = append(sourceArgs, fmt.Sprintf("--incremental-lsn=%s", startLsn))
	}

	targetArgs := []string{
		"ssh", fmt.Sprintf("%s@%s", config.BackupUser, config.BackupHostName),
		filepath.Join(config.BinPath, "xbstream"), "-x", "-C", targetPath,
	}

	args := append(append(sourceArgs, []string{"|"}...), targetArgs...)

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", config.DbUser, config.DbHostName), strings.Join(args, " "))
	log.Infof("cmd: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filepath.Join(targetPath, "xtrabackup_checkpoints"))
	if err != nil {
		return nil, err
	}

	checkpoints, err := b.parseCheckpoints(string(content))
	if err != nil {
		return nil, err
	}

	return &xtarbackupResult{
		FromLSN: checkpoints["from_lsn"],
		ToLSN:   checkpoints["to_lsn"],
	}, nil
}

func (b *Backuper) parseCheckpoints(content string) (map[string]string, error) {
	checkpointsMap := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			checkpointsMap[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return checkpointsMap, nil
}
