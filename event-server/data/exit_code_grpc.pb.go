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

// RunStatusClient is the client API for RunStatus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RunStatusClient interface {
	Update(ctx context.Context, in *RunStatusReq, opts ...grpc.CallOption) (*RunStatusRes, error)
}

type runStatusClient struct {
	cc grpc.ClientConnInterface
}

func NewRunStatusClient(cc grpc.ClientConnInterface) RunStatusClient {
	return &runStatusClient{cc}
}

func (c *runStatusClient) Update(ctx context.Context, in *RunStatusReq, opts ...grpc.CallOption) (*RunStatusRes, error) {
	out := new(RunStatusRes)
	err := c.cc.Invoke(ctx, "/data.RunStatus/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RunStatusServer is the server API for RunStatus service.
// All implementations must embed UnimplementedRunStatusServer
// for forward compatibility
type RunStatusServer interface {
	Update(context.Context, *RunStatusReq) (*RunStatusRes, error)
	mustEmbedUnimplementedRunStatusServer()
}

// UnimplementedRunStatusServer must be embedded to have forward compatible implementations.
type UnimplementedRunStatusServer struct {
}

func (UnimplementedRunStatusServer) Update(context.Context, *RunStatusReq) (*RunStatusRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedRunStatusServer) mustEmbedUnimplementedRunStatusServer() {}

// UnsafeRunStatusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RunStatusServer will
// result in compilation errors.
type UnsafeRunStatusServer interface {
	mustEmbedUnimplementedRunStatusServer()
}

func RegisterRunStatusServer(s grpc.ServiceRegistrar, srv RunStatusServer) {
	s.RegisterService(&RunStatus_ServiceDesc, srv)
}

func _RunStatus_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunStatusServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/data.RunStatus/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunStatusServer).Update(ctx, req.(*RunStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

// RunStatus_ServiceDesc is the grpc.ServiceDesc for RunStatus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RunStatus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data.RunStatus",
	HandlerType: (*RunStatusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _RunStatus_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "exit_code.proto",
}
