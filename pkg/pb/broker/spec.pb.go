// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0-devel
// 	protoc        v3.11.4
// source: spec.proto

package broker

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

type ConsumeRequest_Action int32

const (
	ConsumeRequest_SUBSCRIBE   ConsumeRequest_Action = 0
	ConsumeRequest_UNSUBSCRIBE ConsumeRequest_Action = 1
)

// Enum value maps for ConsumeRequest_Action.
var (
	ConsumeRequest_Action_name = map[int32]string{
		0: "SUBSCRIBE",
		1: "UNSUBSCRIBE",
	}
	ConsumeRequest_Action_value = map[string]int32{
		"SUBSCRIBE":   0,
		"UNSUBSCRIBE": 1,
	}
)

func (x ConsumeRequest_Action) Enum() *ConsumeRequest_Action {
	p := new(ConsumeRequest_Action)
	*p = x
	return p
}

func (x ConsumeRequest_Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsumeRequest_Action) Descriptor() protoreflect.EnumDescriptor {
	return file_spec_proto_enumTypes[0].Descriptor()
}

func (ConsumeRequest_Action) Type() protoreflect.EnumType {
	return &file_spec_proto_enumTypes[0]
}

func (x ConsumeRequest_Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsumeRequest_Action.Descriptor instead.
func (ConsumeRequest_Action) EnumDescriptor() ([]byte, []int) {
	return file_spec_proto_rawDescGZIP(), []int{0, 0}
}

type ConsumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action ConsumeRequest_Action `protobuf:"varint,1,opt,name=action,proto3,enum=mbproto.ConsumeRequest_Action" json:"action,omitempty"`
	// Keys is a list of words divided by dot. May be specified as exact key or pattern.
	// Word should contain only 0-9a-zA-Z characters.
	// Instead of word may be specified:
	// * (star) can substitute for exactly one word.
	// # (hash) can substitute for zero or more words.
	Keys []string `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *ConsumeRequest) Reset() {
	*x = ConsumeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spec_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeRequest) ProtoMessage() {}

func (x *ConsumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_spec_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeRequest.ProtoReflect.Descriptor instead.
func (*ConsumeRequest) Descriptor() ([]byte, []int) {
	return file_spec_proto_rawDescGZIP(), []int{0}
}

func (x *ConsumeRequest) GetAction() ConsumeRequest_Action {
	if x != nil {
		return x.Action
	}
	return ConsumeRequest_SUBSCRIBE
}

func (x *ConsumeRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type ConsumeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key is a list of words divided by dot. Word should contain only 0-9a-zA-Z characters.
	// For example: aaa, aaa.bbb, ccc.123.ddd, etc.
	Key     string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ConsumeResponse) Reset() {
	*x = ConsumeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spec_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeResponse) ProtoMessage() {}

func (x *ConsumeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_spec_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeResponse.ProtoReflect.Descriptor instead.
func (*ConsumeResponse) Descriptor() ([]byte, []int) {
	return file_spec_proto_rawDescGZIP(), []int{1}
}

func (x *ConsumeResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ConsumeResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type ProduceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key is a list of words divided by dot. Word should contain only 0-9a-zA-Z characters.
	// For example: aaa, aaa.bbb, ccc.123.ddd, etc.
	Key     string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ProduceRequest) Reset() {
	*x = ProduceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spec_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProduceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProduceRequest) ProtoMessage() {}

func (x *ProduceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_spec_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProduceRequest.ProtoReflect.Descriptor instead.
func (*ProduceRequest) Descriptor() ([]byte, []int) {
	return file_spec_proto_rawDescGZIP(), []int{2}
}

func (x *ProduceRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ProduceRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type ProduceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProduceResponse) Reset() {
	*x = ProduceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spec_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProduceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProduceResponse) ProtoMessage() {}

func (x *ProduceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_spec_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProduceResponse.ProtoReflect.Descriptor instead.
func (*ProduceResponse) Descriptor() ([]byte, []int) {
	return file_spec_proto_rawDescGZIP(), []int{3}
}

var File_spec_proto protoreflect.FileDescriptor

var file_spec_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x62,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x86, 0x01, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x6d, 0x62, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x6b, 0x65, 0x79, 0x73, 0x22, 0x28, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0d,
	0x0a, 0x09, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x10, 0x00, 0x12, 0x0f, 0x0a,
	0x0b, 0x55, 0x4e, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x10, 0x01, 0x22, 0x3d,
	0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x3c, 0x0a,
	0x0e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x91,
	0x01, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x72, 0x6f, 0x6b, 0x65, 0x72,
	0x12, 0x3e, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x12, 0x17, 0x2e, 0x6d, 0x62,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6d, 0x62, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01,
	0x12, 0x40, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x12, 0x17, 0x2e, 0x6d, 0x62,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6d, 0x62, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43,
	0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01,
	0x30, 0x01, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_spec_proto_rawDescOnce sync.Once
	file_spec_proto_rawDescData = file_spec_proto_rawDesc
)

func file_spec_proto_rawDescGZIP() []byte {
	file_spec_proto_rawDescOnce.Do(func() {
		file_spec_proto_rawDescData = protoimpl.X.CompressGZIP(file_spec_proto_rawDescData)
	})
	return file_spec_proto_rawDescData
}

var file_spec_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_spec_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_spec_proto_goTypes = []interface{}{
	(ConsumeRequest_Action)(0), // 0: mbproto.ConsumeRequest.Action
	(*ConsumeRequest)(nil),     // 1: mbproto.ConsumeRequest
	(*ConsumeResponse)(nil),    // 2: mbproto.ConsumeResponse
	(*ProduceRequest)(nil),     // 3: mbproto.ProduceRequest
	(*ProduceResponse)(nil),    // 4: mbproto.ProduceResponse
}
var file_spec_proto_depIdxs = []int32{
	0, // 0: mbproto.ConsumeRequest.action:type_name -> mbproto.ConsumeRequest.Action
	3, // 1: mbproto.MessageBroker.Produce:input_type -> mbproto.ProduceRequest
	1, // 2: mbproto.MessageBroker.Consume:input_type -> mbproto.ConsumeRequest
	4, // 3: mbproto.MessageBroker.Produce:output_type -> mbproto.ProduceResponse
	2, // 4: mbproto.MessageBroker.Consume:output_type -> mbproto.ConsumeResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_spec_proto_init() }
func file_spec_proto_init() {
	if File_spec_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_spec_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeRequest); i {
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
		file_spec_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeResponse); i {
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
		file_spec_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProduceRequest); i {
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
		file_spec_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProduceResponse); i {
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
			RawDescriptor: file_spec_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_spec_proto_goTypes,
		DependencyIndexes: file_spec_proto_depIdxs,
		EnumInfos:         file_spec_proto_enumTypes,
		MessageInfos:      file_spec_proto_msgTypes,
	}.Build()
	File_spec_proto = out.File
	file_spec_proto_rawDesc = nil
	file_spec_proto_goTypes = nil
	file_spec_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MessageBrokerClient is the client API for MessageBroker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MessageBrokerClient interface {
	Produce(ctx context.Context, opts ...grpc.CallOption) (MessageBroker_ProduceClient, error)
	Consume(ctx context.Context, opts ...grpc.CallOption) (MessageBroker_ConsumeClient, error)
}

type messageBrokerClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageBrokerClient(cc grpc.ClientConnInterface) MessageBrokerClient {
	return &messageBrokerClient{cc}
}

func (c *messageBrokerClient) Produce(ctx context.Context, opts ...grpc.CallOption) (MessageBroker_ProduceClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MessageBroker_serviceDesc.Streams[0], "/mbproto.MessageBroker/Produce", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageBrokerProduceClient{stream}
	return x, nil
}

type MessageBroker_ProduceClient interface {
	Send(*ProduceRequest) error
	CloseAndRecv() (*ProduceResponse, error)
	grpc.ClientStream
}

type messageBrokerProduceClient struct {
	grpc.ClientStream
}

func (x *messageBrokerProduceClient) Send(m *ProduceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageBrokerProduceClient) CloseAndRecv() (*ProduceResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ProduceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *messageBrokerClient) Consume(ctx context.Context, opts ...grpc.CallOption) (MessageBroker_ConsumeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MessageBroker_serviceDesc.Streams[1], "/mbproto.MessageBroker/Consume", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageBrokerConsumeClient{stream}
	return x, nil
}

type MessageBroker_ConsumeClient interface {
	Send(*ConsumeRequest) error
	Recv() (*ConsumeResponse, error)
	grpc.ClientStream
}

type messageBrokerConsumeClient struct {
	grpc.ClientStream
}

func (x *messageBrokerConsumeClient) Send(m *ConsumeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageBrokerConsumeClient) Recv() (*ConsumeResponse, error) {
	m := new(ConsumeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageBrokerServer is the server API for MessageBroker service.
type MessageBrokerServer interface {
	Produce(MessageBroker_ProduceServer) error
	Consume(MessageBroker_ConsumeServer) error
}

// UnimplementedMessageBrokerServer can be embedded to have forward compatible implementations.
type UnimplementedMessageBrokerServer struct {
}

func (*UnimplementedMessageBrokerServer) Produce(MessageBroker_ProduceServer) error {
	return status.Errorf(codes.Unimplemented, "method Produce not implemented")
}
func (*UnimplementedMessageBrokerServer) Consume(MessageBroker_ConsumeServer) error {
	return status.Errorf(codes.Unimplemented, "method Consume not implemented")
}

func RegisterMessageBrokerServer(s *grpc.Server, srv MessageBrokerServer) {
	s.RegisterService(&_MessageBroker_serviceDesc, srv)
}

func _MessageBroker_Produce_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageBrokerServer).Produce(&messageBrokerProduceServer{stream})
}

type MessageBroker_ProduceServer interface {
	SendAndClose(*ProduceResponse) error
	Recv() (*ProduceRequest, error)
	grpc.ServerStream
}

type messageBrokerProduceServer struct {
	grpc.ServerStream
}

func (x *messageBrokerProduceServer) SendAndClose(m *ProduceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageBrokerProduceServer) Recv() (*ProduceRequest, error) {
	m := new(ProduceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _MessageBroker_Consume_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageBrokerServer).Consume(&messageBrokerConsumeServer{stream})
}

type MessageBroker_ConsumeServer interface {
	Send(*ConsumeResponse) error
	Recv() (*ConsumeRequest, error)
	grpc.ServerStream
}

type messageBrokerConsumeServer struct {
	grpc.ServerStream
}

func (x *messageBrokerConsumeServer) Send(m *ConsumeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageBrokerConsumeServer) Recv() (*ConsumeRequest, error) {
	m := new(ConsumeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MessageBroker_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mbproto.MessageBroker",
	HandlerType: (*MessageBrokerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Produce",
			Handler:       _MessageBroker_Produce_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Consume",
			Handler:       _MessageBroker_Consume_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "spec.proto",
}
