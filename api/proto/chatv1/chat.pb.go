// Code generated manually to provide protobuf bindings for chat.proto. DO NOT EDIT.
package chatv1

import (
	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ServerNotice_Type enumerates server notice kinds.
type ServerNotice_Type int32

const (
	ServerNotice_TYPE_GENERIC     ServerNotice_Type = 0
	ServerNotice_TYPE_USER_JOINED ServerNotice_Type = 1
	ServerNotice_TYPE_USER_LEFT   ServerNotice_Type = 2
	ServerNotice_TYPE_ERROR       ServerNotice_Type = 3
)

// Enum value maps for ServerNotice_Type.
var (
	ServerNotice_Type_name = map[int32]string{
		0: "TYPE_GENERIC",
		1: "TYPE_USER_JOINED",
		2: "TYPE_USER_LEFT",
		3: "TYPE_ERROR",
	}
	ServerNotice_Type_value = map[string]int32{
		"TYPE_GENERIC":     0,
		"TYPE_USER_JOINED": 1,
		"TYPE_USER_LEFT":   2,
		"TYPE_ERROR":       3,
	}
)

func (x ServerNotice_Type) Enum() *ServerNotice_Type {
	p := new(ServerNotice_Type)
	*p = x
	return p
}

func (x ServerNotice_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServerNotice_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_chat_proto_enumTypes[0].Descriptor()
}

func (ServerNotice_Type) Type() protoreflect.EnumType {
	return &file_chat_proto_enumTypes[0]
}

func (x ServerNotice_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

func (ServerNotice_Type) EnumDescriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{6, 0}
}

// JoinRequest describes the information a client must send to join a room.
type JoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Room        string `protobuf:"bytes,2,opt,name=room,proto3" json:"room,omitempty"`
	DisplayName string `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
}

func (x *JoinRequest) Reset() {
	*x = JoinRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinRequest) ProtoMessage() {}

func (x *JoinRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*JoinRequest) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

func (x *JoinRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *JoinRequest) GetRoom() string {
	if x != nil {
		return x.Room
	}
	return ""
}

func (x *JoinRequest) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

// ChatPayload represents an arbitrary message sent by a client.
type ChatPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Room         string `protobuf:"bytes,2,opt,name=room,proto3" json:"room,omitempty"`
	Content      string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	TimestampUtc int64  `protobuf:"varint,4,opt,name=timestamp_utc,json=timestampUtc,proto3" json:"timestamp_utc,omitempty"`
}

func (x *ChatPayload) Reset() {
	*x = ChatPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatPayload) ProtoMessage() {}

func (x *ChatPayload) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ChatPayload) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

func (x *ChatPayload) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ChatPayload) GetRoom() string {
	if x != nil {
		return x.Room
	}
	return ""
}

func (x *ChatPayload) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *ChatPayload) GetTimestampUtc() int64 {
	if x != nil {
		return x.TimestampUtc
	}
	return 0
}

// LeaveRequest notifies the server that a client wants to disconnect from a room.
type LeaveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Room   string `protobuf:"bytes,2,opt,name=room,proto3" json:"room,omitempty"`
}

func (x *LeaveRequest) Reset() {
	*x = LeaveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveRequest) ProtoMessage() {}

func (x *LeaveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*LeaveRequest) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{2}
}

func (x *LeaveRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *LeaveRequest) GetRoom() string {
	if x != nil {
		return x.Room
	}
	return ""
}

type ClientEnvelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*ClientEnvelope_Join
	//	*ClientEnvelope_Chat
	//	*ClientEnvelope_Leave
	Message isClientEnvelope_Message `protobuf_oneof:"message"`
}

func (x *ClientEnvelope) Reset() {
	*x = ClientEnvelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientEnvelope) ProtoMessage() {}

func (x *ClientEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ClientEnvelope) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{3}
}

func (m *ClientEnvelope) GetMessage() isClientEnvelope_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *ClientEnvelope) GetJoin() *JoinRequest {
	if x, ok := x.GetMessage().(*ClientEnvelope_Join); ok {
		return x.Join
	}
	return nil
}

func (x *ClientEnvelope) GetChat() *ChatPayload {
	if x, ok := x.GetMessage().(*ClientEnvelope_Chat); ok {
		return x.Chat
	}
	return nil
}

func (x *ClientEnvelope) GetLeave() *LeaveRequest {
	if x, ok := x.GetMessage().(*ClientEnvelope_Leave); ok {
		return x.Leave
	}
	return nil
}

type isClientEnvelope_Message interface {
	isClientEnvelope_Message()
}

type ClientEnvelope_Join struct {
	Join *JoinRequest `protobuf:"bytes,1,opt,name=join,proto3,oneof"`
}

type ClientEnvelope_Chat struct {
	Chat *ChatPayload `protobuf:"bytes,2,opt,name=chat,proto3,oneof"`
}

type ClientEnvelope_Leave struct {
	Leave *LeaveRequest `protobuf:"bytes,3,opt,name=leave,proto3,oneof"`
}

func (*ClientEnvelope_Join) isClientEnvelope_Message()  {}
func (*ClientEnvelope_Chat) isClientEnvelope_Message()  {}
func (*ClientEnvelope_Leave) isClientEnvelope_Message() {}

func (x *ClientEnvelope) SetMessage(is isClientEnvelope_Message) {
	x.Message = is
}

func (x *ClientEnvelope) hasMessage() bool {
	return x.Message != nil
}

// JoinAck confirms that the user joined the requested room.
type JoinAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId         string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Room           string `protobuf:"bytes,2,opt,name=room,proto3" json:"room,omitempty"`
	WelcomeMessage string `protobuf:"bytes,3,opt,name=welcome_message,json=welcomeMessage,proto3" json:"welcome_message,omitempty"`
}

func (x *JoinAck) Reset() {
	*x = JoinAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinAck) ProtoMessage() {}

func (x *JoinAck) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*JoinAck) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{4}
}

func (x *JoinAck) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *JoinAck) GetRoom() string {
	if x != nil {
		return x.Room
	}
	return ""
}

func (x *JoinAck) GetWelcomeMessage() string {
	if x != nil {
		return x.WelcomeMessage
	}
	return ""
}

type ServerEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Event:
	//
	//	*ServerEvent_Joined
	//	*ServerEvent_Broadcast
	//	*ServerEvent_Notice
	Event isServerEvent_Event `protobuf_oneof:"event"`
}

func (x *ServerEvent) Reset() {
	*x = ServerEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerEvent) ProtoMessage() {}

func (x *ServerEvent) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ServerEvent) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{5}
}

func (m *ServerEvent) GetEvent() isServerEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *ServerEvent) GetJoined() *JoinAck {
	if x, ok := x.GetEvent().(*ServerEvent_Joined); ok {
		return x.Joined
	}
	return nil
}

func (x *ServerEvent) GetBroadcast() *ChatPayload {
	if x, ok := x.GetEvent().(*ServerEvent_Broadcast); ok {
		return x.Broadcast
	}
	return nil
}

func (x *ServerEvent) GetNotice() *ServerNotice {
	if x, ok := x.GetEvent().(*ServerEvent_Notice); ok {
		return x.Notice
	}
	return nil
}

type isServerEvent_Event interface {
	isServerEvent_Event()
}

type ServerEvent_Joined struct {
	Joined *JoinAck `protobuf:"bytes,1,opt,name=joined,proto3,oneof"`
}

type ServerEvent_Broadcast struct {
	Broadcast *ChatPayload `protobuf:"bytes,2,opt,name=broadcast,proto3,oneof"`
}

type ServerEvent_Notice struct {
	Notice *ServerNotice `protobuf:"bytes,3,opt,name=notice,proto3,oneof"`
}

func (*ServerEvent_Joined) isServerEvent_Event()    {}
func (*ServerEvent_Broadcast) isServerEvent_Event() {}
func (*ServerEvent_Notice) isServerEvent_Event()    {}

// ServerNotice conveys system-level announcements (errors, user events).
type ServerNotice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    ServerNotice_Type `protobuf:"varint,1,opt,name=type,proto3,enum=chat.v1.ServerNotice_Type" json:"type,omitempty"`
	Message string            `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	UserId  string            `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Room    string            `protobuf:"bytes,4,opt,name=room,proto3" json:"room,omitempty"`
}

func (x *ServerNotice) Reset() {
	*x = ServerNotice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerNotice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerNotice) ProtoMessage() {}

func (x *ServerNotice) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ServerNotice) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{6}
}

func (x *ServerNotice) GetType() ServerNotice_Type {
	if x != nil {
		return x.Type
	}
	return ServerNotice_TYPE_GENERIC
}

func (x *ServerNotice) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ServerNotice) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ServerNotice) GetRoom() string {
	if x != nil {
		return x.Room
	}
	return ""
}

var File_chat_proto protoreflect.FileDescriptor

var file_chat_proto_rawDescOnce sync.Once
var file_chat_proto_rawDescData = file_chat_proto_rawDesc

func file_chat_proto_rawDescGZIP() []byte {
	file_chat_proto_rawDescOnce.Do(func() {
		file_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_proto_rawDescData)
	})
	return file_chat_proto_rawDescData
}

var file_chat_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_chat_proto_goTypes = []interface{}{
	(ServerNotice_Type)(0), // 0: chat.v1.ServerNotice.Type
	(*JoinRequest)(nil),    // 1: chat.v1.JoinRequest
	(*ChatPayload)(nil),    // 2: chat.v1.ChatPayload
	(*LeaveRequest)(nil),   // 3: chat.v1.LeaveRequest
	(*ClientEnvelope)(nil), // 4: chat.v1.ClientEnvelope
	(*JoinAck)(nil),        // 5: chat.v1.JoinAck
	(*ServerEvent)(nil),    // 6: chat.v1.ServerEvent
	(*ServerNotice)(nil),   // 7: chat.v1.ServerNotice
}
var file_chat_proto_depIdxs = []int32{
	1, // 0: chat.v1.ClientEnvelope.join:type_name -> chat.v1.JoinRequest
	2, // 1: chat.v1.ClientEnvelope.chat:type_name -> chat.v1.ChatPayload
	3, // 2: chat.v1.ClientEnvelope.leave:type_name -> chat.v1.LeaveRequest
	5, // 3: chat.v1.ServerEvent.joined:type_name -> chat.v1.JoinAck
	2, // 4: chat.v1.ServerEvent.broadcast:type_name -> chat.v1.ChatPayload
	7, // 5: chat.v1.ServerEvent.notice:type_name -> chat.v1.ServerNotice
	0, // 6: chat.v1.ServerNotice.type:type_name -> chat.v1.ServerNotice.Type
	4, // 7: chat.v1.ChatService.Channel:input_type -> chat.v1.ClientEnvelope
	6, // 8: chat.v1.ChatService.Channel:output_type -> chat.v1.ServerEvent
	8, // 9: [start of service method outputs]
	7, // 10: [start of service method inputs]
	7, // 11: [start of extension field dependencies]
	7, // 12: [start of extension field targets]
	0, // 13: [start of message field dependencies]
}

func init() { file_chat_proto_init() }
func file_chat_proto_init() {
	if File_chat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinRequest); i {
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
		file_chat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatPayload); i {
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
		file_chat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaveRequest); i {
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
		file_chat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientEnvelope); i {
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
		file_chat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinAck); i {
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
		file_chat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerEvent); i {
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
		file_chat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerNotice); i {
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
	file_chat_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*ClientEnvelope_Join)(nil),
		(*ClientEnvelope_Chat)(nil),
		(*ClientEnvelope_Leave)(nil),
	}
	file_chat_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*ServerEvent_Joined)(nil),
		(*ServerEvent_Broadcast)(nil),
		(*ServerEvent_Notice)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chat_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_proto_goTypes,
		DependencyIndexes: file_chat_proto_depIdxs,
		EnumInfos:         file_chat_proto_enumTypes,
		MessageInfos:      file_chat_proto_msgTypes,
	}.Build()
	File_chat_proto = out.File
	file_chat_proto_rawDesc = nil
	file_chat_proto_goTypes = nil
	file_chat_proto_depIdxs = nil
}

var file_chat_proto_rawDesc = func() []byte {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("chat.proto"),
		Package: proto.String("chat.v1"),
		Syntax:  proto.String("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("github.com/lechitz/chat-grpc/api/proto/chatv1;chatv1"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("JoinRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("user_id"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: proto.String("userId"),
					},
					{
						Name:   proto.String("room"),
						Number: proto.Int32(2),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:     proto.String("display_name"),
						Number:   proto.Int32(3),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: proto.String("displayName"),
					},
				},
			},
			{
				Name: proto.String("ChatPayload"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("user_id"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: proto.String("userId"),
					},
					{
						Name:   proto.String("room"),
						Number: proto.Int32(2),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:   proto.String("content"),
						Number: proto.Int32(3),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:     proto.String("timestamp_utc"),
						Number:   proto.Int32(4),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum(),
						JsonName: proto.String("timestampUtc"),
					},
				},
			},
			{
				Name: proto.String("LeaveRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("user_id"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: proto.String("userId"),
					},
					{
						Name:   proto.String("room"),
						Number: proto.Int32(2),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
				},
			},
			{
				Name: proto.String("ClientEnvelope"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:       proto.String("join"),
						Number:     proto.Int32(1),
						Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:       descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName:   proto.String(".chat.v1.JoinRequest"),
						OneofIndex: proto.Int32(0),
					},
					{
						Name:       proto.String("chat"),
						Number:     proto.Int32(2),
						Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:       descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName:   proto.String(".chat.v1.ChatPayload"),
						OneofIndex: proto.Int32(0),
					},
					{
						Name:       proto.String("leave"),
						Number:     proto.Int32(3),
						Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:       descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName:   proto.String(".chat.v1.LeaveRequest"),
						OneofIndex: proto.Int32(0),
					},
				},
				OneofDecl: []*descriptorpb.OneofDescriptorProto{
					{Name: proto.String("message")},
				},
			},
			{
				Name: proto.String("JoinAck"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("user_id"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: proto.String("userId"),
					},
					{
						Name:   proto.String("room"),
						Number: proto.Int32(2),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:     proto.String("welcome_message"),
						Number:   proto.Int32(3),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: proto.String("welcomeMessage"),
					},
				},
			},
			{
				Name: proto.String("ServerEvent"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:       proto.String("joined"),
						Number:     proto.Int32(1),
						Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:       descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName:   proto.String(".chat.v1.JoinAck"),
						OneofIndex: proto.Int32(0),
					},
					{
						Name:       proto.String("broadcast"),
						Number:     proto.Int32(2),
						Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:       descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName:   proto.String(".chat.v1.ChatPayload"),
						OneofIndex: proto.Int32(0),
					},
					{
						Name:       proto.String("notice"),
						Number:     proto.Int32(3),
						Label:      descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:       descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
						TypeName:   proto.String(".chat.v1.ServerNotice"),
						OneofIndex: proto.Int32(0),
					},
				},
				OneofDecl: []*descriptorpb.OneofDescriptorProto{
					{Name: proto.String("event")},
				},
			},
			{
				Name: proto.String("ServerNotice"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("type"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
						TypeName: proto.String(".chat.v1.ServerNotice.Type"),
					},
					{
						Name:   proto.String("message"),
						Number: proto.Int32(2),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
					{
						Name:     proto.String("user_id"),
						Number:   proto.Int32(3),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						JsonName: proto.String("userId"),
					},
					{
						Name:   proto.String("room"),
						Number: proto.Int32(4),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
				},
				NestedType: []*descriptorpb.DescriptorProto{},
				EnumType: []*descriptorpb.EnumDescriptorProto{
					{
						Name: proto.String("Type"),
						Value: []*descriptorpb.EnumValueDescriptorProto{
							{Name: proto.String("TYPE_GENERIC"), Number: proto.Int32(0)},
							{Name: proto.String("TYPE_USER_JOINED"), Number: proto.Int32(1)},
							{Name: proto.String("TYPE_USER_LEFT"), Number: proto.Int32(2)},
							{Name: proto.String("TYPE_ERROR"), Number: proto.Int32(3)},
						},
					},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("ChatService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:            proto.String("Channel"),
						InputType:       proto.String(".chat.v1.ClientEnvelope"),
						OutputType:      proto.String(".chat.v1.ServerEvent"),
						ClientStreaming: proto.Bool(true),
						ServerStreaming: proto.Bool(true),
					},
				},
			},
		},
	}
	raw, _ := proto.Marshal(fd)
	return raw
}()
