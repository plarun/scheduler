// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: entity.proto

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

// ParsedActionServiceClient is the client API for ParsedActionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ParsedActionServiceClient interface {
	// Submit sends the parsed entities to eventserver
	Submit(ctx context.Context, in *ParsedEntitiesRequest, opts ...grpc.CallOption) (*EntityActionResponse, error)
}

type parsedActionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewParsedActionServiceClient(cc grpc.ClientConnInterface) ParsedActionServiceClient {
	return &parsedActionServiceClient{cc}
}

func (c *parsedActionServiceClient) Submit(ctx context.Context, in *ParsedEntitiesRequest, opts ...grpc.CallOption) (*EntityActionResponse, error) {
	out := new(EntityActionResponse)
	err := c.cc.Invoke(ctx, "/proto.ParsedActionService/Submit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ParsedActionServiceServer is the server API for ParsedActionService service.
// All implementations must embed UnimplementedParsedActionServiceServer
// for forward compatibility
type ParsedActionServiceServer interface {
	// Submit sends the parsed entities to eventserver
	Submit(context.Context, *ParsedEntitiesRequest) (*EntityActionResponse, error)
	mustEmbedUnimplementedParsedActionServiceServer()
}

// UnimplementedParsedActionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedParsedActionServiceServer struct {
}

func (UnimplementedParsedActionServiceServer) Submit(context.Context, *ParsedEntitiesRequest) (*EntityActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Submit not implemented")
}
func (UnimplementedParsedActionServiceServer) mustEmbedUnimplementedParsedActionServiceServer() {}

// UnsafeParsedActionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ParsedActionServiceServer will
// result in compilation errors.
type UnsafeParsedActionServiceServer interface {
	mustEmbedUnimplementedParsedActionServiceServer()
}

func RegisterParsedActionServiceServer(s grpc.ServiceRegistrar, srv ParsedActionServiceServer) {
	s.RegisterService(&ParsedActionService_ServiceDesc, srv)
}

func _ParsedActionService_Submit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParsedEntitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParsedActionServiceServer).Submit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ParsedActionService/Submit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParsedActionServiceServer).Submit(ctx, req.(*ParsedEntitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ParsedActionService_ServiceDesc is the grpc.ServiceDesc for ParsedActionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ParsedActionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ParsedActionService",
	HandlerType: (*ParsedActionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Submit",
			Handler:    _ParsedActionService_Submit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entity.proto",
}

// ValidatedActionServiceClient is the client API for ValidatedActionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ValidatedActionServiceClient interface {
	// Route sends the parsed entities from client to validator
	Route(ctx context.Context, in *ParsedEntitiesRequest, opts ...grpc.CallOption) (*ValidatedEntitiesResponse, error)
}

type validatedActionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewValidatedActionServiceClient(cc grpc.ClientConnInterface) ValidatedActionServiceClient {
	return &validatedActionServiceClient{cc}
}

func (c *validatedActionServiceClient) Route(ctx context.Context, in *ParsedEntitiesRequest, opts ...grpc.CallOption) (*ValidatedEntitiesResponse, error) {
	out := new(ValidatedEntitiesResponse)
	err := c.cc.Invoke(ctx, "/proto.ValidatedActionService/Route", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ValidatedActionServiceServer is the server API for ValidatedActionService service.
// All implementations must embed UnimplementedValidatedActionServiceServer
// for forward compatibility
type ValidatedActionServiceServer interface {
	// Route sends the parsed entities from client to validator
	Route(context.Context, *ParsedEntitiesRequest) (*ValidatedEntitiesResponse, error)
	mustEmbedUnimplementedValidatedActionServiceServer()
}

// UnimplementedValidatedActionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedValidatedActionServiceServer struct {
}

func (UnimplementedValidatedActionServiceServer) Route(context.Context, *ParsedEntitiesRequest) (*ValidatedEntitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Route not implemented")
}
func (UnimplementedValidatedActionServiceServer) mustEmbedUnimplementedValidatedActionServiceServer() {
}

// UnsafeValidatedActionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ValidatedActionServiceServer will
// result in compilation errors.
type UnsafeValidatedActionServiceServer interface {
	mustEmbedUnimplementedValidatedActionServiceServer()
}

func RegisterValidatedActionServiceServer(s grpc.ServiceRegistrar, srv ValidatedActionServiceServer) {
	s.RegisterService(&ValidatedActionService_ServiceDesc, srv)
}

func _ValidatedActionService_Route_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParsedEntitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidatedActionServiceServer).Route(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ValidatedActionService/Route",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidatedActionServiceServer).Route(ctx, req.(*ParsedEntitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ValidatedActionService_ServiceDesc is the grpc.ServiceDesc for ValidatedActionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ValidatedActionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ValidatedActionService",
	HandlerType: (*ValidatedActionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Route",
			Handler:    _ValidatedActionService_Route_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entity.proto",
}

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	// GetDefinition gets the existing task's definition
	GetDefinition(ctx context.Context, in *TaskDefinitionRequest, opts ...grpc.CallOption) (*TaskDefinitionResponse, error)
	// GetStatus gets the current status with last start time and last end time of task
	GetStatus(ctx context.Context, in *TaskLatestStatusRequest, opts ...grpc.CallOption) (*TaskLatestStatusResponse, error)
	// GetNRuns gets the last N runs of task
	GetNRuns(ctx context.Context, in *TaskRunsRequest, opts ...grpc.CallOption) (*TaskRunsResponse, error)
	// GetRunsOn gets the runs on given date of task
	GetRunsOn(ctx context.Context, in *TaskRunsOnRequest, opts ...grpc.CallOption) (*TaskRunsOnResponse, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) GetDefinition(ctx context.Context, in *TaskDefinitionRequest, opts ...grpc.CallOption) (*TaskDefinitionResponse, error) {
	out := new(TaskDefinitionResponse)
	err := c.cc.Invoke(ctx, "/proto.TaskService/GetDefinition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) GetStatus(ctx context.Context, in *TaskLatestStatusRequest, opts ...grpc.CallOption) (*TaskLatestStatusResponse, error) {
	out := new(TaskLatestStatusResponse)
	err := c.cc.Invoke(ctx, "/proto.TaskService/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) GetNRuns(ctx context.Context, in *TaskRunsRequest, opts ...grpc.CallOption) (*TaskRunsResponse, error) {
	out := new(TaskRunsResponse)
	err := c.cc.Invoke(ctx, "/proto.TaskService/GetNRuns", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) GetRunsOn(ctx context.Context, in *TaskRunsOnRequest, opts ...grpc.CallOption) (*TaskRunsOnResponse, error) {
	out := new(TaskRunsOnResponse)
	err := c.cc.Invoke(ctx, "/proto.TaskService/GetRunsOn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations must embed UnimplementedTaskServiceServer
// for forward compatibility
type TaskServiceServer interface {
	// GetDefinition gets the existing task's definition
	GetDefinition(context.Context, *TaskDefinitionRequest) (*TaskDefinitionResponse, error)
	// GetStatus gets the current status with last start time and last end time of task
	GetStatus(context.Context, *TaskLatestStatusRequest) (*TaskLatestStatusResponse, error)
	// GetNRuns gets the last N runs of task
	GetNRuns(context.Context, *TaskRunsRequest) (*TaskRunsResponse, error)
	// GetRunsOn gets the runs on given date of task
	GetRunsOn(context.Context, *TaskRunsOnRequest) (*TaskRunsOnResponse, error)
	mustEmbedUnimplementedTaskServiceServer()
}

// UnimplementedTaskServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (UnimplementedTaskServiceServer) GetDefinition(context.Context, *TaskDefinitionRequest) (*TaskDefinitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefinition not implemented")
}
func (UnimplementedTaskServiceServer) GetStatus(context.Context, *TaskLatestStatusRequest) (*TaskLatestStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedTaskServiceServer) GetNRuns(context.Context, *TaskRunsRequest) (*TaskRunsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNRuns not implemented")
}
func (UnimplementedTaskServiceServer) GetRunsOn(context.Context, *TaskRunsOnRequest) (*TaskRunsOnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRunsOn not implemented")
}
func (UnimplementedTaskServiceServer) mustEmbedUnimplementedTaskServiceServer() {}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_GetDefinition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskDefinitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetDefinition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskService/GetDefinition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetDefinition(ctx, req.(*TaskDefinitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskLatestStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskService/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetStatus(ctx, req.(*TaskLatestStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_GetNRuns_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRunsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetNRuns(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskService/GetNRuns",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetNRuns(ctx, req.(*TaskRunsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_GetRunsOn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRunsOnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).GetRunsOn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskService/GetRunsOn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).GetRunsOn(ctx, req.(*TaskRunsOnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDefinition",
			Handler:    _TaskService_GetDefinition_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _TaskService_GetStatus_Handler,
		},
		{
			MethodName: "GetNRuns",
			Handler:    _TaskService_GetNRuns_Handler,
		},
		{
			MethodName: "GetRunsOn",
			Handler:    _TaskService_GetRunsOn_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entity.proto",
}
