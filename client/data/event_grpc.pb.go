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

// SendEventClient is the client API for SendEvent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SendEventClient interface {
	Event(ctx context.Context, in *SendEventReq, opts ...grpc.CallOption) (*SendEventRes, error)
}

type sendEventClient struct {
	cc grpc.ClientConnInterface
}

func NewSendEventClient(cc grpc.ClientConnInterface) SendEventClient {
	return &sendEventClient{cc}
}

func (c *sendEventClient) Event(ctx context.Context, in *SendEventReq, opts ...grpc.CallOption) (*SendEventRes, error) {
	out := new(SendEventRes)
	err := c.cc.Invoke(ctx, "/data.SendEvent/Event", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SendEventServer is the server API for SendEvent service.
// All implementations must embed UnimplementedSendEventServer
// for forward compatibility
type SendEventServer interface {
	Event(context.Context, *SendEventReq) (*SendEventRes, error)
	mustEmbedUnimplementedSendEventServer()
}

// UnimplementedSendEventServer must be embedded to have forward compatible implementations.
type UnimplementedSendEventServer struct {
}

func (UnimplementedSendEventServer) Event(context.Context, *SendEventReq) (*SendEventRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Event not implemented")
}
func (UnimplementedSendEventServer) mustEmbedUnimplementedSendEventServer() {}

// UnsafeSendEventServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SendEventServer will
// result in compilation errors.
type UnsafeSendEventServer interface {
	mustEmbedUnimplementedSendEventServer()
}

func RegisterSendEventServer(s grpc.ServiceRegistrar, srv SendEventServer) {
	s.RegisterService(&SendEvent_ServiceDesc, srv)
}

func _SendEvent_Event_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendEventReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendEventServer).Event(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/data.SendEvent/Event",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendEventServer).Event(ctx, req.(*SendEventReq))
	}
	return interceptor(ctx, in, info, handler)
}

// SendEvent_ServiceDesc is the grpc.ServiceDesc for SendEvent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SendEvent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "data.SendEvent",
	HandlerType: (*SendEventServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Event",
			Handler:    _SendEvent_Event_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event.proto",
}
