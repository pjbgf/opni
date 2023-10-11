// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v1.0.0
// source: github.com/rancher/opni/pkg/apis/alerting/v2/receiver.proto

package v2

import (
	alertmanager "github.com/rancher/opni/internal/alertmanager"
	v1 "github.com/rancher/opni/pkg/apis/core/v1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OpniReceiver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reference *v1.Reference          `protobuf:"bytes,1,opt,name=reference,proto3" json:"reference,omitempty"`
	Receiver  *alertmanager.Receiver `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Clusters  []string               `protobuf:"bytes,3,rep,name=clusters,proto3" json:"clusters,omitempty"`
	// TODO : this should be a label matcher proto
	Matchers map[string]string `protobuf:"bytes,4,rep,name=matchers,proto3" json:"matchers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Metadata map[string]string `protobuf:"bytes,5,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Revision int64             `protobuf:"varint,6,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (x *OpniReceiver) Reset() {
	*x = OpniReceiver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpniReceiver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpniReceiver) ProtoMessage() {}

func (x *OpniReceiver) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpniReceiver.ProtoReflect.Descriptor instead.
func (*OpniReceiver) Descriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescGZIP(), []int{0}
}

func (x *OpniReceiver) GetReference() *v1.Reference {
	if x != nil {
		return x.Reference
	}
	return nil
}

func (x *OpniReceiver) GetReceiver() *alertmanager.Receiver {
	if x != nil {
		return x.Receiver
	}
	return nil
}

func (x *OpniReceiver) GetClusters() []string {
	if x != nil {
		return x.Clusters
	}
	return nil
}

func (x *OpniReceiver) GetMatchers() map[string]string {
	if x != nil {
		return x.Matchers
	}
	return nil
}

func (x *OpniReceiver) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *OpniReceiver) GetRevision() int64 {
	if x != nil {
		return x.Revision
	}
	return 0
}

type DeleteReceiverRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reference *v1.Reference `protobuf:"bytes,1,opt,name=reference,proto3" json:"reference,omitempty"`
	Revision  int64         `protobuf:"varint,2,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (x *DeleteReceiverRequest) Reset() {
	*x = DeleteReceiverRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteReceiverRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteReceiverRequest) ProtoMessage() {}

func (x *DeleteReceiverRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteReceiverRequest.ProtoReflect.Descriptor instead.
func (*DeleteReceiverRequest) Descriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescGZIP(), []int{1}
}

func (x *DeleteReceiverRequest) GetReference() *v1.Reference {
	if x != nil {
		return x.Reference
	}
	return nil
}

func (x *DeleteReceiverRequest) GetRevision() int64 {
	if x != nil {
		return x.Revision
	}
	return 0
}

type OpniReceiverList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*OpniReceiver `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *OpniReceiverList) Reset() {
	*x = OpniReceiverList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpniReceiverList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpniReceiverList) ProtoMessage() {}

func (x *OpniReceiverList) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpniReceiverList.ProtoReflect.Descriptor instead.
func (*OpniReceiverList) Descriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescGZIP(), []int{2}
}

func (x *OpniReceiverList) GetItems() []*OpniReceiver {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto protoreflect.FileDescriptor

var file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDesc = []byte{
	0x0a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x6e, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x32, 0x2f, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61,
	0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x6e, 0x69,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x6e, 0x69, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7, 0x03, 0x0a,
	0x0c, 0x4f, 0x70, 0x6e, 0x69, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x2d, 0x0a,
	0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x08,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x73, 0x12, 0x40, 0x0a, 0x08,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x6e, 0x69, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x12, 0x40,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x24, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x6e, 0x69,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x3b, 0x0a, 0x0d,
	0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x62, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2d, 0x0a, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x09, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x40, 0x0a, 0x10, 0x4f, 0x70,
	0x6e, 0x69, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2c,
	0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x6e, 0x69, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x72, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x32, 0xbb, 0x03, 0x0a,
	0x0e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12,
	0x49, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x0f,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x1a,
	0x16, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x6e, 0x69, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12,
	0x09, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x57, 0x0a, 0x0d, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x1a, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f,
	0x70, 0x6e, 0x69, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x22,
	0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x73, 0x12, 0x4c, 0x0a, 0x0b, 0x50, 0x75, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x12, 0x16, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70,
	0x6e, 0x69, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x1a, 0x0f, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x0e, 0x3a, 0x01, 0x2a, 0x1a, 0x09, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x12, 0x5c, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x11, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0b, 0x2a, 0x09, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12,
	0x59, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12,
	0x16, 0x2e, 0x61, 0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x4f, 0x70, 0x6e, 0x69, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x72,
	0x2f, 0x6f, 0x70, 0x6e, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61,
	0x6c, 0x65, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescOnce sync.Once
	file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescData = file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDesc
)

func file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescGZIP() []byte {
	file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescOnce.Do(func() {
		file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescData)
	})
	return file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDescData
}

var file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_goTypes = []interface{}{
	(*OpniReceiver)(nil),          // 0: alerting.OpniReceiver
	(*DeleteReceiverRequest)(nil), // 1: alerting.DeleteReceiverRequest
	(*OpniReceiverList)(nil),      // 2: alerting.OpniReceiverList
	nil,                           // 3: alerting.OpniReceiver.MatchersEntry
	nil,                           // 4: alerting.OpniReceiver.MetadataEntry
	(*v1.Reference)(nil),          // 5: core.Reference
	(*alertmanager.Receiver)(nil), // 6: alertmanager.Receiver
	(*emptypb.Empty)(nil),         // 7: google.protobuf.Empty
}
var file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_depIdxs = []int32{
	5,  // 0: alerting.OpniReceiver.reference:type_name -> core.Reference
	6,  // 1: alerting.OpniReceiver.receiver:type_name -> alertmanager.Receiver
	3,  // 2: alerting.OpniReceiver.matchers:type_name -> alerting.OpniReceiver.MatchersEntry
	4,  // 3: alerting.OpniReceiver.metadata:type_name -> alerting.OpniReceiver.MetadataEntry
	5,  // 4: alerting.DeleteReceiverRequest.reference:type_name -> core.Reference
	0,  // 5: alerting.OpniReceiverList.items:type_name -> alerting.OpniReceiver
	5,  // 6: alerting.ReceiverServer.GetReceiver:input_type -> core.Reference
	7,  // 7: alerting.ReceiverServer.ListReceivers:input_type -> google.protobuf.Empty
	0,  // 8: alerting.ReceiverServer.PutReceiver:input_type -> alerting.OpniReceiver
	1,  // 9: alerting.ReceiverServer.DeleteReceiver:input_type -> alerting.DeleteReceiverRequest
	0,  // 10: alerting.ReceiverServer.TestReceiver:input_type -> alerting.OpniReceiver
	0,  // 11: alerting.ReceiverServer.GetReceiver:output_type -> alerting.OpniReceiver
	2,  // 12: alerting.ReceiverServer.ListReceivers:output_type -> alerting.OpniReceiverList
	5,  // 13: alerting.ReceiverServer.PutReceiver:output_type -> core.Reference
	7,  // 14: alerting.ReceiverServer.DeleteReceiver:output_type -> google.protobuf.Empty
	7,  // 15: alerting.ReceiverServer.TestReceiver:output_type -> google.protobuf.Empty
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_init() }
func file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_init() {
	if File_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpniReceiver); i {
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
		file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteReceiverRequest); i {
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
		file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpniReceiverList); i {
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
			RawDescriptor: file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_goTypes,
		DependencyIndexes: file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_depIdxs,
		MessageInfos:      file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_msgTypes,
	}.Build()
	File_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto = out.File
	file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_rawDesc = nil
	file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_goTypes = nil
	file_github_com_rancher_opni_pkg_apis_alerting_v2_receiver_proto_depIdxs = nil
}
