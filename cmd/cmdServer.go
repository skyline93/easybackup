package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/skyline93/easybackup/internal/log"
	"github.com/skyline93/easybackup/internal/raft"
	"github.com/skyline93/easybackup/internal/server"
	"github.com/skyline93/easybackup/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdRunServer = &cobra.Command{
	Use:   "run --addr 0.0.0.0:50051 --id 1 --raft_addr 127.0.0.1:5000",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
		rootDir := "./"

		storeDir := filepath.Join(rootDir, "nodes", raftId)

		if _, err := os.Stat(storeDir); os.IsNotExist(err) {
			os.MkdirAll(storeDir, 0700)
		}

		fsm := raft.NewFSM()
		rf, err := raft.NewRaft(raftAddr, raftId, storeDir, fsm)
		if err != nil {
			panic(err)
		}

		raft.Bootstrap(rf, raftId, raftAddr)

		srv := server.NewServer(rf)
		if err := srv.Run(serverAddr); err != nil {
			log.Error("run server failed")
			os.Exit(1)
		}
	},
}

var cmdAddNode = &cobra.Command{
	Use:   "addnode",
	Short: "addnode",
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		conn, err := grpc.Dial(serverAddr, opts...)
		if err != nil {
			panic(err)
		}

		client := proto.NewClusterManageServiceClient(conn)

		resp, err := client.AddNode(context.TODO(), &proto.AddNodeRequest{RaftId: raftId, RaftAddr: raftAddr})
		if err != nil {
			panic(err)
		}

		log.Infof("isOk: %v", resp.IsOk)
	},
}

var (
	serverAddr string
	raftId     string
	raftAddr   string
)

func init() {
	cmdRoot.AddCommand(cmdServer)
	cmdServer.AddCommand(cmdRunServer)
	cmdServer.AddCommand(cmdAddNode)

	cmdRunServer.Flags().StringVarP(&serverAddr, "addr", "a", "0.0.0.0:50051", "server addr")
	cmdRunServer.Flags().StringVarP(&raftId, "id", "i", "1", "raft id")
	cmdRunServer.Flags().StringVarP(&raftAddr, "raft_addr", "r", "127.0.0.1:5000", "raft addr")

	cmdAddNode.Flags().StringVarP(&serverAddr, "addr", "a", "0.0.0.0:50051", "server addr")
	cmdAddNode.Flags().StringVarP(&raftId, "id", "i", "1", "raft id")
	cmdAddNode.Flags().StringVarP(&raftAddr, "raft_addr", "r", "127.0.0.1:5000", "raft addr")

	cmdRunServer.MarkFlagRequired("addr")
	cmdRunServer.MarkFlagRequired("id")
	cmdRunServer.MarkFlagRequired("raft_addr")

	cmdAddNode.MarkFlagRequired("addr")
	cmdAddNode.MarkFlagRequired("id")
	cmdAddNode.MarkFlagRequired("raft_addr")
}
