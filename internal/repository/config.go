package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Identifer   string `json:"identifer"`
	Version     string `json:"version"`
	LoginPath   string `json:"login_path"`
	DbHostName  string `json:"db_hostname"`
	DbUser      string `json:"db_user"`
	Throttle    int    `json:"throttle"`
	TryCompress bool   `json:"try_compress"`

	BinPath        string `json:"bin_path"`
	DataPath       string `json:"data_path"`
	BackupUser     string `json:"backup_user"`
	BackupHostName string `json:"backup_hostname"`
}

func saveConfigToRepo(config *Config, repoPath string) error {
	d, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(repoPath, "config"), d, 0664); err != nil {
		return err
	}

	return nil
}

func loadConfigFromRepo(config *Config, path string) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(d, config); err != nil {
		return err
	}

	return nil
}
