// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: rpc/embeddings.proto

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	EmbeddingsService_GetEmbeddings_FullMethodName = "/EmbeddingsService/GetEmbeddings"
)

// EmbeddingsServiceClient is the client API for EmbeddingsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmbeddingsServiceClient interface {
	GetEmbeddings(ctx context.Context, in *GetEmbeddingsRequest, opts ...grpc.CallOption) (*GetEmbeddingsResponse, error)
}

type embeddingsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmbeddingsServiceClient(cc grpc.ClientConnInterface) EmbeddingsServiceClient {
	return &embeddingsServiceClient{cc}
}

func (c *embeddingsServiceClient) GetEmbeddings(ctx context.Context, in *GetEmbeddingsRequest, opts ...grpc.CallOption) (*GetEmbeddingsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEmbeddingsResponse)
	err := c.cc.Invoke(ctx, EmbeddingsService_GetEmbeddings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmbeddingsServiceServer is the server API for EmbeddingsService service.
// All implementations must embed UnimplementedEmbeddingsServiceServer
// for forward compatibility
type EmbeddingsServiceServer interface {
	GetEmbeddings(context.Context, *GetEmbeddingsRequest) (*GetEmbeddingsResponse, error)
	mustEmbedUnimplementedEmbeddingsServiceServer()
}

// UnimplementedEmbeddingsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEmbeddingsServiceServer struct {
}

func (UnimplementedEmbeddingsServiceServer) GetEmbeddings(context.Context, *GetEmbeddingsRequest) (*GetEmbeddingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmbeddings not implemented")
}
func (UnimplementedEmbeddingsServiceServer) mustEmbedUnimplementedEmbeddingsServiceServer() {}

// UnsafeEmbeddingsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmbeddingsServiceServer will
// result in compilation errors.
type UnsafeEmbeddingsServiceServer interface {
	mustEmbedUnimplementedEmbeddingsServiceServer()
}

func RegisterEmbeddingsServiceServer(s grpc.ServiceRegistrar, srv EmbeddingsServiceServer) {
	s.RegisterService(&EmbeddingsService_ServiceDesc, srv)
}

func _EmbeddingsService_GetEmbeddings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEmbeddingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbeddingsServiceServer).GetEmbeddings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmbeddingsService_GetEmbeddings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmbeddingsServiceServer).GetEmbeddings(ctx, req.(*GetEmbeddingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EmbeddingsService_ServiceDesc is the grpc.ServiceDesc for EmbeddingsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmbeddingsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "EmbeddingsService",
	HandlerType: (*EmbeddingsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEmbeddings",
			Handler:    _EmbeddingsService_GetEmbeddings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/embeddings.proto",
}
