// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        v4.22.2
// source: base.proto

package dove

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

type DoveMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CrcId     uint64 `protobuf:"varint,1,opt,name=crcId,proto3" json:"crcId,omitempty"`         //请求ID
	AckId     uint64 `protobuf:"varint,2,opt,name=ackId,proto3" json:"ackId,omitempty"`         //回执ID
	Seq       string `protobuf:"bytes,3,opt,name=seq,proto3" json:"seq,omitempty"`              //随机数
	Timestamp int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"` //时间
}

func (x *DoveMetadata) Reset() {
	*x = DoveMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoveMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoveMetadata) ProtoMessage() {}

func (x *DoveMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoveMetadata.ProtoReflect.Descriptor instead.
func (*DoveMetadata) Descriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{0}
}

func (x *DoveMetadata) GetCrcId() uint64 {
	if x != nil {
		return x.CrcId
	}
	return 0
}

func (x *DoveMetadata) GetAckId() uint64 {
	if x != nil {
		return x.AckId
	}
	return 0
}

func (x *DoveMetadata) GetSeq() string {
	if x != nil {
		return x.Seq
	}
	return ""
}

func (x *DoveMetadata) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type DoveBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`  // 数据
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`    // 描述
	Code uint64 `protobuf:"varint,3,opt,name=code,proto3" json:"code,omitempty"` // 状态码
}

func (x *DoveBody) Reset() {
	*x = DoveBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_base_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoveBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoveBody) ProtoMessage() {}

func (x *DoveBody) ProtoReflect() protoreflect.Message {
	mi := &file_base_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoveBody.ProtoReflect.Descriptor instead.
func (*DoveBody) Descriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{1}
}

func (x *DoveBody) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DoveBody) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *DoveBody) GetCode() uint64 {
	if x != nil {
		return x.Code
	}
	return 0
}

type Dove struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *DoveMetadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Body     *DoveBody     `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Dove) Reset() {
	*x = Dove{}
	if protoimpl.UnsafeEnabled {
		mi := &file_base_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dove) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dove) ProtoMessage() {}

func (x *Dove) ProtoReflect() protoreflect.Message {
	mi := &file_base_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dove.ProtoReflect.Descriptor instead.
func (*Dove) Descriptor() ([]byte, []int) {
	return file_base_proto_rawDescGZIP(), []int{2}
}

func (x *Dove) GetMetadata() *DoveMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Dove) GetBody() *DoveBody {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_base_proto protoreflect.FileDescriptor

var file_base_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x6f,
	0x76, 0x65, 0x22, 0x6a, 0x0a, 0x0c, 0x44, 0x6f, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x72, 0x63, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x63, 0x72, 0x63, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x63, 0x6b, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x10,
	0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x65, 0x71,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x44,
	0x0a, 0x08, 0x44, 0x6f, 0x76, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10,
	0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x22, 0x5a, 0x0a, 0x04, 0x44, 0x6f, 0x76, 0x65, 0x12, 0x2e, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x64, 0x6f, 0x76, 0x65, 0x2e, 0x44, 0x6f, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x22, 0x0a, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x64, 0x6f, 0x76,
	0x65, 0x2e, 0x44, 0x6f, 0x76, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x64, 0x6f, 0x76, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_base_proto_rawDescOnce sync.Once
	file_base_proto_rawDescData = file_base_proto_rawDesc
)

func file_base_proto_rawDescGZIP() []byte {
	file_base_proto_rawDescOnce.Do(func() {
		file_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_base_proto_rawDescData)
	})
	return file_base_proto_rawDescData
}

var file_base_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_base_proto_goTypes = []interface{}{
	(*DoveMetadata)(nil), // 0: dove.DoveMetadata
	(*DoveBody)(nil),     // 1: dove.DoveBody
	(*Dove)(nil),         // 2: dove.Dove
}
var file_base_proto_depIdxs = []int32{
	0, // 0: dove.Dove.metadata:type_name -> dove.DoveMetadata
	1, // 1: dove.Dove.body:type_name -> dove.DoveBody
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_base_proto_init() }
func file_base_proto_init() {
	if File_base_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoveMetadata); i {
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
		file_base_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoveBody); i {
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
		file_base_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Dove); i {
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
			RawDescriptor: file_base_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_base_proto_goTypes,
		DependencyIndexes: file_base_proto_depIdxs,
		MessageInfos:      file_base_proto_msgTypes,
	}.Build()
	File_base_proto = out.File
	file_base_proto_rawDesc = nil
	file_base_proto_goTypes = nil
	file_base_proto_depIdxs = nil
}
