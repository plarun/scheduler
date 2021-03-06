// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: depend.proto

package data

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Message to represent job with its current status
type JobWithStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName    string `protobuf:"bytes,1,opt,name=JobName,proto3" json:"JobName,omitempty"`
	StatusType Status `protobuf:"varint,2,opt,name=StatusType,proto3,enum=data.Status" json:"StatusType,omitempty"`
}

func (x *JobWithStatus) Reset() {
	*x = JobWithStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depend_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobWithStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobWithStatus) ProtoMessage() {}

func (x *JobWithStatus) ProtoReflect() protoreflect.Message {
	mi := &file_depend_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobWithStatus.ProtoReflect.Descriptor instead.
func (*JobWithStatus) Descriptor() ([]byte, []int) {
	return file_depend_proto_rawDescGZIP(), []int{0}
}

func (x *JobWithStatus) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *JobWithStatus) GetStatusType() Status {
	if x != nil {
		return x.StatusType
	}
	return Status_IDLE
}

// Request message to get preceded and succeeded jobs of job
type GetJobDependsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName string `protobuf:"bytes,1,opt,name=JobName,proto3" json:"JobName,omitempty"`
}

func (x *GetJobDependsReq) Reset() {
	*x = GetJobDependsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depend_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetJobDependsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJobDependsReq) ProtoMessage() {}

func (x *GetJobDependsReq) ProtoReflect() protoreflect.Message {
	mi := &file_depend_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJobDependsReq.ProtoReflect.Descriptor instead.
func (*GetJobDependsReq) Descriptor() ([]byte, []int) {
	return file_depend_proto_rawDescGZIP(), []int{1}
}

func (x *GetJobDependsReq) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

// Response message of preceded and succeeded jobs of job
type GetJobDependsRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName       string           `protobuf:"bytes,1,opt,name=JobName,proto3" json:"JobName,omitempty"`
	PrecededJobs  []*JobWithStatus `protobuf:"bytes,2,rep,name=PrecededJobs,proto3" json:"PrecededJobs,omitempty"`
	SucceededJobs []*JobWithStatus `protobuf:"bytes,3,rep,name=SucceededJobs,proto3" json:"SucceededJobs,omitempty"`
}

func (x *GetJobDependsRes) Reset() {
	*x = GetJobDependsRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depend_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetJobDependsRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJobDependsRes) ProtoMessage() {}

func (x *GetJobDependsRes) ProtoReflect() protoreflect.Message {
	mi := &file_depend_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJobDependsRes.ProtoReflect.Descriptor instead.
func (*GetJobDependsRes) Descriptor() ([]byte, []int) {
	return file_depend_proto_rawDescGZIP(), []int{2}
}

func (x *GetJobDependsRes) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *GetJobDependsRes) GetPrecededJobs() []*JobWithStatus {
	if x != nil {
		return x.PrecededJobs
	}
	return nil
}

func (x *GetJobDependsRes) GetSucceededJobs() []*JobWithStatus {
	if x != nil {
		return x.SucceededJobs
	}
	return nil
}

// Request message to get next runtime of job
type GetNextRunReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName string `protobuf:"bytes,1,opt,name=JobName,proto3" json:"JobName,omitempty"`
}

func (x *GetNextRunReq) Reset() {
	*x = GetNextRunReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depend_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNextRunReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNextRunReq) ProtoMessage() {}

func (x *GetNextRunReq) ProtoReflect() protoreflect.Message {
	mi := &file_depend_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNextRunReq.ProtoReflect.Descriptor instead.
func (*GetNextRunReq) Descriptor() ([]byte, []int) {
	return file_depend_proto_rawDescGZIP(), []int{3}
}

func (x *GetNextRunReq) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

// Response message of next runtime of job
type GetNextRunRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName string                 `protobuf:"bytes,1,opt,name=JobName,proto3" json:"JobName,omitempty"`
	NextRun *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=NextRun,proto3" json:"NextRun,omitempty"`
}

func (x *GetNextRunRes) Reset() {
	*x = GetNextRunRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_depend_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNextRunRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNextRunRes) ProtoMessage() {}

func (x *GetNextRunRes) ProtoReflect() protoreflect.Message {
	mi := &file_depend_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNextRunRes.ProtoReflect.Descriptor instead.
func (*GetNextRunRes) Descriptor() ([]byte, []int) {
	return file_depend_proto_rawDescGZIP(), []int{4}
}

func (x *GetNextRunRes) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *GetNextRunRes) GetNextRun() *timestamppb.Timestamp {
	if x != nil {
		return x.NextRun
	}
	return nil
}

var File_depend_proto protoreflect.FileDescriptor

var file_depend_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x57, 0x0a, 0x0d, 0x4a, 0x6f, 0x62, 0x57, 0x69, 0x74, 0x68, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2c,
	0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2c, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x71,
	0x12, 0x18, 0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xa0, 0x01, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x4a, 0x6f, 0x62, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x0c, 0x50, 0x72, 0x65,
	0x63, 0x65, 0x64, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4a, 0x6f, 0x62, 0x57, 0x69, 0x74, 0x68, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x0c, 0x50, 0x72, 0x65, 0x63, 0x65, 0x64, 0x65, 0x64, 0x4a, 0x6f,
	0x62, 0x73, 0x12, 0x39, 0x0a, 0x0d, 0x53, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x4a,
	0x6f, 0x62, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x4a, 0x6f, 0x62, 0x57, 0x69, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0d,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x4a, 0x6f, 0x62, 0x73, 0x22, 0x29, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x18,
	0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x5f, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4e,
	0x65, 0x78, 0x74, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x4a, 0x6f, 0x62,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4a, 0x6f, 0x62, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x75, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x07, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x75, 0x6e, 0x32, 0x89, 0x01, 0x0a, 0x0a, 0x4a, 0x6f,
	0x62, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x12, 0x41, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4a,
	0x6f, 0x62, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x12, 0x16, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x52, 0x65,
	0x71, 0x1a, 0x16, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x44,
	0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x73, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x75, 0x6e, 0x12, 0x13, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x47, 0x65, 0x74, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x13,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x75, 0x6e,
	0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x3b, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_depend_proto_rawDescOnce sync.Once
	file_depend_proto_rawDescData = file_depend_proto_rawDesc
)

func file_depend_proto_rawDescGZIP() []byte {
	file_depend_proto_rawDescOnce.Do(func() {
		file_depend_proto_rawDescData = protoimpl.X.CompressGZIP(file_depend_proto_rawDescData)
	})
	return file_depend_proto_rawDescData
}

var file_depend_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_depend_proto_goTypes = []interface{}{
	(*JobWithStatus)(nil),         // 0: data.JobWithStatus
	(*GetJobDependsReq)(nil),      // 1: data.GetJobDependsReq
	(*GetJobDependsRes)(nil),      // 2: data.GetJobDependsRes
	(*GetNextRunReq)(nil),         // 3: data.GetNextRunReq
	(*GetNextRunRes)(nil),         // 4: data.GetNextRunRes
	(Status)(0),                   // 5: data.Status
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_depend_proto_depIdxs = []int32{
	5, // 0: data.JobWithStatus.StatusType:type_name -> data.Status
	0, // 1: data.GetJobDependsRes.PrecededJobs:type_name -> data.JobWithStatus
	0, // 2: data.GetJobDependsRes.SucceededJobs:type_name -> data.JobWithStatus
	6, // 3: data.GetNextRunRes.NextRun:type_name -> google.protobuf.Timestamp
	1, // 4: data.JobDepends.GetJobDepends:input_type -> data.GetJobDependsReq
	3, // 5: data.JobDepends.GetNextRun:input_type -> data.GetNextRunReq
	2, // 6: data.JobDepends.GetJobDepends:output_type -> data.GetJobDependsRes
	4, // 7: data.JobDepends.GetNextRun:output_type -> data.GetNextRunRes
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_depend_proto_init() }
func file_depend_proto_init() {
	if File_depend_proto != nil {
		return
	}
	file_status_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_depend_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobWithStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depend_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetJobDependsReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depend_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetJobDependsRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depend_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNextRunReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_depend_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNextRunRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_depend_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_depend_proto_goTypes,
		DependencyIndexes: file_depend_proto_depIdxs,
		MessageInfos:      file_depend_proto_msgTypes,
	}.Build()
	File_depend_proto = out.File
	file_depend_proto_rawDesc = nil
	file_depend_proto_goTypes = nil
	file_depend_proto_depIdxs = nil
}
