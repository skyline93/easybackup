package repository

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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
	BackupTime int64  `json:"backup_time"`
	DataType   string `json:"data_type"`
	Snapshot   string `json:"snapshot"`

	repo *Repository
}

func (b *BackupSet) String() string {
	d, _ := json.Marshal(b)
	return string(d)
}

func (bs *BackupSet) Update(fieldsToUpdate map[string]interface{}) error {
	if fieldsToUpdate == nil {
		return errors.New("fieldsToUpdate is nil")
	}

	for key, value := range fieldsToUpdate {
		switch key {
		case "FromLSN":
			bs.FromLSN = value.(string)
		case "ToLSN":
			bs.ToLSN = value.(string)
		case "Size":
			bs.Size = value.(int64)
		case "BackupTime":
			bs.BackupTime = value.(int64)
		case "DataType":
			bs.DataType = value.(string)
		case "Snapshot":
			bs.Snapshot = value.(string)
		default:
			return errors.New("unsupported field: " + key)
		}
	}

	return nil
}

type Repository struct {
	colMap map[string]*stor.Collection
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
		colMap: make(map[string]*stor.Collection),
		Name:   name,
		Config: config,
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

	repo.colMap[TypeData] = &dataCol
	repo.colMap[TypeLog] = &logCol
	repo.Config = &conf
	repo.Name = conf.Identifer
	repo.Path = path

	return nil
}

func (r *Repository) getCollection(backupDataType string) (*stor.Collection, error) {
	col, found := r.colMap[backupDataType]
	if !found {
		return nil, errors.New("unsupported backup data type: " + backupDataType)
	}
	return col, nil
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

	if err = stor.Serialize(r.colMap[TypeData], filepath.Join(r.Path, "data.index")); err != nil {
		return err
	}

	if err = stor.Serialize(r.colMap[TypeLog], filepath.Join(r.Path, "log.index")); err != nil {
		return err
	}

	return nil
}

func (bs *BackupSet) Path() string {
	path, _ := filepath.Abs(filepath.Join(bs.repo.DataPath(), bs.Id))
	return path
}

func (bs *BackupSet) GetSize() (uint64, error) {
	cmd := exec.Command("du", "-sb", bs.Path())

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

func (r *Repository) CreateBackupSet(backupSetType string) (*BackupSet, error) {
	id := uuid.New().String()

	bakSet := &BackupSet{
		Id:   id,
		Type: backupSetType,
	}

	if _, err := os.Stat(bakSet.Path()); os.IsNotExist(err) {
		err = os.MkdirAll(bakSet.Path(), 0775)
		if err != nil {
			return nil, err
		}
	}

	v, err := json.Marshal(bakSet)
	if err != nil {
		return nil, err
	}

	col, err := r.getCollection(bakSet.DataType)
	if err != nil {
		return nil, err
	}

	if bakSet.Type == TypeBackupSetFull {
		_, err := col.NewNode(bakSet.Id, v, true)
		if err != nil {
			return nil, err
		}
	} else if bakSet.Type == TypeBackupSetIncr {
		_, err := col.NewNode(bakSet.Id, v, false)
		if err != nil {
			return nil, err
		}
	}

	return bakSet, nil
}

func (c *Repository) DeleteBackupSet(bakSet *BackupSet) error {
	return os.RemoveAll(bakSet.Path())
}

func (r *Repository) GetBackupSet(backupDataType string, backupSetId string) (*BackupSet, error) {
	col, err := r.getCollection(backupDataType)
	if err != nil {
		return nil, err
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
	)

	col, err := r.getCollection(backupDataType)
	if err != nil {
		return nil, err
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
	col, err := r.getCollection(backupDataType)
	if err != nil {
		return nil, err
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
	)

	col, err := r.getCollection(backupDataType)
	if err != nil {
		return nil, err
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
