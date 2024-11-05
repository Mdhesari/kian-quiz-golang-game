// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v4.25.1
// source: protobuf/presence/presence.proto

package presence

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

type GetPresenceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId []string `protobuf:"bytes,1,rep,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetPresenceRequest) Reset() {
	*x = GetPresenceRequest{}
	mi := &file_protobuf_presence_presence_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPresenceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPresenceRequest) ProtoMessage() {}

func (x *GetPresenceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_presence_presence_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPresenceRequest.ProtoReflect.Descriptor instead.
func (*GetPresenceRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_presence_presence_proto_rawDescGZIP(), []int{0}
}

func (x *GetPresenceRequest) GetUserId() []string {
	if x != nil {
		return x.UserId
	}
	return nil
}

type GetPresenceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*GetPresenceItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GetPresenceResponse) Reset() {
	*x = GetPresenceResponse{}
	mi := &file_protobuf_presence_presence_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPresenceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPresenceResponse) ProtoMessage() {}

func (x *GetPresenceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_presence_presence_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPresenceResponse.ProtoReflect.Descriptor instead.
func (*GetPresenceResponse) Descriptor() ([]byte, []int) {
	return file_protobuf_presence_presence_proto_rawDescGZIP(), []int{1}
}

func (x *GetPresenceResponse) GetItems() []*GetPresenceItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetPresenceItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Timestamp uint64 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *GetPresenceItem) Reset() {
	*x = GetPresenceItem{}
	mi := &file_protobuf_presence_presence_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPresenceItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPresenceItem) ProtoMessage() {}

func (x *GetPresenceItem) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_presence_presence_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPresenceItem.ProtoReflect.Descriptor instead.
func (*GetPresenceItem) Descriptor() ([]byte, []int) {
	return file_protobuf_presence_presence_proto_rawDescGZIP(), []int{2}
}

func (x *GetPresenceItem) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetPresenceItem) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_protobuf_presence_presence_proto protoreflect.FileDescriptor

var file_protobuf_presence_presence_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x65, 0x73, 0x65,
	0x6e, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x2c, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x46, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2f, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x50,
	0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x22, 0x47, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x32, 0x5d, 0x0a, 0x0f, 0x50,
	0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4a,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x2e,
	0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73,
	0x65, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72,
	0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1a, 0x5a, 0x18, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x70, 0x72,
	0x65, 0x73, 0x65, 0x6e, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_presence_presence_proto_rawDescOnce sync.Once
	file_protobuf_presence_presence_proto_rawDescData = file_protobuf_presence_presence_proto_rawDesc
)

func file_protobuf_presence_presence_proto_rawDescGZIP() []byte {
	file_protobuf_presence_presence_proto_rawDescOnce.Do(func() {
		file_protobuf_presence_presence_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_presence_presence_proto_rawDescData)
	})
	return file_protobuf_presence_presence_proto_rawDescData
}

var file_protobuf_presence_presence_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_protobuf_presence_presence_proto_goTypes = []any{
	(*GetPresenceRequest)(nil),  // 0: presence.GetPresenceRequest
	(*GetPresenceResponse)(nil), // 1: presence.GetPresenceResponse
	(*GetPresenceItem)(nil),     // 2: presence.GetPresenceItem
}
var file_protobuf_presence_presence_proto_depIdxs = []int32{
	2, // 0: presence.GetPresenceResponse.items:type_name -> presence.GetPresenceItem
	0, // 1: presence.PresenceService.GetPresence:input_type -> presence.GetPresenceRequest
	1, // 2: presence.PresenceService.GetPresence:output_type -> presence.GetPresenceResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protobuf_presence_presence_proto_init() }
func file_protobuf_presence_presence_proto_init() {
	if File_protobuf_presence_presence_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protobuf_presence_presence_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_presence_presence_proto_goTypes,
		DependencyIndexes: file_protobuf_presence_presence_proto_depIdxs,
		MessageInfos:      file_protobuf_presence_presence_proto_msgTypes,
	}.Build()
	File_protobuf_presence_presence_proto = out.File
	file_protobuf_presence_presence_proto_rawDesc = nil
	file_protobuf_presence_presence_proto_goTypes = nil
	file_protobuf_presence_presence_proto_depIdxs = nil
}
