// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: proto/easybackup.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ClusterManageService_AddNode_FullMethodName = "/proto.ClusterManageService/AddNode"
)

// ClusterManageServiceClient is the client API for ClusterManageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterManageServiceClient interface {
	AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error)
}

type clusterManageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterManageServiceClient(cc grpc.ClientConnInterface) ClusterManageServiceClient {
	return &clusterManageServiceClient{cc}
}

func (c *clusterManageServiceClient) AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error) {
	out := new(AddNodeResponse)
	err := c.cc.Invoke(ctx, ClusterManageService_AddNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterManageServiceServer is the server API for ClusterManageService service.
// All implementations must embed UnimplementedClusterManageServiceServer
// for forward compatibility
type ClusterManageServiceServer interface {
	AddNode(context.Context, *AddNodeRequest) (*AddNodeResponse, error)
	mustEmbedUnimplementedClusterManageServiceServer()
}

// UnimplementedClusterManageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedClusterManageServiceServer struct {
}

func (UnimplementedClusterManageServiceServer) AddNode(context.Context, *AddNodeRequest) (*AddNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNode not implemented")
}
func (UnimplementedClusterManageServiceServer) mustEmbedUnimplementedClusterManageServiceServer() {}

// UnsafeClusterManageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterManageServiceServer will
// result in compilation errors.
type UnsafeClusterManageServiceServer interface {
	mustEmbedUnimplementedClusterManageServiceServer()
}

func RegisterClusterManageServiceServer(s grpc.ServiceRegistrar, srv ClusterManageServiceServer) {
	s.RegisterService(&ClusterManageService_ServiceDesc, srv)
}

func _ClusterManageService_AddNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterManageServiceServer).AddNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterManageService_AddNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterManageServiceServer).AddNode(ctx, req.(*AddNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClusterManageService_ServiceDesc is the grpc.ServiceDesc for ClusterManageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterManageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ClusterManageService",
	HandlerType: (*ClusterManageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddNode",
			Handler:    _ClusterManageService_AddNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/easybackup.proto",
}