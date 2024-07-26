// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: internal/services/proto/gophkeeper.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	GophKeeper_Login_FullMethodName      = "/goyandexgophkeeper.GophKeeper/Login"
	GophKeeper_AddItem_FullMethodName    = "/goyandexgophkeeper.GophKeeper/AddItem"
	GophKeeper_DelItem_FullMethodName    = "/goyandexgophkeeper.GophKeeper/DelItem"
	GophKeeper_UpdateItem_FullMethodName = "/goyandexgophkeeper.GophKeeper/UpdateItem"
	GophKeeper_CreateUser_FullMethodName = "/goyandexgophkeeper.GophKeeper/CreateUser"
	GophKeeper_ListItems_FullMethodName  = "/goyandexgophkeeper.GophKeeper/ListItems"
)

// GophKeeperClient is the client API for GophKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error)
	DelItem(ctx context.Context, in *DelItemRequest, opts ...grpc.CallOption) (*DelItemResponse, error)
	UpdateItem(ctx context.Context, in *UpdateItemRequest, opts ...grpc.CallOption) (*UpdateItemResponse, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (*ListItemsResponse, error)
}

type gophKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperClient(cc grpc.ClientConnInterface) GophKeeperClient {
	return &gophKeeperClient{cc}
}

func (c *gophKeeperClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddItemResponse)
	err := c.cc.Invoke(ctx, GophKeeper_AddItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) DelItem(ctx context.Context, in *DelItemRequest, opts ...grpc.CallOption) (*DelItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DelItemResponse)
	err := c.cc.Invoke(ctx, GophKeeper_DelItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) UpdateItem(ctx context.Context, in *UpdateItemRequest, opts ...grpc.CallOption) (*UpdateItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateItemResponse)
	err := c.cc.Invoke(ctx, GophKeeper_UpdateItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, GophKeeper_CreateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (*ListItemsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListItemsResponse)
	err := c.cc.Invoke(ctx, GophKeeper_ListItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServer is the server API for GophKeeper service.
// All implementations must embed UnimplementedGophKeeperServer
// for forward compatibility
type GophKeeperServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error)
	DelItem(context.Context, *DelItemRequest) (*DelItemResponse, error)
	UpdateItem(context.Context, *UpdateItemRequest) (*UpdateItemResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	ListItems(context.Context, *ListItemsRequest) (*ListItemsResponse, error)
	mustEmbedUnimplementedGophKeeperServer()
}

// UnimplementedGophKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedGophKeeperServer struct {
}

func (UnimplementedGophKeeperServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedGophKeeperServer) AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedGophKeeperServer) DelItem(context.Context, *DelItemRequest) (*DelItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelItem not implemented")
}
func (UnimplementedGophKeeperServer) UpdateItem(context.Context, *UpdateItemRequest) (*UpdateItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateItem not implemented")
}
func (UnimplementedGophKeeperServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedGophKeeperServer) ListItems(context.Context, *ListItemsRequest) (*ListItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListItems not implemented")
}
func (UnimplementedGophKeeperServer) mustEmbedUnimplementedGophKeeperServer() {}

// UnsafeGophKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServer will
// result in compilation errors.
type UnsafeGophKeeperServer interface {
	mustEmbedUnimplementedGophKeeperServer()
}

func RegisterGophKeeperServer(s grpc.ServiceRegistrar, srv GophKeeperServer) {
	s.RegisterService(&GophKeeper_ServiceDesc, srv)
}

func _GophKeeper_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_AddItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).AddItem(ctx, req.(*AddItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_DelItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).DelItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_DelItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).DelItem(ctx, req.(*DelItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_UpdateItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).UpdateItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_UpdateItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).UpdateItem(ctx, req.(*UpdateItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_ListItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).ListItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_ListItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).ListItems(ctx, req.(*ListItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeper_ServiceDesc is the grpc.ServiceDesc for GophKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "goyandexgophkeeper.GophKeeper",
	HandlerType: (*GophKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _GophKeeper_Login_Handler,
		},
		{
			MethodName: "AddItem",
			Handler:    _GophKeeper_AddItem_Handler,
		},
		{
			MethodName: "DelItem",
			Handler:    _GophKeeper_DelItem_Handler,
		},
		{
			MethodName: "UpdateItem",
			Handler:    _GophKeeper_UpdateItem_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _GophKeeper_CreateUser_Handler,
		},
		{
			MethodName: "ListItems",
			Handler:    _GophKeeper_ListItems_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/services/proto/gophkeeper.proto",
}

const (
	FilseService_Upload_FullMethodName     = "/goyandexgophkeeper.FilseService/Upload"
	FilseService_Download_FullMethodName   = "/goyandexgophkeeper.FilseService/Download"
	FilseService_AddItem_FullMethodName    = "/goyandexgophkeeper.FilseService/AddItem"
	FilseService_UpdateItem_FullMethodName = "/goyandexgophkeeper.FilseService/UpdateItem"
)

// FilseServiceClient is the client API for FilseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FilseServiceClient interface {
	Upload(ctx context.Context, opts ...grpc.CallOption) (FilseService_UploadClient, error)
	Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (FilseService_DownloadClient, error)
	AddItem(ctx context.Context, in *AddItemWithTokenRequest, opts ...grpc.CallOption) (*AddItemResponse, error)
	UpdateItem(ctx context.Context, in *UpdateItemWithTokenRequest, opts ...grpc.CallOption) (*UpdateItemResponse, error)
}

type filseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFilseServiceClient(cc grpc.ClientConnInterface) FilseServiceClient {
	return &filseServiceClient{cc}
}

func (c *filseServiceClient) Upload(ctx context.Context, opts ...grpc.CallOption) (FilseService_UploadClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FilseService_ServiceDesc.Streams[0], FilseService_Upload_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &filseServiceUploadClient{ClientStream: stream}
	return x, nil
}

type FilseService_UploadClient interface {
	Send(*FileUploadRequest) error
	CloseAndRecv() (*FileUploadResponse, error)
	grpc.ClientStream
}

type filseServiceUploadClient struct {
	grpc.ClientStream
}

func (x *filseServiceUploadClient) Send(m *FileUploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filseServiceUploadClient) CloseAndRecv() (*FileUploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileUploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filseServiceClient) Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (FilseService_DownloadClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FilseService_ServiceDesc.Streams[1], FilseService_Download_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &filseServiceDownloadClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FilseService_DownloadClient interface {
	Recv() (*FileDownloadResponse, error)
	grpc.ClientStream
}

type filseServiceDownloadClient struct {
	grpc.ClientStream
}

func (x *filseServiceDownloadClient) Recv() (*FileDownloadResponse, error) {
	m := new(FileDownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filseServiceClient) AddItem(ctx context.Context, in *AddItemWithTokenRequest, opts ...grpc.CallOption) (*AddItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddItemResponse)
	err := c.cc.Invoke(ctx, FilseService_AddItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filseServiceClient) UpdateItem(ctx context.Context, in *UpdateItemWithTokenRequest, opts ...grpc.CallOption) (*UpdateItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateItemResponse)
	err := c.cc.Invoke(ctx, FilseService_UpdateItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilseServiceServer is the server API for FilseService service.
// All implementations must embed UnimplementedFilseServiceServer
// for forward compatibility
type FilseServiceServer interface {
	Upload(FilseService_UploadServer) error
	Download(*FileDownloadRequest, FilseService_DownloadServer) error
	AddItem(context.Context, *AddItemWithTokenRequest) (*AddItemResponse, error)
	UpdateItem(context.Context, *UpdateItemWithTokenRequest) (*UpdateItemResponse, error)
	mustEmbedUnimplementedFilseServiceServer()
}

// UnimplementedFilseServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFilseServiceServer struct {
}

func (UnimplementedFilseServiceServer) Upload(FilseService_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedFilseServiceServer) Download(*FileDownloadRequest, FilseService_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedFilseServiceServer) AddItem(context.Context, *AddItemWithTokenRequest) (*AddItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedFilseServiceServer) UpdateItem(context.Context, *UpdateItemWithTokenRequest) (*UpdateItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateItem not implemented")
}
func (UnimplementedFilseServiceServer) mustEmbedUnimplementedFilseServiceServer() {}

// UnsafeFilseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FilseServiceServer will
// result in compilation errors.
type UnsafeFilseServiceServer interface {
	mustEmbedUnimplementedFilseServiceServer()
}

func RegisterFilseServiceServer(s grpc.ServiceRegistrar, srv FilseServiceServer) {
	s.RegisterService(&FilseService_ServiceDesc, srv)
}

func _FilseService_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilseServiceServer).Upload(&filseServiceUploadServer{ServerStream: stream})
}

type FilseService_UploadServer interface {
	SendAndClose(*FileUploadResponse) error
	Recv() (*FileUploadRequest, error)
	grpc.ServerStream
}

type filseServiceUploadServer struct {
	grpc.ServerStream
}

func (x *filseServiceUploadServer) SendAndClose(m *FileUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filseServiceUploadServer) Recv() (*FileUploadRequest, error) {
	m := new(FileUploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FilseService_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileDownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FilseServiceServer).Download(m, &filseServiceDownloadServer{ServerStream: stream})
}

type FilseService_DownloadServer interface {
	Send(*FileDownloadResponse) error
	grpc.ServerStream
}

type filseServiceDownloadServer struct {
	grpc.ServerStream
}

func (x *filseServiceDownloadServer) Send(m *FileDownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _FilseService_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemWithTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilseServiceServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilseService_AddItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilseServiceServer).AddItem(ctx, req.(*AddItemWithTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FilseService_UpdateItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateItemWithTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilseServiceServer).UpdateItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilseService_UpdateItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilseServiceServer).UpdateItem(ctx, req.(*UpdateItemWithTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FilseService_ServiceDesc is the grpc.ServiceDesc for FilseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FilseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "goyandexgophkeeper.FilseService",
	HandlerType: (*FilseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddItem",
			Handler:    _FilseService_AddItem_Handler,
		},
		{
			MethodName: "UpdateItem",
			Handler:    _FilseService_UpdateItem_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _FilseService_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _FilseService_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/services/proto/gophkeeper.proto",
}
