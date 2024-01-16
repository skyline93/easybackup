package server

import (
	"context"
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/raft"
	"github.com/skyline93/easybackup/internal/log"
	myraft "github.com/skyline93/easybackup/internal/raft"
	"github.com/skyline93/easybackup/proto"
)

type ClusterManageService struct {
	proto.UnimplementedClusterManageServiceServer
	rf  *raft.Raft
	fsm *myraft.FSM
}

func (s *ClusterManageService) AddNode(ctx context.Context, in *proto.AddNodeRequest) (*proto.AddNodeResponse, error) {
	s.rf.AddVoter(raft.ServerID(in.RaftId), raft.ServerAddress(in.RaftAddr), 0, time.Duration(60))
	log.Infof("add one non voter")
	return &proto.AddNodeResponse{IsOk: true}, nil
}

func (s *ClusterManageService) Set(ctx context.Context, in *proto.SetRequest) (*proto.SetResponse, error) {
	data := "set" + "," + in.Key + "," + in.Value

	future := s.rf.Apply([]byte(data), 5*time.Second)
	if err := future.Error(); err != nil {
		log.Info("set error")
		return &proto.SetResponse{IsOk: false}, nil
	}

	return &proto.SetResponse{IsOk: true}, nil
}

func (s *ClusterManageService) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetResponse, error) {
	value := s.fsm.Stor.Get(in.Key)

	v, ok := value.(string)
	if !ok {
		return nil, errors.New("get error")
	}

	return &proto.GetResponse{Value: v}, nil
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
