// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chn.proto

/*
Package chn is a generated protocol buffer package.

It is generated from these files:
	chn.proto

It has these top-level messages:
	StoryReq
	StoryResp
	Story
*/
package handler

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StoryReq struct {
	Category string `protobuf:"bytes,1,opt,name=category" json:"category,omitempty"`
	Limit    int64  `protobuf:"varint,2,opt,name=limit" json:"limit,omitempty"`
}

func (m *StoryReq) Reset()                    { *m = StoryReq{} }
func (m *StoryReq) String() string            { return proto.CompactTextString(m) }
func (*StoryReq) ProtoMessage()               {}
func (*StoryReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StoryReq) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

func (m *StoryReq) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type StoryResp struct {
	Stories []*Story `protobuf:"bytes,1,rep,name=stories" json:"stories,omitempty"`
}

func (m *StoryResp) Reset()                    { *m = StoryResp{} }
func (m *StoryResp) String() string            { return proto.CompactTextString(m) }
func (*StoryResp) ProtoMessage()               {}
func (*StoryResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StoryResp) GetStories() []*Story {
	if m != nil {
		return m.Stories
	}
	return nil
}

type Story struct {
	By          string  `protobuf:"bytes,1,opt,name=by" json:"by,omitempty"`
	Descendants int64   `protobuf:"varint,2,opt,name=descendants" json:"descendants,omitempty"`
	Id          int64   `protobuf:"varint,3,opt,name=id" json:"id,omitempty"`
	Kids        []int64 `protobuf:"varint,4,rep,packed,name=kids" json:"kids,omitempty"`
	Score       int64   `protobuf:"varint,5,opt,name=score" json:"score,omitempty"`
	Type        string  `protobuf:"bytes,6,opt,name=type" json:"type,omitempty"`
	Title       string  `protobuf:"bytes,7,opt,name=title" json:"title,omitempty"`
	Url         string  `protobuf:"bytes,8,opt,name=url" json:"url,omitempty"`
	DomainName  string  `protobuf:"bytes,9,opt,name=domainName" json:"domainName,omitempty"`
}

func (m *Story) Reset()                    { *m = Story{} }
func (m *Story) String() string            { return proto.CompactTextString(m) }
func (*Story) ProtoMessage()               {}
func (*Story) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Story) GetBy() string {
	if m != nil {
		return m.By
	}
	return ""
}

func (m *Story) GetDescendants() int64 {
	if m != nil {
		return m.Descendants
	}
	return 0
}

func (m *Story) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Story) GetKids() []int64 {
	if m != nil {
		return m.Kids
	}
	return nil
}

func (m *Story) GetScore() int64 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *Story) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Story) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Story) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Story) GetDomainName() string {
	if m != nil {
		return m.DomainName
	}
	return ""
}

func init() {
	proto.RegisterType((*StoryReq)(nil), "service.chn.StoryReq")
	proto.RegisterType((*StoryResp)(nil), "service.chn.StoryResp")
	proto.RegisterType((*Story)(nil), "service.chn.Story")
}

func init() { proto.RegisterFile("chn.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x31, 0x4f, 0xfb, 0x30,
	0x10, 0xc5, 0x95, 0xb8, 0x69, 0x9b, 0xeb, 0x5f, 0x7f, 0xa1, 0x13, 0x20, 0xab, 0x03, 0x8a, 0x3a,
	0x65, 0x40, 0x19, 0xca, 0x02, 0x12, 0x13, 0x13, 0x13, 0x43, 0xba, 0xb1, 0xa5, 0xf6, 0x89, 0x5a,
	0x24, 0x76, 0x6a, 0x1b, 0xa4, 0x7c, 0x4a, 0xbe, 0x12, 0x8a, 0xdd, 0xa2, 0x0c, 0xdd, 0xee, 0xfd,
	0xee, 0x9d, 0xde, 0x9d, 0x0d, 0xb9, 0x38, 0xe8, 0xaa, 0xb7, 0xc6, 0x1b, 0x5c, 0x39, 0xb2, 0xdf,
	0x4a, 0x50, 0x25, 0x0e, 0x7a, 0xf3, 0x0c, 0xcb, 0x9d, 0x37, 0x76, 0xa8, 0xe9, 0x88, 0x6b, 0x58,
	0x8a, 0xc6, 0xd3, 0x87, 0xb1, 0x03, 0x4f, 0x8a, 0xa4, 0xcc, 0xeb, 0x3f, 0x8d, 0xd7, 0x90, 0xb5,
	0xaa, 0x53, 0x9e, 0xa7, 0x45, 0x52, 0xb2, 0x3a, 0x8a, 0xcd, 0x13, 0xe4, 0xa7, 0x69, 0xd7, 0xe3,
	0x3d, 0x2c, 0x9c, 0x37, 0x56, 0x91, 0xe3, 0x49, 0xc1, 0xca, 0xd5, 0x16, 0xab, 0x49, 0x52, 0x15,
	0x8d, 0x67, 0xcb, 0xe6, 0x27, 0x81, 0x2c, 0x20, 0xfc, 0x0f, 0xe9, 0xfe, 0x1c, 0x98, 0xee, 0x07,
	0x2c, 0x60, 0x25, 0xc9, 0x09, 0xd2, 0xb2, 0xd1, 0xde, 0x9d, 0x02, 0xa7, 0x68, 0x9c, 0x50, 0x92,
	0xb3, 0xd0, 0x48, 0x95, 0x44, 0x84, 0xd9, 0xa7, 0x92, 0x8e, 0xcf, 0x0a, 0x56, 0xb2, 0x3a, 0xd4,
	0xe3, 0xc2, 0x4e, 0x18, 0x4b, 0x3c, 0x8b, 0x0b, 0x07, 0x31, 0x3a, 0xfd, 0xd0, 0x13, 0x9f, 0x87,
	0xb4, 0x50, 0x8f, 0x4e, 0xaf, 0x7c, 0x4b, 0x7c, 0x11, 0x60, 0x14, 0x78, 0x05, 0xec, 0xcb, 0xb6,
	0x7c, 0x19, 0xd8, 0x58, 0xe2, 0x1d, 0x80, 0x34, 0x5d, 0xa3, 0xf4, 0x5b, 0xd3, 0x11, 0xcf, 0x43,
	0x63, 0x42, 0xb6, 0xaf, 0xf0, 0x2f, 0x1c, 0xb4, 0x8b, 0x47, 0xe3, 0x23, 0x2c, 0x76, 0xf1, 0x58,
	0xbc, 0xb9, 0xf0, 0x12, 0x74, 0x5c, 0xdf, 0x5e, 0xc2, 0xae, 0x7f, 0xc9, 0xde, 0x99, 0x38, 0xe8,
	0xfd, 0x3c, 0xfc, 0xd7, 0xc3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x6a, 0x3f, 0x29, 0xbc,
	0x01, 0x00, 0x00,
}
