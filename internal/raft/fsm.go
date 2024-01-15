package raft

import (
	"io"

	"github.com/hashicorp/raft"
)

type FSM struct {
}

func NewFSM() *FSM {
	return &FSM{}
}

func (f *FSM) Apply(l *raft.Log) interface{} {
	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *FSM) Restore(io.ReadCloser) error {
	return nil
}
