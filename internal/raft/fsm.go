package raft

import (
	"io"
	"strings"

	"github.com/hashicorp/raft"
	"github.com/skyline93/gokv"
)

type FSM struct {
	Stor Stor
}

func NewFSM() *FSM {
	return &FSM{
		Stor: Stor{Cache: gokv.New(1000)},
	}
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	data := strings.Split(string(l.Data), ",")

	op := data[0]
	if op == "set" {
		key := data[1]
		value := data[2]
		f.Stor.Put(key, value)
	}

	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &f.Stor, nil
}

func (f *FSM) Restore(io.ReadCloser) error {
	return nil
}

type Stor struct {
	*gokv.Cache
}

func (s *Stor) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (s *Stor) Release() {}
