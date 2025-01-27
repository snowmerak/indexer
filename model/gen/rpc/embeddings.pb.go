// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: rpc/embeddings.proto

package rpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetEmbeddingsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Model         string                 `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
	Contents      string                 `protobuf:"bytes,2,opt,name=contents,proto3" json:"contents,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetEmbeddingsRequest) Reset() {
	*x = GetEmbeddingsRequest{}
	mi := &file_rpc_embeddings_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetEmbeddingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEmbeddingsRequest) ProtoMessage() {}

func (x *GetEmbeddingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_embeddings_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEmbeddingsRequest.ProtoReflect.Descriptor instead.
func (*GetEmbeddingsRequest) Descriptor() ([]byte, []int) {
	return file_rpc_embeddings_proto_rawDescGZIP(), []int{0}
}

func (x *GetEmbeddingsRequest) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *GetEmbeddingsRequest) GetContents() string {
	if x != nil {
		return x.Contents
	}
	return ""
}

type GetEmbeddingsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Embeddings    []float32              `protobuf:"fixed32,1,rep,packed,name=embeddings,proto3" json:"embeddings,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetEmbeddingsResponse) Reset() {
	*x = GetEmbeddingsResponse{}
	mi := &file_rpc_embeddings_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetEmbeddingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEmbeddingsResponse) ProtoMessage() {}

func (x *GetEmbeddingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_embeddings_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEmbeddingsResponse.ProtoReflect.Descriptor instead.
func (*GetEmbeddingsResponse) Descriptor() ([]byte, []int) {
	return file_rpc_embeddings_proto_rawDescGZIP(), []int{1}
}

func (x *GetEmbeddingsResponse) GetEmbeddings() []float32 {
	if x != nil {
		return x.Embeddings
	}
	return nil
}

var File_rpc_embeddings_proto protoreflect.FileDescriptor

var file_rpc_embeddings_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73,
	0x22, 0x37, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x02, 0x52, 0x0a, 0x65,
	0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x32, 0x55, 0x0a, 0x11, 0x45, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40,
	0x0a, 0x0d, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x12,
	0x15, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6d, 0x62, 0x65,
	0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x3f, 0x42, 0x0f, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6e, 0x6f, 0x77, 0x6d, 0x65, 0x72, 0x61, 0x6b, 0x2f, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x65, 0x72, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x72, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_rpc_embeddings_proto_rawDescOnce sync.Once
	file_rpc_embeddings_proto_rawDescData []byte
)

func file_rpc_embeddings_proto_rawDescGZIP() []byte {
	file_rpc_embeddings_proto_rawDescOnce.Do(func() {
		file_rpc_embeddings_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_rpc_embeddings_proto_rawDesc), len(file_rpc_embeddings_proto_rawDesc)))
	})
	return file_rpc_embeddings_proto_rawDescData
}

var file_rpc_embeddings_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_embeddings_proto_goTypes = []any{
	(*GetEmbeddingsRequest)(nil),  // 0: GetEmbeddingsRequest
	(*GetEmbeddingsResponse)(nil), // 1: GetEmbeddingsResponse
}
var file_rpc_embeddings_proto_depIdxs = []int32{
	0, // 0: EmbeddingsService.GetEmbeddings:input_type -> GetEmbeddingsRequest
	1, // 1: EmbeddingsService.GetEmbeddings:output_type -> GetEmbeddingsResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_embeddings_proto_init() }
func file_rpc_embeddings_proto_init() {
	if File_rpc_embeddings_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_rpc_embeddings_proto_rawDesc), len(file_rpc_embeddings_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_embeddings_proto_goTypes,
		DependencyIndexes: file_rpc_embeddings_proto_depIdxs,
		MessageInfos:      file_rpc_embeddings_proto_msgTypes,
	}.Build()
	File_rpc_embeddings_proto = out.File
	file_rpc_embeddings_proto_goTypes = nil
	file_rpc_embeddings_proto_depIdxs = nil
}
