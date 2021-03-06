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

// UpdateStatusClient is the client API for UpdateStatus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UpdateStatusClient interface {
	Update(ctx context.Context, in *UpdateStatusReq, opts ...grpc.CallOption) (*UpdateStatusRes, error)
}

type updateStatusClient struct {
	cc grpc.ClientConnInterface
}

func NewUpdateStatusClient(cc grpc.ClientConnInterface) UpdateStatusClient {
	return &updateStatusClient{cc}
}

func (c *updateStatusClient) Update(ctx context.Context, in *UpdateStatusReq, opts ...grpc.CallOption) (*UpdateStatusRes, error) {
	out := new(UpdateStatusRes)
	err := c.cc.Invoke(ctx, "/data.UpdateStatus/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateStatusServer is the server API for UpdateStatus service.
// All implementations must embed UnimplementedUpdateStatusServer
// for forward compatibility
type UpdateStatusServer interface {
	Update(context.Context, *UpdateStatusReq) (*UpdateStatusRes, error)
	mustEmbedUnimplementedUpdateStatusServer()
}

// UnimplementedUpdateStatusServer must be embedded to have forward compatible implementations.
type UnimplementedUpdateStatusServer struct {
}

func (UnimplementedUpdateStatusServer) Update(context.Context, *UpdateStatusReq) (*UpdateStatusRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedUpdateStatusServer) mustEmbedUnimplementedUpdateStatusServer() {}

// UnsafeUpdateStatusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UpdateStatusServer will
// result in compilation errors.
type UnsafeUpdateStatusServer interface {
	mustEmbedUnimplementedUpdateStatusServer()
}

func RegisterUpdateStatusServer(s grpc.ServiceRegistrar, srv UpdateStatusServer) {
	s.RegisterService(&UpdateStatus_ServiceDesc, srv)
}

func _UpdateStatus_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateStatusServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/data.UpdateStatus/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpdateStatusServer).Update(ctx, req.(*UpdateStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

// UpdateStatus_ServiceDesc is the grpc.ServiceDesc for UpdateStatus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UpdateStatus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data.UpdateStatus",
	HandlerType: (*UpdateStatusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _UpdateStatus_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "status.proto",
}
