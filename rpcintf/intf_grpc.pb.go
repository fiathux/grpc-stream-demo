// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: protos/intf.proto

package rpcintf

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

// StreamingClient is the client API for Streaming service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamingClient interface {
	TestStream(ctx context.Context, opts ...grpc.CallOption) (Streaming_TestStreamClient, error)
}

type streamingClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamingClient(cc grpc.ClientConnInterface) StreamingClient {
	return &streamingClient{cc}
}

func (c *streamingClient) TestStream(ctx context.Context, opts ...grpc.CallOption) (Streaming_TestStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Streaming_ServiceDesc.Streams[0], "/rpcintf.Streaming/TestStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamingTestStreamClient{stream}
	return x, nil
}

type Streaming_TestStreamClient interface {
	Send(*MsgAgent) error
	Recv() (*MsgAgent, error)
	grpc.ClientStream
}

type streamingTestStreamClient struct {
	grpc.ClientStream
}

func (x *streamingTestStreamClient) Send(m *MsgAgent) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamingTestStreamClient) Recv() (*MsgAgent, error) {
	m := new(MsgAgent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamingServer is the server API for Streaming service.
// All implementations must embed UnimplementedStreamingServer
// for forward compatibility
type StreamingServer interface {
	TestStream(Streaming_TestStreamServer) error
	mustEmbedUnimplementedStreamingServer()
}

// UnimplementedStreamingServer must be embedded to have forward compatible implementations.
type UnimplementedStreamingServer struct {
}

func (UnimplementedStreamingServer) TestStream(Streaming_TestStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method TestStream not implemented")
}
func (UnimplementedStreamingServer) mustEmbedUnimplementedStreamingServer() {}

// UnsafeStreamingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamingServer will
// result in compilation errors.
type UnsafeStreamingServer interface {
	mustEmbedUnimplementedStreamingServer()
}

func RegisterStreamingServer(s grpc.ServiceRegistrar, srv StreamingServer) {
	s.RegisterService(&Streaming_ServiceDesc, srv)
}

func _Streaming_TestStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamingServer).TestStream(&streamingTestStreamServer{stream})
}

type Streaming_TestStreamServer interface {
	Send(*MsgAgent) error
	Recv() (*MsgAgent, error)
	grpc.ServerStream
}

type streamingTestStreamServer struct {
	grpc.ServerStream
}

func (x *streamingTestStreamServer) Send(m *MsgAgent) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamingTestStreamServer) Recv() (*MsgAgent, error) {
	m := new(MsgAgent)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Streaming_ServiceDesc is the grpc.ServiceDesc for Streaming service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Streaming_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpcintf.Streaming",
	HandlerType: (*StreamingServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TestStream",
			Handler:       _Streaming_TestStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protos/intf.proto",
}