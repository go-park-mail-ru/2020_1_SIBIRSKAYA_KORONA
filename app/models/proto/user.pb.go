// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

func init() {
	proto.RegisterFile("user.proto", fileDescriptor_116e343673f7ffaf)
}

var fileDescriptor_116e343673f7ffaf = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xbc, 0xb9, 0xa9, 0xc5, 0xc5,
	0x89, 0xe9, 0xa9, 0x10, 0x51, 0x29, 0xe9, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0x7d, 0x30, 0x2f,
	0xa9, 0x34, 0x4d, 0x3f, 0x35, 0xb7, 0xa0, 0xa4, 0x12, 0x22, 0x69, 0xb4, 0x99, 0x89, 0x8b, 0x25,
	0x14, 0x68, 0x82, 0x90, 0x26, 0x17, 0x9b, 0x73, 0x51, 0x6a, 0x62, 0x49, 0xaa, 0x10, 0x3f, 0x44,
	0x4a, 0x0f, 0x24, 0xec, 0x0b, 0x34, 0x49, 0x8a, 0x0f, 0x26, 0x90, 0x99, 0x02, 0xe2, 0x0b, 0x69,
	0x71, 0xb1, 0xbb, 0xa7, 0x96, 0x38, 0x55, 0x7a, 0xba, 0x08, 0xa1, 0x49, 0x49, 0xa1, 0xeb, 0x15,
	0x32, 0xe5, 0xe2, 0x05, 0xab, 0xf5, 0xcb, 0x4c, 0xce, 0xce, 0x4b, 0xcc, 0x4d, 0x15, 0x12, 0x86,
	0xaa, 0x80, 0x09, 0x60, 0xd7, 0x66, 0xc6, 0xc5, 0xeb, 0x9c, 0x91, 0x9a, 0x9c, 0x1d, 0x90, 0x58,
	0x5c, 0x5c, 0x9e, 0x5f, 0x94, 0x22, 0x24, 0x02, 0x55, 0x01, 0x17, 0x45, 0xd1, 0xe7, 0x94, 0x9f,
	0x9f, 0x03, 0xd6, 0x67, 0xcc, 0xc5, 0x16, 0x5a, 0x90, 0x02, 0xf2, 0x85, 0x20, 0xcc, 0x48, 0x30,
	0x17, 0xac, 0x5a, 0x4c, 0x0f, 0x12, 0x12, 0x7a, 0xb0, 0x90, 0xd0, 0x73, 0x05, 0x85, 0x84, 0x90,
	0x01, 0x17, 0x9b, 0x4b, 0x6a, 0x4e, 0x2a, 0x50, 0x13, 0xba, 0x77, 0x70, 0xe8, 0x48, 0x62, 0x03,
	0xf3, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa8, 0xd9, 0x5c, 0x10, 0x7d, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserClient interface {
	Create(ctx context.Context, in *UserMess, opts ...grpc.CallOption) (*UidMess, error)
	GetByID(ctx context.Context, in *UidMess, opts ...grpc.CallOption) (*UserMess, error)
	GetByNickname(ctx context.Context, in *NicknameMess, opts ...grpc.CallOption) (*UserMess, error)
	CheckPassword(ctx context.Context, in *CheckPassMess, opts ...grpc.CallOption) (*BoolMess, error)
	Update(ctx context.Context, in *UpdateMess, opts ...grpc.CallOption) (*empty.Empty, error)
	Delete(ctx context.Context, in *UidMess, opts ...grpc.CallOption) (*empty.Empty, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) Create(ctx context.Context, in *UserMess, opts ...grpc.CallOption) (*UidMess, error) {
	out := new(UidMess)
	err := c.cc.Invoke(ctx, "/proto.User/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetByID(ctx context.Context, in *UidMess, opts ...grpc.CallOption) (*UserMess, error) {
	out := new(UserMess)
	err := c.cc.Invoke(ctx, "/proto.User/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetByNickname(ctx context.Context, in *NicknameMess, opts ...grpc.CallOption) (*UserMess, error) {
	out := new(UserMess)
	err := c.cc.Invoke(ctx, "/proto.User/GetByNickname", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CheckPassword(ctx context.Context, in *CheckPassMess, opts ...grpc.CallOption) (*BoolMess, error) {
	out := new(BoolMess)
	err := c.cc.Invoke(ctx, "/proto.User/CheckPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Update(ctx context.Context, in *UpdateMess, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.User/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Delete(ctx context.Context, in *UidMess, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.User/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
type UserServer interface {
	Create(context.Context, *UserMess) (*UidMess, error)
	GetByID(context.Context, *UidMess) (*UserMess, error)
	GetByNickname(context.Context, *NicknameMess) (*UserMess, error)
	CheckPassword(context.Context, *CheckPassMess) (*BoolMess, error)
	Update(context.Context, *UpdateMess) (*empty.Empty, error)
	Delete(context.Context, *UidMess) (*empty.Empty, error)
}

// UnimplementedUserServer can be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (*UnimplementedUserServer) Create(ctx context.Context, req *UserMess) (*UidMess, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedUserServer) GetByID(ctx context.Context, req *UidMess) (*UserMess, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (*UnimplementedUserServer) GetByNickname(ctx context.Context, req *NicknameMess) (*UserMess, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByNickname not implemented")
}
func (*UnimplementedUserServer) CheckPassword(ctx context.Context, req *CheckPassMess) (*BoolMess, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPassword not implemented")
}
func (*UnimplementedUserServer) Update(ctx context.Context, req *UpdateMess) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedUserServer) Delete(ctx context.Context, req *UidMess) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func RegisterUserServer(s *grpc.Server, srv UserServer) {
	s.RegisterService(&_User_serviceDesc, srv)
}

func _User_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMess)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.User/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Create(ctx, req.(*UserMess))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UidMess)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.User/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetByID(ctx, req.(*UidMess))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetByNickname_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NicknameMess)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetByNickname(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.User/GetByNickname",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetByNickname(ctx, req.(*NicknameMess))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CheckPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPassMess)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CheckPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.User/CheckPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CheckPassword(ctx, req.(*CheckPassMess))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMess)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.User/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Update(ctx, req.(*UpdateMess))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UidMess)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.User/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Delete(ctx, req.(*UidMess))
	}
	return interceptor(ctx, in, info, handler)
}

var _User_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _User_Create_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _User_GetByID_Handler,
		},
		{
			MethodName: "GetByNickname",
			Handler:    _User_GetByNickname_Handler,
		},
		{
			MethodName: "CheckPassword",
			Handler:    _User_CheckPassword_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _User_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _User_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}