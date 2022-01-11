// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package data

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

// SubmitJilClient is the client API for SubmitJil service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubmitJilClient interface {
	Submit(ctx context.Context, in *SubmitJilReq, opts ...grpc.CallOption) (*SubmitJilRes, error)
}

type submitJilClient struct {
	cc grpc.ClientConnInterface
}

func NewSubmitJilClient(cc grpc.ClientConnInterface) SubmitJilClient {
	return &submitJilClient{cc}
}

func (c *submitJilClient) Submit(ctx context.Context, in *SubmitJilReq, opts ...grpc.CallOption) (*SubmitJilRes, error) {
	out := new(SubmitJilRes)
	err := c.cc.Invoke(ctx, "/data.SubmitJil/Submit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubmitJilServer is the server API for SubmitJil service.
// All implementations must embed UnimplementedSubmitJilServer
// for forward compatibility
type SubmitJilServer interface {
	Submit(context.Context, *SubmitJilReq) (*SubmitJilRes, error)
	mustEmbedUnimplementedSubmitJilServer()
}

// UnimplementedSubmitJilServer must be embedded to have forward compatible implementations.
type UnimplementedSubmitJilServer struct {
}

func (UnimplementedSubmitJilServer) Submit(context.Context, *SubmitJilReq) (*SubmitJilRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Submit not implemented")
}
func (UnimplementedSubmitJilServer) mustEmbedUnimplementedSubmitJilServer() {}

// UnsafeSubmitJilServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubmitJilServer will
// result in compilation errors.
type UnsafeSubmitJilServer interface {
	mustEmbedUnimplementedSubmitJilServer()
}

func RegisterSubmitJilServer(s grpc.ServiceRegistrar, srv SubmitJilServer) {
	s.RegisterService(&SubmitJil_ServiceDesc, srv)
}

func _SubmitJil_Submit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitJilReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubmitJilServer).Submit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/data.SubmitJil/Submit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubmitJilServer).Submit(ctx, req.(*SubmitJilReq))
	}
	return interceptor(ctx, in, info, handler)
}

// SubmitJil_ServiceDesc is the grpc.ServiceDesc for SubmitJil service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SubmitJil_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data.SubmitJil",
	HandlerType: (*SubmitJilServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Submit",
			Handler:    _SubmitJil_Submit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "jil.proto",
}
