// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package raftpb

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

// RaftCoreServiceClient is the client API for RaftCoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RaftCoreServiceClient interface {
	// *
	// tell server to time out immediately to start a election.
	// no need of request and response parameters,
	// RequestVote and RequestVoteRes are just placeholders.
	HeartBeat(ctx context.Context, in *RequestVote, opts ...grpc.CallOption) (*RequestVoteRes, error)
	// *
	// pre-vote phase to prevent disruptions, we use same parameter
	// as request_vote since most of the logic resembles.
	PreVote(ctx context.Context, in *RequestVote, opts ...grpc.CallOption) (*RequestVoteRes, error)
	RequestVote(ctx context.Context, in *RequestVote, opts ...grpc.CallOption) (*RequestVoteRes, error)
	AppendEntries(ctx context.Context, in *AppendEntries, opts ...grpc.CallOption) (*AppendEntriesRes, error)
}

type raftCoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftCoreServiceClient(cc grpc.ClientConnInterface) RaftCoreServiceClient {
	return &raftCoreServiceClient{cc}
}

func (c *raftCoreServiceClient) HeartBeat(ctx context.Context, in *RequestVote, opts ...grpc.CallOption) (*RequestVoteRes, error) {
	out := new(RequestVoteRes)
	err := c.cc.Invoke(ctx, "/raftpb.RaftCoreService/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftCoreServiceClient) PreVote(ctx context.Context, in *RequestVote, opts ...grpc.CallOption) (*RequestVoteRes, error) {
	out := new(RequestVoteRes)
	err := c.cc.Invoke(ctx, "/raftpb.RaftCoreService/preVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftCoreServiceClient) RequestVote(ctx context.Context, in *RequestVote, opts ...grpc.CallOption) (*RequestVoteRes, error) {
	out := new(RequestVoteRes)
	err := c.cc.Invoke(ctx, "/raftpb.RaftCoreService/requestVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftCoreServiceClient) AppendEntries(ctx context.Context, in *AppendEntries, opts ...grpc.CallOption) (*AppendEntriesRes, error) {
	out := new(AppendEntriesRes)
	err := c.cc.Invoke(ctx, "/raftpb.RaftCoreService/appendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftCoreServiceServer is the server API for RaftCoreService service.
// All implementations must embed UnimplementedRaftCoreServiceServer
// for forward compatibility
type RaftCoreServiceServer interface {
	// *
	// tell server to time out immediately to start a election.
	// no need of request and response parameters,
	// RequestVote and RequestVoteRes are just placeholders.
	HeartBeat(context.Context, *RequestVote) (*RequestVoteRes, error)
	// *
	// pre-vote phase to prevent disruptions, we use same parameter
	// as request_vote since most of the logic resembles.
	PreVote(context.Context, *RequestVote) (*RequestVoteRes, error)
	RequestVote(context.Context, *RequestVote) (*RequestVoteRes, error)
	AppendEntries(context.Context, *AppendEntries) (*AppendEntriesRes, error)
	mustEmbedUnimplementedRaftCoreServiceServer()
}

// UnimplementedRaftCoreServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRaftCoreServiceServer struct {
}

func (UnimplementedRaftCoreServiceServer) HeartBeat(context.Context, *RequestVote) (*RequestVoteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}
func (UnimplementedRaftCoreServiceServer) PreVote(context.Context, *RequestVote) (*RequestVoteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreVote not implemented")
}
func (UnimplementedRaftCoreServiceServer) RequestVote(context.Context, *RequestVote) (*RequestVoteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedRaftCoreServiceServer) AppendEntries(context.Context, *AppendEntries) (*AppendEntriesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedRaftCoreServiceServer) mustEmbedUnimplementedRaftCoreServiceServer() {}

// UnsafeRaftCoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RaftCoreServiceServer will
// result in compilation errors.
type UnsafeRaftCoreServiceServer interface {
	mustEmbedUnimplementedRaftCoreServiceServer()
}

func RegisterRaftCoreServiceServer(s grpc.ServiceRegistrar, srv RaftCoreServiceServer) {
	s.RegisterService(&RaftCoreService_ServiceDesc, srv)
}

func _RaftCoreService_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftCoreServiceServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raftpb.RaftCoreService/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftCoreServiceServer).HeartBeat(ctx, req.(*RequestVote))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftCoreService_PreVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftCoreServiceServer).PreVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raftpb.RaftCoreService/preVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftCoreServiceServer).PreVote(ctx, req.(*RequestVote))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftCoreService_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftCoreServiceServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raftpb.RaftCoreService/requestVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftCoreServiceServer).RequestVote(ctx, req.(*RequestVote))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftCoreService_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntries)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftCoreServiceServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/raftpb.RaftCoreService/appendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftCoreServiceServer).AppendEntries(ctx, req.(*AppendEntries))
	}
	return interceptor(ctx, in, info, handler)
}

// RaftCoreService_ServiceDesc is the grpc.ServiceDesc for RaftCoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RaftCoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raftpb.RaftCoreService",
	HandlerType: (*RaftCoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HeartBeat",
			Handler:    _RaftCoreService_HeartBeat_Handler,
		},
		{
			MethodName: "preVote",
			Handler:    _RaftCoreService_PreVote_Handler,
		},
		{
			MethodName: "requestVote",
			Handler:    _RaftCoreService_RequestVote_Handler,
		},
		{
			MethodName: "appendEntries",
			Handler:    _RaftCoreService_AppendEntries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft.proto",
}