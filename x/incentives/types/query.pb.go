// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mars/incentives/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// QueryScheduleRequest is the request type for the Query/Schedule RPC method
type QueryScheduleRequest struct {
	// ID is the identifier of the incentives schedule to be queried
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *QueryScheduleRequest) Reset()         { *m = QueryScheduleRequest{} }
func (m *QueryScheduleRequest) String() string { return proto.CompactTextString(m) }
func (*QueryScheduleRequest) ProtoMessage()    {}
func (*QueryScheduleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b5ccb4babaf29c00, []int{0}
}
func (m *QueryScheduleRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryScheduleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryScheduleRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryScheduleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryScheduleRequest.Merge(m, src)
}
func (m *QueryScheduleRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryScheduleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryScheduleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryScheduleRequest proto.InternalMessageInfo

// QueryScheduleResponse is the response type for the Query/Schedule RPC method
type QueryScheduleResponse struct {
	// Schedule is the parameters of the incentives schedule
	Schedule Schedule `protobuf:"bytes,1,opt,name=schedule,proto3" json:"schedule"`
}

func (m *QueryScheduleResponse) Reset()         { *m = QueryScheduleResponse{} }
func (m *QueryScheduleResponse) String() string { return proto.CompactTextString(m) }
func (*QueryScheduleResponse) ProtoMessage()    {}
func (*QueryScheduleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b5ccb4babaf29c00, []int{1}
}
func (m *QueryScheduleResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryScheduleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryScheduleResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryScheduleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryScheduleResponse.Merge(m, src)
}
func (m *QueryScheduleResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryScheduleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryScheduleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryScheduleResponse proto.InternalMessageInfo

func (m *QueryScheduleResponse) GetSchedule() Schedule {
	if m != nil {
		return m.Schedule
	}
	return Schedule{}
}

// QuerySchedulesRequest is the request type for the Query/Schedules RPC method
type QuerySchedulesRequest struct {
	// Pagination defines an optional pagination for the request
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QuerySchedulesRequest) Reset()         { *m = QuerySchedulesRequest{} }
func (m *QuerySchedulesRequest) String() string { return proto.CompactTextString(m) }
func (*QuerySchedulesRequest) ProtoMessage()    {}
func (*QuerySchedulesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b5ccb4babaf29c00, []int{2}
}
func (m *QuerySchedulesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySchedulesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySchedulesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySchedulesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySchedulesRequest.Merge(m, src)
}
func (m *QuerySchedulesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QuerySchedulesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySchedulesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySchedulesRequest proto.InternalMessageInfo

// QueryScheduleResponse is the response type for the Query/Schedules RPC method
type QuerySchedulesResponse struct {
	// Schedule is the parameters of the incentives schedule
	Schedules []Schedule `protobuf:"bytes,1,rep,name=schedules,proto3" json:"schedules"`
	// Pagination defines the pagination in the response
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QuerySchedulesResponse) Reset()         { *m = QuerySchedulesResponse{} }
func (m *QuerySchedulesResponse) String() string { return proto.CompactTextString(m) }
func (*QuerySchedulesResponse) ProtoMessage()    {}
func (*QuerySchedulesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b5ccb4babaf29c00, []int{3}
}
func (m *QuerySchedulesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySchedulesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySchedulesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySchedulesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySchedulesResponse.Merge(m, src)
}
func (m *QuerySchedulesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QuerySchedulesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySchedulesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySchedulesResponse proto.InternalMessageInfo

func (m *QuerySchedulesResponse) GetSchedules() []Schedule {
	if m != nil {
		return m.Schedules
	}
	return nil
}

func (m *QuerySchedulesResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryScheduleRequest)(nil), "mars.incentives.v1.QueryScheduleRequest")
	proto.RegisterType((*QueryScheduleResponse)(nil), "mars.incentives.v1.QueryScheduleResponse")
	proto.RegisterType((*QuerySchedulesRequest)(nil), "mars.incentives.v1.QuerySchedulesRequest")
	proto.RegisterType((*QuerySchedulesResponse)(nil), "mars.incentives.v1.QuerySchedulesResponse")
}

func init() { proto.RegisterFile("mars/incentives/v1/query.proto", fileDescriptor_b5ccb4babaf29c00) }

var fileDescriptor_b5ccb4babaf29c00 = []byte{
	// 452 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x8b, 0xd4, 0x30,
	0x14, 0xc7, 0x9b, 0xba, 0xca, 0x6c, 0x04, 0x0f, 0x61, 0x95, 0xa5, 0xac, 0xed, 0x5a, 0x51, 0x77,
	0x17, 0x4c, 0xec, 0x7a, 0xf3, 0x20, 0xb2, 0x07, 0xc5, 0x9b, 0xd6, 0x83, 0xe0, 0x2d, 0x6d, 0x43,
	0x27, 0x30, 0xd3, 0x74, 0x9b, 0xb4, 0xb8, 0x88, 0x17, 0x4f, 0xea, 0x49, 0xf0, 0x0b, 0x0c, 0xf8,
	0x65, 0xf6, 0x38, 0xe0, 0xc5, 0x93, 0xc8, 0x8c, 0x07, 0x3f, 0x86, 0x34, 0x4d, 0x67, 0xea, 0x58,
	0x76, 0xe6, 0x16, 0xf2, 0xde, 0xff, 0xfd, 0x7f, 0xef, 0x9f, 0x40, 0x77, 0x4c, 0x0b, 0x49, 0x78,
	0x16, 0xb3, 0x4c, 0xf1, 0x8a, 0x49, 0x52, 0x05, 0xe4, 0xb4, 0x64, 0xc5, 0x19, 0xce, 0x0b, 0xa1,
	0x04, 0x42, 0x75, 0x1d, 0x2f, 0xeb, 0xb8, 0x0a, 0x9c, 0xa3, 0x58, 0xc8, 0xb1, 0x90, 0x24, 0xa2,
	0x92, 0x35, 0xcd, 0xa4, 0x0a, 0x22, 0xa6, 0x68, 0x40, 0x72, 0x9a, 0xf2, 0x8c, 0x2a, 0x2e, 0xb2,
	0x46, 0xef, 0xec, 0xa4, 0x22, 0x15, 0xfa, 0x48, 0xea, 0x93, 0xb9, 0xdd, 0x4b, 0x85, 0x48, 0x47,
	0x8c, 0xd0, 0x9c, 0x13, 0x9a, 0x65, 0x42, 0x69, 0x89, 0x34, 0xd5, 0x3e, 0x26, 0xa9, 0x44, 0xc1,
	0x9a, 0xba, 0xff, 0x00, 0xee, 0xbc, 0xac, 0x5d, 0x5f, 0xc5, 0x43, 0x96, 0x94, 0x23, 0x16, 0xb2,
	0xd3, 0x92, 0x49, 0x85, 0xae, 0x41, 0x9b, 0x27, 0xbb, 0x60, 0x1f, 0x1c, 0x6c, 0x85, 0x36, 0x4f,
	0x1e, 0x0d, 0x3e, 0x4e, 0x3c, 0xeb, 0xcf, 0xc4, 0xb3, 0xfc, 0xd7, 0xf0, 0xfa, 0x8a, 0x42, 0xe6,
	0x22, 0x93, 0x0c, 0x3d, 0x86, 0x03, 0x69, 0xee, 0xb4, 0xf0, 0xea, 0xf1, 0x1e, 0xfe, 0x7f, 0x63,
	0xdc, 0xea, 0x4e, 0xb6, 0xce, 0x7f, 0x7a, 0x56, 0xb8, 0xd0, 0xf8, 0x7c, 0x65, 0xb0, 0x6c, 0x59,
	0x9e, 0x42, 0xb8, 0xcc, 0xc2, 0x8c, 0xbe, 0x8b, 0x9b, 0xe0, 0x70, 0x1d, 0x1c, 0x6e, 0x52, 0x36,
	0xc1, 0xe1, 0x17, 0x34, 0x6d, 0xf7, 0x08, 0x3b, 0xca, 0xce, 0x0e, 0xdf, 0x00, 0xbc, 0xb1, 0xea,
	0x65, 0xb6, 0x78, 0x02, 0xb7, 0x5b, 0x22, 0xb9, 0x0b, 0xf6, 0x2f, 0x6d, 0xb8, 0xc6, 0x52, 0x84,
	0x9e, 0xfd, 0x83, 0x6b, 0x6b, 0xdc, 0x7b, 0x6b, 0x71, 0x1b, 0xfb, 0x2e, 0xef, 0xf1, 0xc4, 0x86,
	0x97, 0x35, 0x25, 0xfa, 0x0c, 0xe0, 0xa0, 0x35, 0x44, 0x07, 0x7d, 0x38, 0x7d, 0x8f, 0xe8, 0x1c,
	0x6e, 0xd0, 0xd9, 0xf8, 0xfa, 0x87, 0x1f, 0xbe, 0xff, 0xfe, 0x6a, 0xdf, 0x46, 0xb7, 0x48, 0xdf,
	0x87, 0x31, 0xdd, 0xe4, 0x1d, 0x4f, 0xde, 0xa3, 0x4f, 0x00, 0x6e, 0x2f, 0x72, 0x43, 0xeb, 0x3d,
	0xda, 0x77, 0x74, 0x8e, 0x36, 0x69, 0x35, 0x3c, 0x77, 0x34, 0x8f, 0x87, 0x6e, 0x5e, 0xc4, 0x23,
	0x4f, 0x9e, 0x9f, 0xcf, 0x5c, 0x30, 0x9d, 0xb9, 0xe0, 0xd7, 0xcc, 0x05, 0x5f, 0xe6, 0xae, 0x35,
	0x9d, 0xbb, 0xd6, 0x8f, 0xb9, 0x6b, 0xbd, 0x21, 0x29, 0x57, 0xc3, 0x32, 0xc2, 0xb1, 0x18, 0xeb,
	0x11, 0xf7, 0xf5, 0x7f, 0x8f, 0xc5, 0x88, 0x0c, 0xcb, 0x88, 0xbc, 0xed, 0x4e, 0x54, 0x67, 0x39,
	0x93, 0xd1, 0x15, 0xdd, 0xf0, 0xf0, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x36, 0x80, 0x39, 0xb7,
	0xc6, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Schedule queries an incentives schedule by identifier
	Schedule(ctx context.Context, in *QueryScheduleRequest, opts ...grpc.CallOption) (*QueryScheduleResponse, error)
	// Schedules queries all incentives schedules
	Schedules(ctx context.Context, in *QuerySchedulesRequest, opts ...grpc.CallOption) (*QuerySchedulesResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Schedule(ctx context.Context, in *QueryScheduleRequest, opts ...grpc.CallOption) (*QueryScheduleResponse, error) {
	out := new(QueryScheduleResponse)
	err := c.cc.Invoke(ctx, "/mars.incentives.v1.Query/Schedule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Schedules(ctx context.Context, in *QuerySchedulesRequest, opts ...grpc.CallOption) (*QuerySchedulesResponse, error) {
	out := new(QuerySchedulesResponse)
	err := c.cc.Invoke(ctx, "/mars.incentives.v1.Query/Schedules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Schedule queries an incentives schedule by identifier
	Schedule(context.Context, *QueryScheduleRequest) (*QueryScheduleResponse, error)
	// Schedules queries all incentives schedules
	Schedules(context.Context, *QuerySchedulesRequest) (*QuerySchedulesResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Schedule(ctx context.Context, req *QueryScheduleRequest) (*QueryScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Schedule not implemented")
}
func (*UnimplementedQueryServer) Schedules(ctx context.Context, req *QuerySchedulesRequest) (*QuerySchedulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Schedules not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Schedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Schedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mars.incentives.v1.Query/Schedule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Schedule(ctx, req.(*QueryScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Schedules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySchedulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Schedules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mars.incentives.v1.Query/Schedules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Schedules(ctx, req.(*QuerySchedulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mars.incentives.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Schedule",
			Handler:    _Query_Schedule_Handler,
		},
		{
			MethodName: "Schedules",
			Handler:    _Query_Schedules_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mars/incentives/v1/query.proto",
}

func (m *QueryScheduleRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryScheduleRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryScheduleRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryScheduleResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryScheduleResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryScheduleResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Schedule.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QuerySchedulesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySchedulesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySchedulesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QuerySchedulesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySchedulesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySchedulesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Schedules) > 0 {
		for iNdEx := len(m.Schedules) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Schedules[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryScheduleRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovQuery(uint64(m.Id))
	}
	return n
}

func (m *QueryScheduleResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Schedule.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QuerySchedulesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QuerySchedulesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Schedules) > 0 {
		for _, e := range m.Schedules {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryScheduleRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryScheduleRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryScheduleRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryScheduleResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryScheduleResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryScheduleResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Schedule", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Schedule.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QuerySchedulesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QuerySchedulesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySchedulesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QuerySchedulesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QuerySchedulesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySchedulesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Schedules", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Schedules = append(m.Schedules, Schedule{})
			if err := m.Schedules[len(m.Schedules)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
