package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/skyline93/mysql-xtrabackup/internal/stor"
)

const (
	TypeBackupSetFull = "full"
	TypeBackupSetIncr = "incr"
)

type BackupSet struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	FromLSN    string `json:"from_lsn"`
	ToLSN      string `json:"to_lsn"`
	Size       int64  `json:"size"`
	BackupTime string `json:"backup_time"`
}

type Repository struct {
	col    *stor.Collection
	Config *Config
	Path   string `json:"path"`
	Name   string `json:"name"`
}

func NewBackupSet(backupSetType string) *BackupSet {
	return &BackupSet{
		Id:   uuid.New().String(),
		Type: backupSetType,
	}
}

func NewRepository(name string, config *Config) *Repository {
	return &Repository{
		col:    stor.NewCollection(),
		Name:   name,
		Config: config,
	}
}

func LoadRepository(repo *Repository, path string) error {
	indexPath := filepath.Join(path, "index")
	col := stor.Collection{}

	if err := stor.Deserialize(&col, indexPath); err != nil {
		return err
	}

	confPath := filepath.Join(path, "config")
	conf := Config{}
	if err := loadConfigFromRepo(&conf, confPath); err != nil {
		return err
	}

	repo.col = &col
	repo.Config = &conf
	repo.Name = conf.Identifer
	repo.Path = path

	return nil
}

func (r *Repository) Init(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	repoPath := filepath.Join(absPath, r.Name)
	r.Path = repoPath

	if err = os.MkdirAll(repoPath, 0764); err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Join(r.Path, "data"), 0764); err != nil {
		return err
	}

	if err = saveConfigToRepo(r.Config, r.Path); err != nil {
		return err
	}

	if err = stor.Serialize(r.col, filepath.Join(r.Path, "index")); err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddBackupSet(backupSet *BackupSet) error {
	v, err := json.Marshal(backupSet)
	if err != nil {
		return err
	}

	if backupSet.Type == TypeBackupSetFull {
		_, err := r.col.NewNode(backupSet.Id, v, true)
		if err != nil {
			return err
		}
	} else if backupSet.Type == TypeBackupSetIncr {
		_, err := r.col.NewNode(backupSet.Id, v, false)
		if err != nil {
			return err
		}
	}

	if err := stor.Serialize(r.col, filepath.Join(r.Path, "index")); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBackupSet(backupSetId string) (*BackupSet, error) {
	n := r.col.GetNode(backupSetId)
	var backupSet BackupSet
	if err := json.Unmarshal(n.Data, &backupSet); err != nil {
		return nil, err
	}

	return &backupSet, nil
}

func (r *Repository) GetBeforeBackupSet(backupSetId string) ([]BackupSet, error) {
	var backupSets []BackupSet

	nodes := r.col.GetBeforeNodes(backupSetId)

	for _, n := range nodes {
		var backupSet BackupSet
		if err := json.Unmarshal(n.Data, &backupSet); err != nil {
			return nil, err
		}

		backupSets = append(backupSets, backupSet)
	}

	return backupSets, nil
}

func (r *Repository) GetLastBackupSet() (*BackupSet, error) {
	n := r.col.GetLastNode()
	if n == nil {
		return nil, errors.New("last backupset is not found")
	}

	var backupSet BackupSet
	if err := json.Unmarshal(n.Data, &backupSet); err != nil {
		return nil, err
	}
	return &backupSet, nil
}

func (r *Repository) DataPath() string {
	return filepath.Join(r.Path, "data")
}

func (r *Repository) ListBackupSets() ([]BackupSet, error) {
	var backupSets []BackupSet

	ns := r.col.GetAllNodes()

	for _, n := range ns {
		var backupSet BackupSet
		if err := json.Unmarshal(n.Data, &backupSet); err != nil {
			return nil, err
		}

		backupSets = append(backupSets, backupSet)
	}

	return backupSets, nil
}
