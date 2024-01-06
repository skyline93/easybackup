package mysql

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/skyline93/mysql-xtrabackup/internal/repository"
	"gopkg.in/ini.v1"
)

type Restorer struct {
}

func NewRestorer() *Restorer {
	return &Restorer{}
}

func (r *Restorer) Restore(repo *repository.Repository, targetPath string, mysqlPath string, backupSetId string) (err error) {
	if _, err = os.Stat(targetPath); os.IsNotExist(err) {
		if err = os.MkdirAll(targetPath, 0775); err != nil {
			return
		}
	}

	defer func() {
		if err != nil {
			cmd := exec.Command("sudo", "rm", "-rf", targetPath)
			log.Printf("run cmd: %s", cmd)
			if err = cmd.Run(); err != nil {
				return
			}
		}
	}()

	var backupSets []repository.BackupSet

	backupSet, err := repo.GetBackupSet(backupSetId)
	if err != nil {
		return err
	}

	if backupSet.Type == repository.TypeBackupSetFull {
		backupSets = append(backupSets, *backupSet)
	} else if backupSet.Type == repository.TypeBackupSetIncr {
		backupSets, err = repo.GetBeforeBackupSet(backupSetId)
		if err != nil {
			return err
		}
	}

	var dataPath string
	for _, bs := range backupSets {
		targetSubPath := filepath.Join(targetPath, bs.Id)

		if bs.Type == repository.TypeBackupSetFull {
			dataPath = targetSubPath
		}

		// 拷贝文件
		cmd := exec.Command("cp", "-r", filepath.Join(repo.DataPath(), bs.Id), targetSubPath)
		log.Printf("run cmd: %s", cmd)
		if err = cmd.Run(); err != nil {
			return err
		}

		// 解压
		cmd = exec.Command(
			filepath.Join(repo.Config.BinPath, "xtrabackup"),
			"--decompress", "--remove-original",
			fmt.Sprintf("--target-dir=%s", targetSubPath),
		)
		log.Printf("run cmd: %s", cmd)
		if err = cmd.Run(); err != nil {
			return err
		}

		// 追日志
		args := []string{
			"--prepare", "--apply-log-only",
			fmt.Sprintf("--target-dir=%s", dataPath),
		}

		if bs.Type == repository.TypeBackupSetIncr {
			args = append(args, fmt.Sprintf("--incremental-dir=%s", targetSubPath))
		}

		cmd = exec.Command(filepath.Join(repo.Config.BinPath, "xtrabackup"), args...)
		log.Printf("run cmd: %s", cmd)
		if err = cmd.Run(); err != nil {
			return err
		}

		// 删除增量目标文件
		if bs.Type == repository.TypeBackupSetIncr {
			if err = os.RemoveAll(targetSubPath); err != nil {
				return err
			}
		}
	}

	bkMyCnf, err := ini.Load(filepath.Join(dataPath, "backup-my.cnf"))
	if err != nil {
		return err
	}

	sec := bkMyCnf.Section("mysqld")

	requiredKeys := []string{
		"innodb_data_file_path",
		"innodb_log_files_in_group",
		"innodb_log_file_size",
		"innodb_page_size",
		"innodb_undo_directory",
		"innodb_undo_tablespaces",
		"server_id",
		"lower-case-table-names",
		"log-bin",
	}

	configs := make(map[string]string)
	for _, k := range requiredKeys {
		v, err := sec.GetKey(k)
		if err != nil {
			continue
		}

		configs[k] = v.String()
	}

	freePort, err := GetFreePort()
	if err != nil {
		return err
	}

	configs["basedir"] = mysqlPath
	configs["datadir"] = dataPath
	configs["socket"] = filepath.Join(targetPath, "mysql.sock")
	configs["pid-file"] = filepath.Join(targetPath, "mysql.pid")
	configs["log-error"] = filepath.Join(targetPath, "mysql.err")
	configs["port"] = fmt.Sprintf("%d", freePort)

	myCnf := ini.Empty(ini.LoadOptions{AllowBooleanKeys: true})
	sec, err = myCnf.NewSection("mysqld")
	if err != nil {
		return err
	}

	for k, v := range configs {
		_, err = sec.NewKey(k, v)
		if err != nil {
			return err
		}
	}

	configPath := filepath.Join(targetPath, "my.cnf")
	if err = myCnf.SaveTo(configPath); err != nil {
		return err
	}

	cmd := exec.Command("sudo", "chown", "-R", "mysql:mysql", targetPath)
	if err = cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("sudo", "-u", "mysql", filepath.Join(mysqlPath, "bin/mysqld_safe"), fmt.Sprintf("--defaults-file=%s", configPath))
	log.Printf("run cmd: %s", cmd)
	if err = cmd.Start(); err != nil {
		return err
	}

	time.Sleep(time.Second * 30)

	log.Printf("restore completed\nport: %d", freePort)
	return nil
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
