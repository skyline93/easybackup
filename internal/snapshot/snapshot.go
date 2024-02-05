package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	os.FileInfo
	Path string
}

type NodeType int

const (
	NodeUniform NodeType = iota
	NodeNew
	NodeModified
	NodeDeleted
)

type Node struct {
	Path    string   `json:"path"`
	IsDir   bool     `json:"isDir"`
	ModTime int64    `json"mtime"`
	Type    NodeType `json:"type"`
}

type Snapshot struct {
	Path  string
	Time  int64           `json:"time"`
	Nodes map[string]Node `json:"nodes"`
}

func New(path string) *Snapshot {
	return &Snapshot{Path: path, Nodes: make(map[string]Node)}
}

func LoadSnapshot(path string) (*Snapshot, error) {
	var snapshot Snapshot

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	dec := json.NewDecoder(fp)
	if err := dec.Decode(&snapshot); err != nil {
		return nil, err
	}

	return &snapshot, nil
}

func (s *Snapshot) Snap(sourcePath string, lastSnap *Snapshot) error {
	fileChan := make(chan *FileInfo)

	s.Time = time.Now().UnixNano()

	go s.scanDir(fileChan, sourcePath)

	for {
		fi, ok := <-fileChan
		if !ok {
			break
		}

		fileInfo := *fi

		nodeType := NodeNew
		if lastSnap != nil {
			node, ok := lastSnap.Nodes[fi.Path]
			if !ok {
				nodeType = NodeNew
			} else if node.ModTime != fi.ModTime().UnixNano() {
				nodeType = NodeModified
			} else {
				nodeType = NodeUniform
			}
		}

		s.Nodes[fileInfo.Path] = Node{
			Path:    fileInfo.Path,
			IsDir:   fileInfo.IsDir(),
			ModTime: fileInfo.ModTime().UnixNano(),
			Type:    nodeType,
		}
	}

	snapData := fmt.Sprintf("%d", s.Time)
	if err := s.writeFile(s, filepath.Join(s.Path, snapData)); err != nil {
		return err
	}

	return nil
}

func (s *Snapshot) scanDir(fileChan chan *FileInfo, sourcePath string) error {
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		fileChan <- &FileInfo{FileInfo: fi, Path: path}

		return nil
	})
	if err != nil {
		return err
	}

	close(fileChan)
	return nil
}

func (s *Snapshot) writeFile(v any, path string) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0665)
	if err != nil {
		return err
	}
	defer fp.Close()

	enc := json.NewEncoder(fp)
	if err = enc.Encode(v); err != nil {
		return err
	}

	return nil
}
