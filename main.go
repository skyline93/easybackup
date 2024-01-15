package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skyline93/easybackup/internal/raft"
)

var raftId string

var cluster = map[string]string{
	"1": "127.0.0.1:5000",
	// "2": "127.0.0.1:6000",
	// "3": "127.0.0.1:7000",
	// "4": "127.0.0.1:8000",
	// "5": "127.0.0.1:9000",
}

func init() {
	flag.StringVar(&raftId, "raft_id", "1", "raft id")
}

func main() {
	flag.Parse()
	if raftId == "" {
		fmt.Println("raft_id error")
		os.Exit(1)
	}

	rootDir := "./"

	id := raftId
	addr := cluster[raftId]

	storeDir := filepath.Join(rootDir, "nodes", id)

	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		os.MkdirAll(storeDir, 0700)
	}

	fsm := raft.NewFSM()
	rf, err := raft.NewRaft(addr, id, storeDir, fsm)
	if err != nil {
		panic(err)
	}

	raft.Bootstrap(rf, cluster)

	for {
	}
}
