package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/skyline93/easybackup/internal/stor"
)

const (
	TypeBackupSetFull = "full"
	TypeBackupSetIncr = "incr"

	TypeData = "data"
	TypeLog  = "log"
)

type BackupSet struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	FromLSN    string `json:"from_lsn"`
	ToLSN      string `json:"to_lsn"`
	Size       int64  `json:"size"`
	BackupTime string `json:"backup_time"`
	DataType   string `json:"data_type"`
	Snapshot   string `json:"snapshot"`
}

type Repository struct {
	DataCol *stor.Collection
	LogCol  *stor.Collection
	Config  *Config
	Path    string `json:"path"`
	Name    string `json:"name"`
}

func NewBackupSet(backupSetType string) *BackupSet {
	return &BackupSet{
		Id:   uuid.New().String(),
		Type: backupSetType,
	}
}

func NewRepository(name string, config *Config) *Repository {
	return &Repository{
		DataCol: stor.NewCollection(),
		LogCol:  stor.NewCollection(),
		Name:    name,
		Config:  config,
	}
}

func LoadRepository(repo *Repository, path string) error {
	dataIndexPath := filepath.Join(path, "data.index")
	dataCol := stor.Collection{}

	if err := stor.Deserialize(&dataCol, dataIndexPath); err != nil {
		return err
	}

	logIndexPath := filepath.Join(path, "log.index")
	logCol := stor.Collection{}

	if err := stor.Deserialize(&logCol, logIndexPath); err != nil {
		return err
	}

	confPath := filepath.Join(path, "config")
	conf := Config{}
	if err := loadConfigFromRepo(&conf, confPath); err != nil {
		return err
	}

	repo.DataCol = &dataCol
	repo.LogCol = &logCol
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

	if err = os.MkdirAll(filepath.Join(r.Path, "snapshot"), 0764); err != nil {
		return err
	}

	if err = saveConfigToRepo(r.Config, r.Path); err != nil {
		return err
	}

	if err = stor.Serialize(r.DataCol, filepath.Join(r.Path, "data.index")); err != nil {
		return err
	}

	if err = stor.Serialize(r.LogCol, filepath.Join(r.Path, "log.index")); err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddBackupSet(backupSet *BackupSet) error {
	v, err := json.Marshal(backupSet)
	if err != nil {
		return err
	}

	var (
		col       *stor.Collection
		indexName string
	)

	if backupSet.DataType == TypeData {
		col = r.DataCol
		indexName = "data.index"
	} else if backupSet.DataType == TypeLog {
		col = r.LogCol
		indexName = "log.index"
	}

	if backupSet.Type == TypeBackupSetFull {
		_, err := col.NewNode(backupSet.Id, v, true)
		if err != nil {
			return err
		}
	} else if backupSet.Type == TypeBackupSetIncr {
		_, err := col.NewNode(backupSet.Id, v, false)
		if err != nil {
			return err
		}
	}

	if err := stor.Serialize(col, filepath.Join(r.Path, indexName)); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBackupSet(backupDataType string, backupSetId string) (*BackupSet, error) {
	var col *stor.Collection

	if backupDataType == TypeData {
		col = r.DataCol
	} else if backupDataType == TypeLog {
		col = r.LogCol
	}

	n := col.GetNode(backupSetId)
	var backupSet BackupSet
	if err := json.Unmarshal(n.Data, &backupSet); err != nil {
		return nil, err
	}

	return &backupSet, nil
}

func (r *Repository) GetBeforeBackupSet(backupDataType string, backupSetId string) ([]BackupSet, error) {
	var (
		backupSets []BackupSet
		col        *stor.Collection
	)

	if backupDataType == TypeData {
		col = r.DataCol
	} else if backupDataType == TypeLog {
		col = r.LogCol
	}

	nodes := col.GetBeforeNodes(backupSetId)

	for _, n := range nodes {
		var backupSet BackupSet
		if err := json.Unmarshal(n.Data, &backupSet); err != nil {
			return nil, err
		}

		backupSets = append(backupSets, backupSet)
	}

	return backupSets, nil
}

func (r *Repository) GetLastBackupSet(backupDataType string) (*BackupSet, error) {
	var col *stor.Collection

	if backupDataType == TypeData {
		col = r.DataCol
	} else if backupDataType == TypeLog {
		col = r.LogCol
	}

	n := col.GetLastNode()
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

func (r *Repository) SnapshotPath() string {
	return filepath.Join(r.Path, "snapshot")
}

func (r *Repository) ListBackupSets(backupDataType string) ([]BackupSet, error) {
	var (
		backupSets []BackupSet
		col        *stor.Collection
	)

	if backupDataType == TypeData {
		col = r.DataCol
	} else if backupDataType == TypeLog {
		col = r.LogCol
	}

	ns := col.GetAllNodes()

	for _, n := range ns {
		var backupSet BackupSet
		if err := json.Unmarshal(n.Data, &backupSet); err != nil {
			return nil, err
		}

		backupSets = append(backupSets, backupSet)
	}

	return backupSets, nil
}
