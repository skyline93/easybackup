package main

import (
	"context"
	"os"

	"github.com/skyline93/easybackup/internal/log"
	"github.com/skyline93/easybackup/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cmdKV = &cobra.Command{
	Use:   "kv",
	Short: "kv",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var cmdSet = &cobra.Command{
	Use:   "set",
	Short: "set",
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		conn, err := grpc.Dial(serverAddr, opts...)
		if err != nil {
			panic(err)
		}

		client := proto.NewClusterManageServiceClient(conn)

		resp, err := client.Set(context.TODO(), &proto.SetRequest{Key: key, Value: value})
		if err != nil {
			panic(err)
		}

		log.Infof("isOk: %v", resp.IsOk)
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "get",
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		conn, err := grpc.Dial(serverAddr, opts...)
		if err != nil {
			panic(err)
		}

		client := proto.NewClusterManageServiceClient(conn)

		resp, err := client.Get(context.TODO(), &proto.GetRequest{Key: key})
		if err != nil {
			panic(err)
		}

		log.Infof("value: %v", resp.Value)
	},
}

var (
	key   string
	value string
)

func init() {
	cmdRoot.AddCommand(cmdKV)
	cmdKV.AddCommand(cmdSet)
	cmdKV.AddCommand(cmdGet)

	cmdSet.Flags().StringVarP(&serverAddr, "addr", "a", "0.0.0.0:50051", "server addr")
	cmdSet.Flags().StringVarP(&key, "key", "k", "", "key")
	cmdSet.Flags().StringVarP(&value, "value", "v", "", "value")

	cmdGet.Flags().StringVarP(&serverAddr, "addr", "a", "0.0.0.0:50051", "server addr")
	cmdGet.Flags().StringVarP(&key, "key", "k", "", "key")

	cmdSet.MarkFlagRequired("addr")
	cmdSet.MarkFlagRequired("key")
	cmdSet.MarkFlagRequired("value")

	cmdGet.MarkFlagRequired("addr")
	cmdGet.MarkFlagRequired("key")
}
