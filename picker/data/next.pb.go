// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: next.proto

package data

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Message to represent job ready for next run
type ReadyJob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName       string   `protobuf:"bytes,1,opt,name=JobName,proto3" json:"JobName,omitempty"`
	ConditionJobs []string `protobuf:"bytes,2,rep,name=ConditionJobs,proto3" json:"ConditionJobs,omitempty"`
}

func (x *ReadyJob) Reset() {
	*x = ReadyJob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_next_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadyJob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadyJob) ProtoMessage() {}

func (x *ReadyJob) ProtoReflect() protoreflect.Message {
	mi := &file_next_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadyJob.ProtoReflect.Descriptor instead.
func (*ReadyJob) Descriptor() ([]byte, []int) {
	return file_next_proto_rawDescGZIP(), []int{0}
}

func (x *ReadyJob) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *ReadyJob) GetConditionJobs() []string {
	if x != nil {
		return x.ConditionJobs
	}
	return nil
}

// Request message to get list of jobs for next run
type NextJobsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NextJobsReq) Reset() {
	*x = NextJobsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_next_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextJobsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextJobsReq) ProtoMessage() {}

func (x *NextJobsReq) ProtoReflect() protoreflect.Message {
	mi := &file_next_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextJobsReq.ProtoReflect.Descriptor instead.
func (*NextJobsReq) Descriptor() ([]byte, []int) {
	return file_next_proto_rawDescGZIP(), []int{1}
}

// Response message with list of jobs
type NextJobsRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobList []*ReadyJob `protobuf:"bytes,1,rep,name=JobList,proto3" json:"JobList,omitempty"`
}

func (x *NextJobsRes) Reset() {
	*x = NextJobsRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_next_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextJobsRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextJobsRes) ProtoMessage() {}

func (x *NextJobsRes) ProtoReflect() protoreflect.Message {
	mi := &file_next_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextJobsRes.ProtoReflect.Descriptor instead.
func (*NextJobsRes) Descriptor() ([]byte, []int) {
	return file_next_proto_rawDescGZIP(), []int{2}
}

func (x *NextJobsRes) GetJobList() []*ReadyJob {
	if x != nil {
		return x.JobList
	}
	return nil
}

var File_next_proto protoreflect.FileDescriptor

var file_next_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x4a, 0x0a, 0x08, 0x52, 0x65, 0x61, 0x64, 0x79, 0x4a, 0x6f, 0x62, 0x12, 0x18,
	0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x4a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0d, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x73, 0x22, 0x0d,
	0x0a, 0x0b, 0x4e, 0x65, 0x78, 0x74, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x22, 0x37, 0x0a,
	0x0b, 0x4e, 0x65, 0x78, 0x74, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x07,
	0x4a, 0x6f, 0x62, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x4a, 0x6f, 0x62, 0x52, 0x07, 0x4a,
	0x6f, 0x62, 0x4c, 0x69, 0x73, 0x74, 0x32, 0x3a, 0x0a, 0x08, 0x4e, 0x65, 0x78, 0x74, 0x4a, 0x6f,
	0x62, 0x73, 0x12, 0x2e, 0x0a, 0x04, 0x4e, 0x65, 0x78, 0x74, 0x12, 0x11, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x73,
	0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x3b, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_next_proto_rawDescOnce sync.Once
	file_next_proto_rawDescData = file_next_proto_rawDesc
)

func file_next_proto_rawDescGZIP() []byte {
	file_next_proto_rawDescOnce.Do(func() {
		file_next_proto_rawDescData = protoimpl.X.CompressGZIP(file_next_proto_rawDescData)
	})
	return file_next_proto_rawDescData
}

var file_next_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_next_proto_goTypes = []interface{}{
	(*ReadyJob)(nil),    // 0: data.ReadyJob
	(*NextJobsReq)(nil), // 1: data.NextJobsReq
	(*NextJobsRes)(nil), // 2: data.NextJobsRes
}
var file_next_proto_depIdxs = []int32{
	0, // 0: data.NextJobsRes.JobList:type_name -> data.ReadyJob
	1, // 1: data.NextJobs.Next:input_type -> data.NextJobsReq
	2, // 2: data.NextJobs.Next:output_type -> data.NextJobsRes
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_next_proto_init() }
func file_next_proto_init() {
	if File_next_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_next_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadyJob); i {
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
		file_next_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextJobsReq); i {
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
		file_next_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextJobsRes); i {
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
			RawDescriptor: file_next_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_next_proto_goTypes,
		DependencyIndexes: file_next_proto_depIdxs,
		MessageInfos:      file_next_proto_msgTypes,
	}.Build()
	File_next_proto = out.File
	file_next_proto_rawDesc = nil
	file_next_proto_goTypes = nil
	file_next_proto_depIdxs = nil
}
