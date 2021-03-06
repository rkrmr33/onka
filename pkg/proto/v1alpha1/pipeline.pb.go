// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pkg/proto/v1alpha1/pipeline.proto

package v1alpha1

import (
	fmt "fmt"
	types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PipelineState int32

const (
	PipelineState_PIPELINE_STATE_UNSPECIFIED PipelineState = 0
	PipelineState_PIPELINE_STATE_PENDING     PipelineState = 1
	PipelineState_PIPELINE_STATE_PREPARE     PipelineState = 2
	PipelineState_PIPELINE_STATE_RUNNING     PipelineState = 3
	PipelineState_PIPELINE_STATE_SUCCESS     PipelineState = 4
	PipelineState_PIPELINE_STATE_FAILURE     PipelineState = 5
	PipelineState_PIPELINE_STATE_ERROR       PipelineState = 6
)

var PipelineState_name = map[int32]string{
	0: "PIPELINE_STATE_UNSPECIFIED",
	1: "PIPELINE_STATE_PENDING",
	2: "PIPELINE_STATE_PREPARE",
	3: "PIPELINE_STATE_RUNNING",
	4: "PIPELINE_STATE_SUCCESS",
	5: "PIPELINE_STATE_FAILURE",
	6: "PIPELINE_STATE_ERROR",
}

var PipelineState_value = map[string]int32{
	"PIPELINE_STATE_UNSPECIFIED": 0,
	"PIPELINE_STATE_PENDING":     1,
	"PIPELINE_STATE_PREPARE":     2,
	"PIPELINE_STATE_RUNNING":     3,
	"PIPELINE_STATE_SUCCESS":     4,
	"PIPELINE_STATE_FAILURE":     5,
	"PIPELINE_STATE_ERROR":       6,
}

func (x PipelineState) String() string {
	return proto.EnumName(PipelineState_name, int32(x))
}

func (PipelineState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3e7df5ab5b2e195a, []int{0}
}

type Pipeline struct {
	Metadata             *Metadata         `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec                 *PipelineSpec     `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
	Statuses             []*PipelineStatus `protobuf:"bytes,3,rep,name=statuses,proto3" json:"statuses,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Pipeline) Reset()         { *m = Pipeline{} }
func (m *Pipeline) String() string { return proto.CompactTextString(m) }
func (*Pipeline) ProtoMessage()    {}
func (*Pipeline) Descriptor() ([]byte, []int) {
	return fileDescriptor_3e7df5ab5b2e195a, []int{0}
}
func (m *Pipeline) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Pipeline) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Pipeline.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Pipeline) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pipeline.Merge(m, src)
}
func (m *Pipeline) XXX_Size() int {
	return m.Size()
}
func (m *Pipeline) XXX_DiscardUnknown() {
	xxx_messageInfo_Pipeline.DiscardUnknown(m)
}

var xxx_messageInfo_Pipeline proto.InternalMessageInfo

func (m *Pipeline) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Pipeline) GetSpec() *PipelineSpec {
	if m != nil {
		return m.Spec
	}
	return nil
}

func (m *Pipeline) GetStatuses() []*PipelineStatus {
	if m != nil {
		return m.Statuses
	}
	return nil
}

type PipelineSpec struct {
	Params               []*PipelineParam     `protobuf:"bytes,1,rep,name=params,proto3" json:"params,omitempty"`
	Tasks                map[string]*TaskSpec `protobuf:"bytes,2,rep,name=tasks,proto3" json:"tasks,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *PipelineSpec) Reset()         { *m = PipelineSpec{} }
func (m *PipelineSpec) String() string { return proto.CompactTextString(m) }
func (*PipelineSpec) ProtoMessage()    {}
func (*PipelineSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_3e7df5ab5b2e195a, []int{1}
}
func (m *PipelineSpec) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PipelineSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PipelineSpec.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PipelineSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PipelineSpec.Merge(m, src)
}
func (m *PipelineSpec) XXX_Size() int {
	return m.Size()
}
func (m *PipelineSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_PipelineSpec.DiscardUnknown(m)
}

var xxx_messageInfo_PipelineSpec proto.InternalMessageInfo

func (m *PipelineSpec) GetParams() []*PipelineParam {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *PipelineSpec) GetTasks() map[string]*TaskSpec {
	if m != nil {
		return m.Tasks
	}
	return nil
}

type PipelineParam struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Default              string   `protobuf:"bytes,2,opt,name=default,proto3" json:"default,omitempty"`
	Value                string   `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PipelineParam) Reset()         { *m = PipelineParam{} }
func (m *PipelineParam) String() string { return proto.CompactTextString(m) }
func (*PipelineParam) ProtoMessage()    {}
func (*PipelineParam) Descriptor() ([]byte, []int) {
	return fileDescriptor_3e7df5ab5b2e195a, []int{2}
}
func (m *PipelineParam) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PipelineParam) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PipelineParam.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PipelineParam) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PipelineParam.Merge(m, src)
}
func (m *PipelineParam) XXX_Size() int {
	return m.Size()
}
func (m *PipelineParam) XXX_DiscardUnknown() {
	xxx_messageInfo_PipelineParam.DiscardUnknown(m)
}

var xxx_messageInfo_PipelineParam proto.InternalMessageInfo

func (m *PipelineParam) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PipelineParam) GetDefault() string {
	if m != nil {
		return m.Default
	}
	return ""
}

func (m *PipelineParam) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type PipelineStatus struct {
	State                PipelineState    `protobuf:"varint,1,opt,name=state,proto3,enum=onka.v1alpha1.PipelineState" json:"state,omitempty"`
	Cause                string           `protobuf:"bytes,2,opt,name=cause,proto3" json:"cause,omitempty"`
	From                 *types.Timestamp `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *PipelineStatus) Reset()         { *m = PipelineStatus{} }
func (m *PipelineStatus) String() string { return proto.CompactTextString(m) }
func (*PipelineStatus) ProtoMessage()    {}
func (*PipelineStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_3e7df5ab5b2e195a, []int{3}
}
func (m *PipelineStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PipelineStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PipelineStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PipelineStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PipelineStatus.Merge(m, src)
}
func (m *PipelineStatus) XXX_Size() int {
	return m.Size()
}
func (m *PipelineStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PipelineStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PipelineStatus proto.InternalMessageInfo

func (m *PipelineStatus) GetState() PipelineState {
	if m != nil {
		return m.State
	}
	return PipelineState_PIPELINE_STATE_UNSPECIFIED
}

func (m *PipelineStatus) GetCause() string {
	if m != nil {
		return m.Cause
	}
	return ""
}

func (m *PipelineStatus) GetFrom() *types.Timestamp {
	if m != nil {
		return m.From
	}
	return nil
}

func init() {
	proto.RegisterEnum("onka.v1alpha1.PipelineState", PipelineState_name, PipelineState_value)
	proto.RegisterType((*Pipeline)(nil), "onka.v1alpha1.Pipeline")
	proto.RegisterType((*PipelineSpec)(nil), "onka.v1alpha1.PipelineSpec")
	proto.RegisterMapType((map[string]*TaskSpec)(nil), "onka.v1alpha1.PipelineSpec.TasksEntry")
	proto.RegisterType((*PipelineParam)(nil), "onka.v1alpha1.PipelineParam")
	proto.RegisterType((*PipelineStatus)(nil), "onka.v1alpha1.PipelineStatus")
}

func init() { proto.RegisterFile("pkg/proto/v1alpha1/pipeline.proto", fileDescriptor_3e7df5ab5b2e195a) }

var fileDescriptor_3e7df5ab5b2e195a = []byte{
	// 533 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0xe7, 0xfe, 0x19, 0xdd, 0x5b, 0x36, 0x45, 0xd6, 0x04, 0x51, 0x60, 0x65, 0xf4, 0x80,
	0xa6, 0x49, 0x4b, 0xb4, 0x8e, 0x03, 0xa0, 0x5d, 0x4a, 0xe7, 0xa1, 0x48, 0x23, 0x04, 0xa7, 0xbd,
	0x70, 0x99, 0xdc, 0xce, 0xed, 0xaa, 0x36, 0x4d, 0xd4, 0x38, 0x93, 0xf6, 0x15, 0xf8, 0x04, 0x7c,
	0x07, 0xbe, 0x08, 0x47, 0x2e, 0xbb, 0xa3, 0xf2, 0x45, 0x90, 0xed, 0x64, 0xa3, 0x25, 0xe5, 0x54,
	0xdb, 0xcf, 0xef, 0x79, 0xfc, 0xf6, 0x91, 0x03, 0x2f, 0xe3, 0xc9, 0xc8, 0x89, 0xe7, 0x91, 0x88,
	0x9c, 0x9b, 0x63, 0x36, 0x8d, 0xaf, 0xd9, 0xb1, 0x13, 0x8f, 0x63, 0x3e, 0x1d, 0xcf, 0xb8, 0xad,
	0xce, 0xf1, 0x76, 0x34, 0x9b, 0x30, 0x3b, 0x57, 0xad, 0x17, 0xa3, 0x28, 0x1a, 0x4d, 0xb9, 0x36,
	0xf5, 0xd3, 0xa1, 0x23, 0xc6, 0x21, 0x4f, 0x04, 0x0b, 0x63, 0xcd, 0x5b, 0x45, 0x91, 0x21, 0x17,
	0xec, 0x8a, 0x09, 0x96, 0x21, 0x7b, 0x05, 0x88, 0x60, 0xc9, 0x44, 0xcb, 0xcd, 0xef, 0x08, 0x6a,
	0x7e, 0x36, 0x04, 0x3e, 0x81, 0x5a, 0xee, 0x36, 0xd1, 0x3e, 0x3a, 0xa8, 0xb7, 0x9e, 0xda, 0x4b,
	0x13, 0xd9, 0x1f, 0x33, 0x99, 0xde, 0x83, 0xd8, 0x81, 0x4a, 0x12, 0xf3, 0x81, 0x59, 0x52, 0x86,
	0x67, 0x2b, 0x86, 0x3c, 0x3b, 0x88, 0xf9, 0x80, 0x2a, 0x10, 0xbf, 0x85, 0x5a, 0x22, 0x98, 0x48,
	0x13, 0x9e, 0x98, 0xe5, 0xfd, 0xf2, 0x41, 0xbd, 0xb5, 0xb7, 0xce, 0xa4, 0x30, 0x7a, 0x8f, 0x37,
	0xef, 0x10, 0x3c, 0xfe, 0x3b, 0x11, 0xbf, 0x86, 0xcd, 0x98, 0xcd, 0x59, 0x98, 0x98, 0x48, 0x25,
	0x3d, 0x5f, 0x93, 0xe4, 0x4b, 0x88, 0x66, 0x2c, 0x3e, 0x85, 0xaa, 0xac, 0x20, 0x31, 0x4b, 0xca,
	0xf4, 0xea, 0x3f, 0x33, 0xdb, 0x5d, 0x09, 0x92, 0x99, 0x98, 0xdf, 0x52, 0x6d, 0xb2, 0x3e, 0x03,
	0x3c, 0x1c, 0x62, 0x03, 0xca, 0x13, 0x7e, 0xab, 0xea, 0xda, 0xa2, 0x72, 0x89, 0x8f, 0xa0, 0x7a,
	0xc3, 0xa6, 0x29, 0xcf, 0x1a, 0x59, 0xad, 0x50, 0x7a, 0x55, 0x1b, 0x9a, 0x7a, 0x57, 0x7a, 0x83,
	0x9a, 0x01, 0x6c, 0x2f, 0x4d, 0x8a, 0x31, 0x54, 0x66, 0x2c, 0xe4, 0x59, 0xac, 0x5a, 0x63, 0x13,
	0x1e, 0x5d, 0xf1, 0x21, 0x4b, 0xa7, 0x42, 0x25, 0x6f, 0xd1, 0x7c, 0x8b, 0x77, 0xf3, 0x1b, 0xcb,
	0xea, 0x5c, 0x6f, 0x9a, 0x5f, 0x11, 0xec, 0x2c, 0x37, 0x89, 0x5b, 0x50, 0x95, 0x5d, 0xea, 0xdc,
	0x9d, 0xb5, 0x6d, 0x49, 0x9a, 0x53, 0x8d, 0xca, 0xf0, 0x01, 0x4b, 0x13, 0x9e, 0x5d, 0xaa, 0x37,
	0xd8, 0x86, 0xca, 0x70, 0x1e, 0x85, 0xea, 0xc6, 0x7a, 0xcb, 0xb2, 0xf5, 0x4b, 0xb5, 0xf3, 0x97,
	0x6a, 0x77, 0xf3, 0x97, 0x4a, 0x15, 0x77, 0x78, 0x87, 0x1e, 0xfe, 0xa2, 0x8a, 0xc7, 0x0d, 0xb0,
	0x7c, 0xd7, 0x27, 0x17, 0xae, 0x47, 0x2e, 0x83, 0x6e, 0xbb, 0x4b, 0x2e, 0x7b, 0x5e, 0xe0, 0x93,
	0x8e, 0x7b, 0xee, 0x92, 0x33, 0x63, 0x03, 0x5b, 0xf0, 0x64, 0x45, 0xf7, 0x89, 0x77, 0xe6, 0x7a,
	0x1f, 0x0c, 0x54, 0xa4, 0x51, 0xe2, 0xb7, 0x29, 0x31, 0x4a, 0x05, 0x1a, 0xed, 0x79, 0x9e, 0xf4,
	0x95, 0x0b, 0xb4, 0xa0, 0xd7, 0xe9, 0x90, 0x20, 0x30, 0x2a, 0x05, 0xda, 0x79, 0xdb, 0xbd, 0xe8,
	0x51, 0x62, 0x54, 0xb1, 0x09, 0xbb, 0x2b, 0x1a, 0xa1, 0xf4, 0x13, 0x35, 0x36, 0xdf, 0x9f, 0xfe,
	0x58, 0x34, 0xd0, 0xcf, 0x45, 0x03, 0xfd, 0x5a, 0x34, 0xd0, 0xb7, 0xdf, 0x8d, 0x8d, 0x2f, 0x87,
	0xa3, 0xb1, 0xb8, 0x4e, 0xfb, 0xf6, 0x20, 0x0a, 0x1d, 0x59, 0xef, 0xd1, 0x60, 0xac, 0x7e, 0x9d,
	0x7f, 0x3f, 0xc4, 0xfe, 0xa6, 0xda, 0x9f, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0x3e, 0x17, 0xb1,
	0x6e, 0x1b, 0x04, 0x00, 0x00,
}

func (m *Pipeline) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pipeline) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Pipeline) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Statuses) > 0 {
		for iNdEx := len(m.Statuses) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Statuses[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPipeline(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.Spec != nil {
		{
			size, err := m.Spec.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPipeline(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Metadata != nil {
		{
			size, err := m.Metadata.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPipeline(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PipelineSpec) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PipelineSpec) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PipelineSpec) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Tasks) > 0 {
		for k := range m.Tasks {
			v := m.Tasks[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintPipeline(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintPipeline(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintPipeline(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Params) > 0 {
		for iNdEx := len(m.Params) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Params[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPipeline(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *PipelineParam) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PipelineParam) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PipelineParam) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintPipeline(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Default) > 0 {
		i -= len(m.Default)
		copy(dAtA[i:], m.Default)
		i = encodeVarintPipeline(dAtA, i, uint64(len(m.Default)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintPipeline(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PipelineStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PipelineStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PipelineStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.From != nil {
		{
			size, err := m.From.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPipeline(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Cause) > 0 {
		i -= len(m.Cause)
		copy(dAtA[i:], m.Cause)
		i = encodeVarintPipeline(dAtA, i, uint64(len(m.Cause)))
		i--
		dAtA[i] = 0x12
	}
	if m.State != 0 {
		i = encodeVarintPipeline(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPipeline(dAtA []byte, offset int, v uint64) int {
	offset -= sovPipeline(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Pipeline) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Metadata != nil {
		l = m.Metadata.Size()
		n += 1 + l + sovPipeline(uint64(l))
	}
	if m.Spec != nil {
		l = m.Spec.Size()
		n += 1 + l + sovPipeline(uint64(l))
	}
	if len(m.Statuses) > 0 {
		for _, e := range m.Statuses {
			l = e.Size()
			n += 1 + l + sovPipeline(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PipelineSpec) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Params) > 0 {
		for _, e := range m.Params {
			l = e.Size()
			n += 1 + l + sovPipeline(uint64(l))
		}
	}
	if len(m.Tasks) > 0 {
		for k, v := range m.Tasks {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovPipeline(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovPipeline(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovPipeline(uint64(mapEntrySize))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PipelineParam) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovPipeline(uint64(l))
	}
	l = len(m.Default)
	if l > 0 {
		n += 1 + l + sovPipeline(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovPipeline(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PipelineStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.State != 0 {
		n += 1 + sovPipeline(uint64(m.State))
	}
	l = len(m.Cause)
	if l > 0 {
		n += 1 + l + sovPipeline(uint64(l))
	}
	if m.From != nil {
		l = m.From.Size()
		n += 1 + l + sovPipeline(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovPipeline(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPipeline(x uint64) (n int) {
	return sovPipeline(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Pipeline) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPipeline
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Pipeline: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pipeline: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Metadata == nil {
				m.Metadata = &Metadata{}
			}
			if err := m.Metadata.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Spec", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Spec == nil {
				m.Spec = &PipelineSpec{}
			}
			if err := m.Spec.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Statuses", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Statuses = append(m.Statuses, &PipelineStatus{})
			if err := m.Statuses[len(m.Statuses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPipeline(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPipeline
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PipelineSpec) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPipeline
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PipelineSpec: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PipelineSpec: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Params = append(m.Params, &PipelineParam{})
			if err := m.Params[len(m.Params)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tasks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Tasks == nil {
				m.Tasks = make(map[string]*TaskSpec)
			}
			var mapkey string
			var mapvalue *TaskSpec
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPipeline
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPipeline
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthPipeline
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthPipeline
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPipeline
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthPipeline
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthPipeline
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &TaskSpec{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipPipeline(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthPipeline
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Tasks[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPipeline(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPipeline
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PipelineParam) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPipeline
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PipelineParam: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PipelineParam: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Default", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Default = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPipeline(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPipeline
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PipelineStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPipeline
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PipelineStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PipelineStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= PipelineState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cause", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cause = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPipeline
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPipeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.From == nil {
				m.From = &types.Timestamp{}
			}
			if err := m.From.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPipeline(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPipeline
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPipeline(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPipeline
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPipeline
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPipeline
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPipeline
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPipeline
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPipeline        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPipeline          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPipeline = fmt.Errorf("proto: unexpected end of group")
)
