// Code generated by protoc-gen-go.
// source: gcgrpc.proto
// DO NOT EDIT!

/*
Package gcgrpc is a generated protocol buffer package.

It is generated from these files:
	gcgrpc.proto

It has these top-level messages:
	RetrieveRequest
	RetrieveResponse
*/
package gcgrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RetrieveRequest struct {
	Group string `protobuf:"bytes,1,opt,name=group" json:"group,omitempty"`
	Key   string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
}

func (m *RetrieveRequest) Reset()                    { *m = RetrieveRequest{} }
func (m *RetrieveRequest) String() string            { return proto.CompactTextString(m) }
func (*RetrieveRequest) ProtoMessage()               {}
func (*RetrieveRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type RetrieveResponse struct {
	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *RetrieveResponse) Reset()                    { *m = RetrieveResponse{} }
func (m *RetrieveResponse) String() string            { return proto.CompactTextString(m) }
func (*RetrieveResponse) ProtoMessage()               {}
func (*RetrieveResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*RetrieveRequest)(nil), "gcgrpc.RetrieveRequest")
	proto.RegisterType((*RetrieveResponse)(nil), "gcgrpc.RetrieveResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Peer service

type PeerClient interface {
	Retrieve(ctx context.Context, in *RetrieveRequest, opts ...grpc.CallOption) (*RetrieveResponse, error)
}

type peerClient struct {
	cc *grpc.ClientConn
}

func NewPeerClient(cc *grpc.ClientConn) PeerClient {
	return &peerClient{cc}
}

func (c *peerClient) Retrieve(ctx context.Context, in *RetrieveRequest, opts ...grpc.CallOption) (*RetrieveResponse, error) {
	out := new(RetrieveResponse)
	err := grpc.Invoke(ctx, "/gcgrpc.Peer/Retrieve", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Peer service

type PeerServer interface {
	Retrieve(context.Context, *RetrieveRequest) (*RetrieveResponse, error)
}

func RegisterPeerServer(s *grpc.Server, srv PeerServer) {
	s.RegisterService(&_Peer_serviceDesc, srv)
}

func _Peer_Retrieve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServer).Retrieve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gcgrpc.Peer/Retrieve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServer).Retrieve(ctx, req.(*RetrieveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Peer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gcgrpc.Peer",
	HandlerType: (*PeerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Retrieve",
			Handler:    _Peer_Retrieve_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("gcgrpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 151 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x4f, 0x4e, 0x2f,
	0x2a, 0x48, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0x94, 0x2c, 0xb9, 0xf8,
	0x83, 0x52, 0x4b, 0x8a, 0x32, 0x53, 0xcb, 0x52, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84,
	0x44, 0xb8, 0x58, 0xd3, 0x8b, 0xf2, 0x4b, 0x0b, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x20,
	0x1c, 0x21, 0x01, 0x2e, 0xe6, 0xec, 0xd4, 0x4a, 0x09, 0x26, 0xb0, 0x18, 0x88, 0xa9, 0xa4, 0xc1,
	0x25, 0x80, 0xd0, 0x5a, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0x0a, 0xd2, 0x5b, 0x96, 0x98, 0x53, 0x9a,
	0x0a, 0xd6, 0xcb, 0x13, 0x04, 0xe1, 0x18, 0xb9, 0x73, 0xb1, 0x04, 0xa4, 0xa6, 0x16, 0x09, 0xd9,
	0x73, 0x71, 0xc0, 0x74, 0x08, 0x89, 0xeb, 0x41, 0xdd, 0x83, 0x66, 0xbd, 0x94, 0x04, 0xa6, 0x04,
	0xc4, 0x70, 0x25, 0x86, 0x24, 0x36, 0xb0, 0xe3, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xbd,
	0x33, 0xf5, 0xbd, 0xcc, 0x00, 0x00, 0x00,
}
