// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ref.proto

/*
Package goref_proto is a generated protocol buffer package.

It is generated from these files:
	ref.proto

It has these top-level messages:
	Ref
	Location
	Position
*/
package goref_proto

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

type Type int32

const (
	Type_Instantiation  Type = 0
	Type_Call           Type = 1
	Type_Implementation Type = 2
	Type_Extension      Type = 3
	Type_Import         Type = 4
	Type_Reference      Type = 5
)

var Type_name = map[int32]string{
	0: "Instantiation",
	1: "Call",
	2: "Implementation",
	3: "Extension",
	4: "Import",
	5: "Reference",
}
var Type_value = map[string]int32{
	"Instantiation":  0,
	"Call":           1,
	"Implementation": 2,
	"Extension":      3,
	"Import":         4,
	"Reference":      5,
}

func (x Type) String() string {
	return proto.EnumName(Type_name, int32(x))
}
func (Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Ref struct {
	Version int64     `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	From    *Location `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	To      *Location `protobuf:"bytes,3,opt,name=to" json:"to,omitempty"`
	Type    Type      `protobuf:"varint,4,opt,name=type,enum=goref.proto.Type" json:"type,omitempty"`
}

func (m *Ref) Reset()                    { *m = Ref{} }
func (m *Ref) String() string            { return proto.CompactTextString(m) }
func (*Ref) ProtoMessage()               {}
func (*Ref) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Ref) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Ref) GetFrom() *Location {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Ref) GetTo() *Location {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Ref) GetType() Type {
	if m != nil {
		return m.Type
	}
	return Type_Instantiation
}

type Location struct {
	Position *Position `protobuf:"bytes,1,opt,name=position" json:"position,omitempty"`
	Package  string    `protobuf:"bytes,2,opt,name=package" json:"package,omitempty"`
	Ident    string    `protobuf:"bytes,3,opt,name=ident" json:"ident,omitempty"`
}

func (m *Location) Reset()                    { *m = Location{} }
func (m *Location) String() string            { return proto.CompactTextString(m) }
func (*Location) ProtoMessage()               {}
func (*Location) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Location) GetPosition() *Position {
	if m != nil {
		return m.Position
	}
	return nil
}

func (m *Location) GetPackage() string {
	if m != nil {
		return m.Package
	}
	return ""
}

func (m *Location) GetIdent() string {
	if m != nil {
		return m.Ident
	}
	return ""
}

type Position struct {
	Filename  string `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
	StartLine int32  `protobuf:"varint,2,opt,name=start_line,json=startLine" json:"start_line,omitempty"`
	StartCol  int32  `protobuf:"varint,3,opt,name=start_col,json=startCol" json:"start_col,omitempty"`
	EndLine   int32  `protobuf:"varint,4,opt,name=end_line,json=endLine" json:"end_line,omitempty"`
	EndCol    int32  `protobuf:"varint,5,opt,name=end_col,json=endCol" json:"end_col,omitempty"`
}

func (m *Position) Reset()                    { *m = Position{} }
func (m *Position) String() string            { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()               {}
func (*Position) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Position) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *Position) GetStartLine() int32 {
	if m != nil {
		return m.StartLine
	}
	return 0
}

func (m *Position) GetStartCol() int32 {
	if m != nil {
		return m.StartCol
	}
	return 0
}

func (m *Position) GetEndLine() int32 {
	if m != nil {
		return m.EndLine
	}
	return 0
}

func (m *Position) GetEndCol() int32 {
	if m != nil {
		return m.EndCol
	}
	return 0
}

func init() {
	proto.RegisterType((*Ref)(nil), "goref.proto.Ref")
	proto.RegisterType((*Location)(nil), "goref.proto.Location")
	proto.RegisterType((*Position)(nil), "goref.proto.Position")
	proto.RegisterEnum("goref.proto.Type", Type_name, Type_value)
}

func init() { proto.RegisterFile("ref.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 347 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x91, 0x4f, 0x4b, 0xf3, 0x40,
	0x10, 0xc6, 0xdf, 0xfc, 0x6b, 0x93, 0x29, 0x2d, 0xe9, 0xf0, 0x8a, 0x51, 0x11, 0x4a, 0xa1, 0x50,
	0x3d, 0x14, 0xac, 0x1f, 0xa1, 0x78, 0x28, 0xf4, 0x20, 0x8b, 0xf7, 0xb2, 0x26, 0x93, 0xb2, 0xb8,
	0xd9, 0x0d, 0x9b, 0x45, 0xec, 0x17, 0x11, 0x3f, 0xae, 0x64, 0x63, 0x8a, 0x7a, 0xf0, 0x96, 0x67,
	0x7e, 0xcf, 0x33, 0x33, 0x99, 0x85, 0xc4, 0x50, 0xb9, 0xaa, 0x8d, 0xb6, 0x1a, 0x47, 0x07, 0x7d,
	0x12, 0xf3, 0x0f, 0x0f, 0x02, 0x46, 0x25, 0x66, 0x30, 0x7c, 0x25, 0xd3, 0x08, 0xad, 0x32, 0x6f,
	0xe6, 0x2d, 0x03, 0xd6, 0x4b, 0xbc, 0x81, 0xb0, 0x34, 0xba, 0xca, 0xfc, 0x99, 0xb7, 0x1c, 0xad,
	0xcf, 0x56, 0xdf, 0xd2, 0xab, 0x9d, 0xce, 0xb9, 0x15, 0x5a, 0x31, 0x67, 0xc1, 0x05, 0xf8, 0x56,
	0x67, 0xc1, 0x5f, 0x46, 0xdf, 0x6a, 0x5c, 0x40, 0x68, 0x8f, 0x35, 0x65, 0xe1, 0xcc, 0x5b, 0x4e,
	0xd6, 0xd3, 0x1f, 0xc6, 0xa7, 0x63, 0x4d, 0xcc, 0xe1, 0x79, 0x05, 0x71, 0x1f, 0xc3, 0x3b, 0x88,
	0x6b, 0xdd, 0x08, 0xdb, 0xef, 0xf7, 0xbb, 0xff, 0xe3, 0x17, 0x64, 0x27, 0x5b, 0xfb, 0x47, 0x35,
	0xcf, 0x5f, 0xf8, 0x81, 0xdc, 0xea, 0x09, 0xeb, 0x25, 0xfe, 0x87, 0x48, 0x14, 0xa4, 0xac, 0xdb,
	0x34, 0x61, 0x9d, 0x98, 0xbf, 0x7b, 0x10, 0xf7, 0x6d, 0xf0, 0x12, 0xe2, 0x52, 0x48, 0x52, 0xbc,
	0x22, 0x37, 0x2f, 0x61, 0x27, 0x8d, 0xd7, 0x00, 0x8d, 0xe5, 0xc6, 0xee, 0xa5, 0x50, 0x5d, 0xef,
	0x88, 0x25, 0xae, 0xb2, 0x13, 0x8a, 0xf0, 0x0a, 0x3a, 0xb1, 0xcf, 0xb5, 0x74, 0x13, 0x22, 0x16,
	0xbb, 0xc2, 0x46, 0x4b, 0xbc, 0x80, 0x98, 0x54, 0xd1, 0x25, 0x43, 0xc7, 0x86, 0xa4, 0x0a, 0x97,
	0x3b, 0x87, 0xf6, 0xd3, 0xa5, 0x22, 0x47, 0x06, 0xa4, 0x8a, 0x8d, 0x96, 0xb7, 0x1c, 0xc2, 0xf6,
	0x2a, 0x38, 0x85, 0xf1, 0x56, 0x35, 0x96, 0x2b, 0x2b, 0xdc, 0x51, 0xd2, 0x7f, 0x18, 0x43, 0xb8,
	0xe1, 0x52, 0xa6, 0x1e, 0x22, 0x4c, 0xb6, 0x55, 0x2d, 0xa9, 0x22, 0x65, 0x3b, 0xea, 0xe3, 0x18,
	0x92, 0x87, 0x37, 0x4b, 0xaa, 0x7d, 0xc6, 0x34, 0x40, 0x80, 0xc1, 0xb6, 0xaa, 0xb5, 0xb1, 0x69,
	0xd8, 0x22, 0x46, 0x25, 0x19, 0x52, 0x39, 0xa5, 0xd1, 0xf3, 0xc0, 0x5d, 0xf1, 0xfe, 0x33, 0x00,
	0x00, 0xff, 0xff, 0x6d, 0xca, 0x56, 0x68, 0x26, 0x02, 0x00, 0x00,
}