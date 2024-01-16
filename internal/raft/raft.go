package raft

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

func NewRaft(addr string, id string, storeDir string, fsm *FSM) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(id)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(addr, tcpAddr, 2, 5*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(storeDir, 2, os.Stderr)
	if err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(storeDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(storeDir, "raft-stable.db"))
	if err != nil {
		return nil, err
	}

	rf, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	return rf, nil
}

func Bootstrap(rf *raft.Raft, raftId string, raftAddr string) {
	servers := rf.GetConfiguration().Configuration().Servers
	if len(servers) > 0 {
		log.Fatal("Error checking existing server")
		return
	}

	var configuration raft.Configuration
	server := raft.Server{
		ID:      raft.ServerID(raftId),
		Address: raft.ServerAddress(raftAddr),
	}
	configuration.Servers = append(configuration.Servers, server)

	rf.BootstrapCluster(configuration)
}
