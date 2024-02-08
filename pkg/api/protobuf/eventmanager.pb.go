// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: api/protobuf/eventmanager.proto

package eventmanager

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

type MakeEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId int64  `protobuf:"varint,1,opt,name=senderId,proto3" json:"senderId,omitempty"`
	Time     int64  `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	Name     string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *MakeEventRequest) Reset() {
	*x = MakeEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MakeEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MakeEventRequest) ProtoMessage() {}

func (x *MakeEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MakeEventRequest.ProtoReflect.Descriptor instead.
func (*MakeEventRequest) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{0}
}

func (x *MakeEventRequest) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *MakeEventRequest) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *MakeEventRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type MakeEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId int64 `protobuf:"varint,1,opt,name=eventId,proto3" json:"eventId,omitempty"`
}

func (x *MakeEventResponse) Reset() {
	*x = MakeEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MakeEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MakeEventResponse) ProtoMessage() {}

func (x *MakeEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MakeEventResponse.ProtoReflect.Descriptor instead.
func (*MakeEventResponse) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{1}
}

func (x *MakeEventResponse) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

type GetEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId int64 `protobuf:"varint,1,opt,name=senderId,proto3" json:"senderId,omitempty"`
	EventId  int64 `protobuf:"varint,2,opt,name=eventId,proto3" json:"eventId,omitempty"`
}

func (x *GetEventRequest) Reset() {
	*x = GetEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventRequest) ProtoMessage() {}

func (x *GetEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventRequest.ProtoReflect.Descriptor instead.
func (*GetEventRequest) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{2}
}

func (x *GetEventRequest) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *GetEventRequest) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

type GetEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId int64  `protobuf:"varint,1,opt,name=senderId,proto3" json:"senderId,omitempty"`
	EventId  int64  `protobuf:"varint,2,opt,name=eventId,proto3" json:"eventId,omitempty"`
	Time     int64  `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	Name     string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetEventResponse) Reset() {
	*x = GetEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventResponse) ProtoMessage() {}

func (x *GetEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventResponse.ProtoReflect.Descriptor instead.
func (*GetEventResponse) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{3}
}

func (x *GetEventResponse) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *GetEventResponse) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *GetEventResponse) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *GetEventResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DeleteEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId int64 `protobuf:"varint,1,opt,name=senderId,proto3" json:"senderId,omitempty"`
	EventId  int64 `protobuf:"varint,2,opt,name=eventId,proto3" json:"eventId,omitempty"`
}

func (x *DeleteEventRequest) Reset() {
	*x = DeleteEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteEventRequest) ProtoMessage() {}

func (x *DeleteEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteEventRequest.ProtoReflect.Descriptor instead.
func (*DeleteEventRequest) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteEventRequest) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *DeleteEventRequest) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

type DeleteEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId int64 `protobuf:"varint,1,opt,name=eventId,proto3" json:"eventId,omitempty"`
}

func (x *DeleteEventResponse) Reset() {
	*x = DeleteEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteEventResponse) ProtoMessage() {}

func (x *DeleteEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteEventResponse.ProtoReflect.Descriptor instead.
func (*DeleteEventResponse) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteEventResponse) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

type GetEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId int64 `protobuf:"varint,1,opt,name=senderId,proto3" json:"senderId,omitempty"`
}

func (x *GetEventsRequest) Reset() {
	*x = GetEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventsRequest) ProtoMessage() {}

func (x *GetEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventsRequest.ProtoReflect.Descriptor instead.
func (*GetEventsRequest) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{6}
}

func (x *GetEventsRequest) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

type GetEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderId int64  `protobuf:"varint,1,opt,name=senderId,proto3" json:"senderId,omitempty"`
	EventId  int64  `protobuf:"varint,2,opt,name=eventId,proto3" json:"eventId,omitempty"`
	Time     int64  `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	Name     string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetEventsResponse) Reset() {
	*x = GetEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_eventmanager_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventsResponse) ProtoMessage() {}

func (x *GetEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_eventmanager_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventsResponse.ProtoReflect.Descriptor instead.
func (*GetEventsResponse) Descriptor() ([]byte, []int) {
	return file_api_protobuf_eventmanager_proto_rawDescGZIP(), []int{7}
}

func (x *GetEventsResponse) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *GetEventsResponse) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *GetEventsResponse) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *GetEventsResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_api_protobuf_eventmanager_proto protoreflect.FileDescriptor

var file_api_protobuf_eventmanager_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x22,
	0x56, 0x0a, 0x10, 0x4d, 0x61, 0x6b, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2d, 0x0a, 0x11, 0x4d, 0x61, 0x6b, 0x65, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x47, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22,
	0x70, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x4a, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x2f, 0x0a,
	0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x2e,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x71,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x32, 0xcd, 0x02, 0x0a, 0x06, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x4e, 0x0a, 0x09,
	0x4d, 0x61, 0x6b, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x6b, 0x65, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x6b, 0x65, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0b, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x50, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1e, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30,
	0x01, 0x42, 0x3e, 0x5a, 0x3c, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x52, 0x61, 0x69, 0x6d, 0x67, 0x75, 0x7a, 0x68,
	0x69, 0x6e, 0x6f, 0x76, 0x2f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2d, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_protobuf_eventmanager_proto_rawDescOnce sync.Once
	file_api_protobuf_eventmanager_proto_rawDescData = file_api_protobuf_eventmanager_proto_rawDesc
)

func file_api_protobuf_eventmanager_proto_rawDescGZIP() []byte {
	file_api_protobuf_eventmanager_proto_rawDescOnce.Do(func() {
		file_api_protobuf_eventmanager_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_protobuf_eventmanager_proto_rawDescData)
	})
	return file_api_protobuf_eventmanager_proto_rawDescData
}

var file_api_protobuf_eventmanager_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_protobuf_eventmanager_proto_goTypes = []interface{}{
	(*MakeEventRequest)(nil),    // 0: eventmanager.MakeEventRequest
	(*MakeEventResponse)(nil),   // 1: eventmanager.MakeEventResponse
	(*GetEventRequest)(nil),     // 2: eventmanager.GetEventRequest
	(*GetEventResponse)(nil),    // 3: eventmanager.GetEventResponse
	(*DeleteEventRequest)(nil),  // 4: eventmanager.DeleteEventRequest
	(*DeleteEventResponse)(nil), // 5: eventmanager.DeleteEventResponse
	(*GetEventsRequest)(nil),    // 6: eventmanager.GetEventsRequest
	(*GetEventsResponse)(nil),   // 7: eventmanager.GetEventsResponse
}
var file_api_protobuf_eventmanager_proto_depIdxs = []int32{
	0, // 0: eventmanager.Events.MakeEvent:input_type -> eventmanager.MakeEventRequest
	2, // 1: eventmanager.Events.GetEvent:input_type -> eventmanager.GetEventRequest
	4, // 2: eventmanager.Events.DeleteEvent:input_type -> eventmanager.DeleteEventRequest
	6, // 3: eventmanager.Events.GetEvents:input_type -> eventmanager.GetEventsRequest
	1, // 4: eventmanager.Events.MakeEvent:output_type -> eventmanager.MakeEventResponse
	3, // 5: eventmanager.Events.GetEvent:output_type -> eventmanager.GetEventResponse
	5, // 6: eventmanager.Events.DeleteEvent:output_type -> eventmanager.DeleteEventResponse
	7, // 7: eventmanager.Events.GetEvents:output_type -> eventmanager.GetEventsResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_protobuf_eventmanager_proto_init() }
func file_api_protobuf_eventmanager_proto_init() {
	if File_api_protobuf_eventmanager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_protobuf_eventmanager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MakeEventRequest); i {
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
		file_api_protobuf_eventmanager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MakeEventResponse); i {
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
		file_api_protobuf_eventmanager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventRequest); i {
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
		file_api_protobuf_eventmanager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventResponse); i {
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
		file_api_protobuf_eventmanager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteEventRequest); i {
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
		file_api_protobuf_eventmanager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteEventResponse); i {
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
		file_api_protobuf_eventmanager_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventsRequest); i {
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
		file_api_protobuf_eventmanager_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventsResponse); i {
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
			RawDescriptor: file_api_protobuf_eventmanager_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_protobuf_eventmanager_proto_goTypes,
		DependencyIndexes: file_api_protobuf_eventmanager_proto_depIdxs,
		MessageInfos:      file_api_protobuf_eventmanager_proto_msgTypes,
	}.Build()
	File_api_protobuf_eventmanager_proto = out.File
	file_api_protobuf_eventmanager_proto_rawDesc = nil
	file_api_protobuf_eventmanager_proto_goTypes = nil
	file_api_protobuf_eventmanager_proto_depIdxs = nil
}
