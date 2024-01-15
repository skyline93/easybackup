package server

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/raft"
	"github.com/skyline93/easybackup/internal/log"
	"github.com/skyline93/easybackup/proto"
)

type ClusterManageService struct {
	proto.UnimplementedClusterManageServiceServer
	rf *raft.Raft
}

func (s *ClusterManageService) AddNode(ctx context.Context, in *proto.AddNodeRequest) (*proto.AddNodeResponse, error) {
	s.rf.AddVoter(raft.ServerID(in.RaftId), raft.ServerAddress(in.RaftAddr), 0, time.Duration(60))
	log.Infof("add one non voter")
	return &proto.AddNodeResponse{IsOk: true}, nil
}

func NewServer(rf *raft.Raft) *ClusterManageService {
	return &ClusterManageService{rf: rf}
}

func (s *ClusterManageService) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("run server failed, err: %v", err)
		return err
	}

	log.Infof("listen at %s", addr)

	grpcServer := grpc.NewServer()
	proto.RegisterClusterManageServiceServer(grpcServer, s)

	return grpcServer.Serve(lis)
}
