// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: proto/texc.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data        []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Exec        []string `protobuf:"bytes,2,rep,name=exec,proto3" json:"exec,omitempty"`
	Dl          string   `protobuf:"bytes,3,opt,name=dl,proto3" json:"dl,omitempty"`
	NoOutStream bool     `protobuf:"varint,4,opt,name=no_out_stream,json=noOutStream,proto3" json:"no_out_stream,omitempty"`
}

func (x *Input) Reset() {
	*x = Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_texc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Input) ProtoMessage() {}

func (x *Input) ProtoReflect() protoreflect.Message {
	mi := &file_proto_texc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Input.ProtoReflect.Descriptor instead.
func (*Input) Descriptor() ([]byte, []int) {
	return file_proto_texc_proto_rawDescGZIP(), []int{0}
}

func (x *Input) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Input) GetExec() []string {
	if x != nil {
		return x.Exec
	}
	return nil
}

func (x *Input) GetDl() string {
	if x != nil {
		return x.Dl
	}
	return ""
}

func (x *Input) GetNoOutStream() bool {
	if x != nil {
		return x.NoOutStream
	}
	return false
}

type Output struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data   []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Stdout []byte `protobuf:"bytes,2,opt,name=stdout,proto3" json:"stdout,omitempty"`
	Stderr []byte `protobuf:"bytes,3,opt,name=stderr,proto3" json:"stderr,omitempty"`
}

func (x *Output) Reset() {
	*x = Output{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_texc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Output) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Output) ProtoMessage() {}

func (x *Output) ProtoReflect() protoreflect.Message {
	mi := &file_proto_texc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Output.ProtoReflect.Descriptor instead.
func (*Output) Descriptor() ([]byte, []int) {
	return file_proto_texc_proto_rawDescGZIP(), []int{1}
}

func (x *Output) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Output) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

func (x *Output) GetStderr() []byte {
	if x != nil {
		return x.Stderr
	}
	return nil
}

var File_proto_texc_proto protoreflect.FileDescriptor

var file_proto_texc_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x65, 0x78, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63, 0x0a, 0x05, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x78, 0x65, 0x63, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x65, 0x78, 0x65, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x64, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x64, 0x6c, 0x12, 0x22, 0x0a, 0x0d, 0x6e, 0x6f,
	0x5f, 0x6f, 0x75, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0b, 0x6e, 0x6f, 0x4f, 0x75, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x4c,
	0x0a, 0x06, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74,
	0x64, 0x6f, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x32, 0x38, 0x0a, 0x0b,
	0x54, 0x65, 0x78, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x53,
	0x79, 0x6e, 0x63, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_texc_proto_rawDescOnce sync.Once
	file_proto_texc_proto_rawDescData = file_proto_texc_proto_rawDesc
)

func file_proto_texc_proto_rawDescGZIP() []byte {
	file_proto_texc_proto_rawDescOnce.Do(func() {
		file_proto_texc_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_texc_proto_rawDescData)
	})
	return file_proto_texc_proto_rawDescData
}

var file_proto_texc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_texc_proto_goTypes = []interface{}{
	(*Input)(nil),  // 0: proto.Input
	(*Output)(nil), // 1: proto.Output
}
var file_proto_texc_proto_depIdxs = []int32{
	0, // 0: proto.TexcService.Sync:input_type -> proto.Input
	1, // 1: proto.TexcService.Sync:output_type -> proto.Output
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_texc_proto_init() }
func file_proto_texc_proto_init() {
	if File_proto_texc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_texc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Input); i {
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
		file_proto_texc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Output); i {
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
			RawDescriptor: file_proto_texc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_texc_proto_goTypes,
		DependencyIndexes: file_proto_texc_proto_depIdxs,
		MessageInfos:      file_proto_texc_proto_msgTypes,
	}.Build()
	File_proto_texc_proto = out.File
	file_proto_texc_proto_rawDesc = nil
	file_proto_texc_proto_goTypes = nil
	file_proto_texc_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TexcServiceClient is the client API for TexcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TexcServiceClient interface {
	Sync(ctx context.Context, opts ...grpc.CallOption) (TexcService_SyncClient, error)
}

type texcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTexcServiceClient(cc grpc.ClientConnInterface) TexcServiceClient {
	return &texcServiceClient{cc}
}

func (c *texcServiceClient) Sync(ctx context.Context, opts ...grpc.CallOption) (TexcService_SyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TexcService_serviceDesc.Streams[0], "/proto.TexcService/Sync", opts...)
	if err != nil {
		return nil, err
	}
	x := &texcServiceSyncClient{stream}
	return x, nil
}

type TexcService_SyncClient interface {
	Send(*Input) error
	Recv() (*Output, error)
	grpc.ClientStream
}

type texcServiceSyncClient struct {
	grpc.ClientStream
}

func (x *texcServiceSyncClient) Send(m *Input) error {
	return x.ClientStream.SendMsg(m)
}

func (x *texcServiceSyncClient) Recv() (*Output, error) {
	m := new(Output)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TexcServiceServer is the server API for TexcService service.
type TexcServiceServer interface {
	Sync(TexcService_SyncServer) error
}

// UnimplementedTexcServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTexcServiceServer struct {
}

func (*UnimplementedTexcServiceServer) Sync(TexcService_SyncServer) error {
	return status.Errorf(codes.Unimplemented, "method Sync not implemented")
}

func RegisterTexcServiceServer(s *grpc.Server, srv TexcServiceServer) {
	s.RegisterService(&_TexcService_serviceDesc, srv)
}

func _TexcService_Sync_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TexcServiceServer).Sync(&texcServiceSyncServer{stream})
}

type TexcService_SyncServer interface {
	Send(*Output) error
	Recv() (*Input, error)
	grpc.ServerStream
}

type texcServiceSyncServer struct {
	grpc.ServerStream
}

func (x *texcServiceSyncServer) Send(m *Output) error {
	return x.ServerStream.SendMsg(m)
}

func (x *texcServiceSyncServer) Recv() (*Input, error) {
	m := new(Input)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _TexcService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TexcService",
	HandlerType: (*TexcServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Sync",
			Handler:       _TexcService_Sync_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/texc.proto",
}
