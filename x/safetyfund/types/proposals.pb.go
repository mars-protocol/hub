// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mars/safetyfund/v1/proposals.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// SafetyFundSpendProposal details a proposal for the use of safety funds, together with how many coins
// are proposed to be spent, and to which recipient account
//
// NOTE: for now, this is just a copy of distribution module's CommunityPoolSpendProposal. in the long
// term, the goal is that the module is able to automatically detect bad debts incurred in the outposts
// and automatically distribute appropriate amount of funds, without having to go through the governance
// process.
type SafetyFundSpendProposal struct {
	// Title is the proposal's title
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Description is the proposal's description
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Recipient is a string representing the account address to which the funds shall be sent to
	Recipient string `protobuf:"bytes,3,opt,name=recipient,proto3" json:"recipient,omitempty"`
	// Amount represents the coins that shall be sent to the recipient
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,4,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func (m *SafetyFundSpendProposal) Reset()      { *m = SafetyFundSpendProposal{} }
func (*SafetyFundSpendProposal) ProtoMessage() {}
func (*SafetyFundSpendProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_263510325c894316, []int{0}
}
func (m *SafetyFundSpendProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SafetyFundSpendProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SafetyFundSpendProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SafetyFundSpendProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SafetyFundSpendProposal.Merge(m, src)
}
func (m *SafetyFundSpendProposal) XXX_Size() int {
	return m.Size()
}
func (m *SafetyFundSpendProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SafetyFundSpendProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SafetyFundSpendProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SafetyFundSpendProposal)(nil), "mars.safetyfund.v1.SafetyFundSpendProposal")
}

func init() {
	proto.RegisterFile("mars/safetyfund/v1/proposals.proto", fileDescriptor_263510325c894316)
}

var fileDescriptor_263510325c894316 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x1d, 0x0a, 0x95, 0x9a, 0x32, 0x45, 0x95, 0x08, 0x15, 0x72, 0xaa, 0x4e, 0x5d, 0x6a,
	0x13, 0xd8, 0x18, 0x8b, 0x84, 0xc4, 0x86, 0xda, 0x8d, 0x2d, 0x71, 0xdc, 0xd6, 0xa2, 0xf1, 0x59,
	0xb1, 0x53, 0xd1, 0x37, 0x60, 0x64, 0x64, 0xec, 0xcc, 0x93, 0x74, 0xec, 0xc8, 0x04, 0x28, 0x5d,
	0x78, 0x0c, 0x14, 0x27, 0x52, 0x33, 0xd9, 0xbe, 0xff, 0xf3, 0x7f, 0xbf, 0xcf, 0xee, 0x30, 0x8d,
	0x32, 0x4d, 0x75, 0x34, 0xe7, 0x66, 0x33, 0xcf, 0x65, 0x42, 0xd7, 0x21, 0x55, 0x19, 0x28, 0xd0,
	0xd1, 0x4a, 0x13, 0x95, 0x81, 0x01, 0xcf, 0x2b, 0x19, 0x72, 0x64, 0xc8, 0x3a, 0xec, 0xf7, 0x16,
	0xb0, 0x00, 0x2b, 0xd3, 0x72, 0x57, 0x91, 0x7d, 0xcc, 0x40, 0xa7, 0xa0, 0x69, 0x1c, 0x69, 0x4e,
	0xd7, 0x61, 0xcc, 0x4d, 0x14, 0x52, 0x06, 0x42, 0x56, 0xfa, 0xb0, 0x70, 0xdc, 0x8b, 0x99, 0xf5,
	0x79, 0xc8, 0x65, 0x32, 0x53, 0x5c, 0x26, 0x4f, 0x75, 0x33, 0xaf, 0xe7, 0x9e, 0x19, 0x61, 0x56,
	0xdc, 0x77, 0x06, 0xce, 0xa8, 0x33, 0xad, 0x0e, 0xde, 0xc0, 0xed, 0x26, 0x5c, 0xb3, 0x4c, 0x28,
	0x23, 0x40, 0xfa, 0x27, 0x56, 0x6b, 0x96, 0xbc, 0x2b, 0xb7, 0x93, 0x71, 0x26, 0x94, 0xe0, 0xd2,
	0xf8, 0x2d, 0xab, 0x1f, 0x0b, 0x1e, 0x73, 0xdb, 0x51, 0x0a, 0xb9, 0x34, 0xfe, 0xe9, 0xa0, 0x35,
	0xea, 0xde, 0x5c, 0x92, 0x2a, 0x22, 0x29, 0x23, 0x92, 0x3a, 0x22, 0xb9, 0x07, 0x21, 0x27, 0xd7,
	0xbb, 0xef, 0x00, 0x7d, 0xfe, 0x04, 0xa3, 0x85, 0x30, 0xcb, 0x3c, 0x26, 0x0c, 0x52, 0x5a, 0xbf,
	0xa7, 0x5a, 0xc6, 0x3a, 0x79, 0xa1, 0x66, 0xa3, 0xb8, 0xb6, 0x17, 0xf4, 0xb4, 0xb6, 0xbe, 0x3b,
	0x7f, 0xdb, 0x06, 0xe8, 0x63, 0x1b, 0xa0, 0xbf, 0x6d, 0x80, 0x26, 0x8f, 0xbb, 0x02, 0x3b, 0xfb,
	0x02, 0x3b, 0xbf, 0x05, 0x76, 0xde, 0x0f, 0x18, 0xed, 0x0f, 0x18, 0x7d, 0x1d, 0x30, 0x7a, 0xa6,
	0x0d, 0xe7, 0x72, 0xa6, 0x63, 0x3b, 0x15, 0x06, 0x2b, 0xba, 0xcc, 0x63, 0xfa, 0xda, 0xfc, 0x06,
	0xdb, 0x26, 0x6e, 0x5b, 0xe0, 0xf6, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x56, 0xff, 0x39, 0xa6,
	0x01, 0x00, 0x00,
}

func (m *SafetyFundSpendProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SafetyFundSpendProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SafetyFundSpendProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProposals(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintProposals(dAtA, i, uint64(len(m.Recipient)))
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
func (m *SafetyFundSpendProposal) Size() (n int) {
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
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovProposals(uint64(l))
	}
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovProposals(uint64(l))
		}
	}
	return n
}

func sovProposals(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProposals(x uint64) (n int) {
	return sovProposals(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SafetyFundSpendProposal) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SafetyFundSpendProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SafetyFundSpendProposal: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
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
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
