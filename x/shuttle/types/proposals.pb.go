// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mars/shuttle/v1beta1/proposals.proto

package types

import (
	fmt "fmt"
	github_com_CosmWasm_wasmd_x_wasm_types "github.com/CosmWasm/wasmd/x/wasm/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ExecuteRemoteContractProposal defines a governance proposal for instructing the interchain
// account to execute a wasm smart contract on the host chain
type ExecuteRemoteContractProposal struct {
	// Title is the title of the proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Description is the description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// ConnectionId is the identifier of the IBC connection through which the wasm message is to be dispatched
	ConnectionId string `protobuf:"bytes,3,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty" yaml:"connection_id"`
	// Contract is the address of the wasm contract which is to be executed
	Contract string `protobuf:"bytes,4,opt,name=contract,proto3" json:"contract,omitempty"`
	// Msg is the execute message which is to be dispatched to the contract
	Msg github_com_CosmWasm_wasmd_x_wasm_types.RawContractMessage `protobuf:"bytes,5,opt,name=msg,proto3,casttype=github.com/CosmWasm/wasmd/x/wasm/types.RawContractMessage" json:"msg,omitempty"`
	// Funds are the coins that will be sent to the contract during the execution
	Funds github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,6,rep,name=funds,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"funds"`
}

func (m *ExecuteRemoteContractProposal) Reset()      { *m = ExecuteRemoteContractProposal{} }
func (*ExecuteRemoteContractProposal) ProtoMessage() {}
func (*ExecuteRemoteContractProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_b64d4c09e797fcb3, []int{0}
}
func (m *ExecuteRemoteContractProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExecuteRemoteContractProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExecuteRemoteContractProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExecuteRemoteContractProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteRemoteContractProposal.Merge(m, src)
}
func (m *ExecuteRemoteContractProposal) XXX_Size() int {
	return m.Size()
}
func (m *ExecuteRemoteContractProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteRemoteContractProposal.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteRemoteContractProposal proto.InternalMessageInfo

// MigrateRemoteContractProposal defines a governance proposal for instructing the interchain
// account to migrate a wasm smart contract on the host chain
type MigrateRemoteContractProposal struct {
	// Title is the title of the proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Description is the description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// ConnectionId is the identifier of the IBC connection through which the wasm message is to be dispatched
	ConnectionId string `protobuf:"bytes,3,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty" yaml:"connection_id"`
	// Contract is the address of the wasm contract which is to be migrated
	Contract string `protobuf:"bytes,4,opt,name=contract,proto3" json:"contract,omitempty"`
	// CodeId is the identifier of the wasm code to which the contract is to be migrated
	CodeId uint64 `protobuf:"varint,5,opt,name=code_id,json=codeId,proto3" json:"code_id,omitempty" yaml:"code_id"`
	// Msg is the migration message which is to be executed during the migration process
	Msg github_com_CosmWasm_wasmd_x_wasm_types.RawContractMessage `protobuf:"bytes,6,opt,name=msg,proto3,casttype=github.com/CosmWasm/wasmd/x/wasm/types.RawContractMessage" json:"msg,omitempty"`
}

func (m *MigrateRemoteContractProposal) Reset()      { *m = MigrateRemoteContractProposal{} }
func (*MigrateRemoteContractProposal) ProtoMessage() {}
func (*MigrateRemoteContractProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_b64d4c09e797fcb3, []int{1}
}
func (m *MigrateRemoteContractProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MigrateRemoteContractProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MigrateRemoteContractProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MigrateRemoteContractProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MigrateRemoteContractProposal.Merge(m, src)
}
func (m *MigrateRemoteContractProposal) XXX_Size() int {
	return m.Size()
}
func (m *MigrateRemoteContractProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_MigrateRemoteContractProposal.DiscardUnknown(m)
}

var xxx_messageInfo_MigrateRemoteContractProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ExecuteRemoteContractProposal)(nil), "mars.shuttle.v1beta1.ExecuteRemoteContractProposal")
	proto.RegisterType((*MigrateRemoteContractProposal)(nil), "mars.shuttle.v1beta1.MigrateRemoteContractProposal")
}

func init() {
	proto.RegisterFile("mars/shuttle/v1beta1/proposals.proto", fileDescriptor_b64d4c09e797fcb3)
}

var fileDescriptor_b64d4c09e797fcb3 = []byte{
	// 481 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x53, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xb6, 0x93, 0x26, 0xc0, 0x35, 0x30, 0x58, 0x19, 0x4c, 0xa4, 0xda, 0x91, 0xc5, 0x10, 0x09,
	0xea, 0xa3, 0x30, 0x51, 0xa9, 0x8b, 0x23, 0x90, 0x3a, 0x54, 0x20, 0x2f, 0x48, 0x2c, 0xe8, 0x7c,
	0xbe, 0x3a, 0x16, 0x3e, 0x3f, 0xcb, 0xef, 0x4c, 0xd3, 0x7f, 0xc0, 0xc8, 0xc8, 0x98, 0x99, 0x1f,
	0xc1, 0xdc, 0xb1, 0x23, 0x53, 0x40, 0xc9, 0xc2, 0x5c, 0x31, 0x31, 0xa1, 0xb3, 0x9d, 0x36, 0xdd,
	0x59, 0x98, 0xee, 0x7d, 0xf7, 0xbe, 0xf7, 0xbd, 0xf3, 0xf7, 0xfc, 0xc8, 0x23, 0xc9, 0x4a, 0xa4,
	0x38, 0xab, 0x94, 0xca, 0x04, 0xfd, 0x78, 0x10, 0x09, 0xc5, 0x0e, 0x68, 0x51, 0x42, 0x01, 0xc8,
	0x32, 0xf4, 0x8b, 0x12, 0x14, 0x58, 0x43, 0xcd, 0xf2, 0x5b, 0x96, 0xdf, 0xb2, 0x46, 0x0e, 0x07,
	0x94, 0x80, 0x34, 0x62, 0x78, 0x53, 0xca, 0x21, 0xcd, 0x9b, 0xaa, 0xd1, 0x30, 0x81, 0x04, 0xea,
	0x90, 0xea, 0xa8, 0xbd, 0x75, 0x13, 0x80, 0x24, 0x13, 0xb4, 0x46, 0x51, 0x75, 0x4a, 0x55, 0x2a,
	0x05, 0x2a, 0x26, 0x8b, 0x86, 0xe0, 0xfd, 0xee, 0x90, 0xbd, 0x97, 0x73, 0xc1, 0x2b, 0x25, 0x42,
	0x21, 0x41, 0x89, 0x29, 0xe4, 0xaa, 0x64, 0x5c, 0xbd, 0x69, 0x5f, 0x65, 0x0d, 0x49, 0x4f, 0xa5,
	0x2a, 0x13, 0xb6, 0x39, 0x36, 0x27, 0xf7, 0xc2, 0x06, 0x58, 0x63, 0xb2, 0x1b, 0x0b, 0xe4, 0x65,
	0x5a, 0xa8, 0x14, 0x72, 0xbb, 0x53, 0xe7, 0xb6, 0xaf, 0xac, 0x23, 0x72, 0x9f, 0x43, 0x9e, 0x0b,
	0xae, 0xd1, 0xfb, 0x34, 0xb6, 0xbb, 0x9a, 0x13, 0xd8, 0x57, 0x4b, 0x77, 0x78, 0xce, 0x64, 0x76,
	0xe8, 0xdd, 0x4a, 0x7b, 0xe1, 0xe0, 0x06, 0x1f, 0xc7, 0xd6, 0x88, 0xdc, 0xe5, 0xed, 0x53, 0xec,
	0x9d, 0x5a, 0xfd, 0x1a, 0x5b, 0xaf, 0x49, 0x57, 0x62, 0x62, 0xf7, 0xc6, 0xe6, 0x64, 0x10, 0x1c,
	0xfd, 0x59, 0xba, 0x2f, 0x92, 0x54, 0xcd, 0xaa, 0xc8, 0xe7, 0x20, 0xe9, 0x14, 0x50, 0xbe, 0x65,
	0x28, 0xe9, 0x19, 0x43, 0x19, 0xd3, 0x79, 0x7d, 0x52, 0x75, 0x5e, 0x08, 0xf4, 0x43, 0x76, 0xb6,
	0xf9, 0xbe, 0x13, 0x81, 0xc8, 0x12, 0x11, 0x6a, 0x25, 0x8b, 0x91, 0xde, 0x69, 0x95, 0xc7, 0x68,
	0xf7, 0xc7, 0xdd, 0xc9, 0xee, 0xb3, 0x87, 0x7e, 0x63, 0xb6, 0xaf, 0xcd, 0xde, 0x4c, 0xc0, 0x9f,
	0x42, 0x9a, 0x07, 0x4f, 0x2f, 0x96, 0xae, 0xf1, 0xf5, 0x87, 0x3b, 0xd9, 0xea, 0xd8, 0x4e, 0xa6,
	0x39, 0xf6, 0x31, 0xfe, 0xd0, 0x76, 0xd3, 0x05, 0x18, 0x36, 0xca, 0x87, 0x83, 0x4f, 0x0b, 0xd7,
	0xf8, 0xb2, 0x70, 0x8d, 0x5f, 0x0b, 0xd7, 0xf0, 0xbe, 0x75, 0xc8, 0xde, 0x49, 0x9a, 0x94, 0xec,
	0x7f, 0xb2, 0xfd, 0x31, 0xb9, 0xc3, 0x21, 0x16, 0x5a, 0x54, 0x5b, 0xbf, 0x13, 0x58, 0x57, 0x4b,
	0xf7, 0xc1, 0x46, 0xb4, 0x4e, 0x78, 0x61, 0x5f, 0x47, 0xc7, 0xf1, 0x66, 0x46, 0xfd, 0x7f, 0x35,
	0xa3, 0xdb, 0x06, 0x06, 0xaf, 0x2e, 0x56, 0x8e, 0x79, 0xb9, 0x72, 0xcc, 0x9f, 0x2b, 0xc7, 0xfc,
	0xbc, 0x76, 0x8c, 0xcb, 0xb5, 0x63, 0x7c, 0x5f, 0x3b, 0xc6, 0xbb, 0x27, 0x5b, 0x7d, 0xf4, 0x26,
	0xed, 0xd7, 0x3f, 0x3a, 0x87, 0x8c, 0xce, 0xaa, 0x88, 0xce, 0xaf, 0xd7, 0xaf, 0xee, 0x16, 0xf5,
	0xeb, 0xec, 0xf3, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x42, 0x44, 0x5f, 0x95, 0x9b, 0x03, 0x00,
	0x00,
}

func (m *ExecuteRemoteContractProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExecuteRemoteContractProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExecuteRemoteContractProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Funds) > 0 {
		for iNdEx := len(m.Funds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Funds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProposals(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ConnectionId) > 0 {
		i -= len(m.ConnectionId)
		copy(dAtA[i:], m.ConnectionId)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.ConnectionId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MigrateRemoteContractProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MigrateRemoteContractProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MigrateRemoteContractProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x32
	}
	if m.CodeId != 0 {
		i = encodeVarintProposals(dAtA, i, uint64(m.CodeId))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ConnectionId) > 0 {
		i -= len(m.ConnectionId)
		copy(dAtA[i:], m.ConnectionId)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.ConnectionId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProposals(dAtA []byte, offset int, v uint64) int {
	offset -= sovProposals(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ExecuteRemoteContractProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	l = len(m.ConnectionId)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	if len(m.Funds) > 0 {
		for _, e := range m.Funds {
			l = e.Size()
			n += 1 + l + sovProposals(uint64(l))
		}
	}
	return n
}

func (m *MigrateRemoteContractProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	l = len(m.ConnectionId)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	if m.CodeId != 0 {
		n += 1 + sovProposals(uint64(m.CodeId))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	return n
}

func sovProposals(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProposals(x uint64) (n int) {
	return sovProposals(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ExecuteRemoteContractProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposals
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
			return fmt.Errorf("proto: ExecuteRemoteContractProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExecuteRemoteContractProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = append(m.Msg[:0], dAtA[iNdEx:postIndex]...)
			if m.Msg == nil {
				m.Msg = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Funds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Funds = append(m.Funds, types.Coin{})
			if err := m.Funds[len(m.Funds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposals(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposals
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MigrateRemoteContractProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposals
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
			return fmt.Errorf("proto: MigrateRemoteContractProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MigrateRemoteContractProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
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
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodeId", wireType)
			}
			m.CodeId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CodeId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposals
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProposals
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProposals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = append(m.Msg[:0], dAtA[iNdEx:postIndex]...)
			if m.Msg == nil {
				m.Msg = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposals(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposals
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProposals(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProposals
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
					return 0, ErrIntOverflowProposals
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
					return 0, ErrIntOverflowProposals
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
				return 0, ErrInvalidLengthProposals
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProposals
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProposals
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProposals        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProposals          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProposals = fmt.Errorf("proto: unexpected end of group")
)
