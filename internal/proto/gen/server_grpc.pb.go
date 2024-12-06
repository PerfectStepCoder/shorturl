// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: server.proto

package shorter

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Shorter_ShorterURL_FullMethodName = "/shorter.Shorter/ShorterURL"
	Shorter_Login_FullMethodName      = "/shorter.Shorter/Login"
	Shorter_Stats_FullMethodName      = "/shorter.Shorter/Stats"
)

// ShorterClient is the client API for Shorter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShorterClient interface {
	ShorterURL(ctx context.Context, in *RequestFullURL, opts ...grpc.CallOption) (*ResponseShortURL, error)
	Login(ctx context.Context, in *RequestLogin, opts ...grpc.CallOption) (*ResponseJWT, error)
	Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ResponseStats, error)
}

type shorterClient struct {
	cc grpc.ClientConnInterface
}

func NewShorterClient(cc grpc.ClientConnInterface) ShorterClient {
	return &shorterClient{cc}
}

func (c *shorterClient) ShorterURL(ctx context.Context, in *RequestFullURL, opts ...grpc.CallOption) (*ResponseShortURL, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseShortURL)
	err := c.cc.Invoke(ctx, Shorter_ShorterURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shorterClient) Login(ctx context.Context, in *RequestLogin, opts ...grpc.CallOption) (*ResponseJWT, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseJWT)
	err := c.cc.Invoke(ctx, Shorter_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shorterClient) Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ResponseStats, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseStats)
	err := c.cc.Invoke(ctx, Shorter_Stats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShorterServer is the server API for Shorter service.
// All implementations must embed UnimplementedShorterServer
// for forward compatibility.
type ShorterServer interface {
	ShorterURL(context.Context, *RequestFullURL) (*ResponseShortURL, error)
	Login(context.Context, *RequestLogin) (*ResponseJWT, error)
	Stats(context.Context, *emptypb.Empty) (*ResponseStats, error)
	mustEmbedUnimplementedShorterServer()
}

// UnimplementedShorterServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedShorterServer struct{}

func (UnimplementedShorterServer) ShorterURL(context.Context, *RequestFullURL) (*ResponseShortURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShorterURL not implemented")
}
func (UnimplementedShorterServer) Login(context.Context, *RequestLogin) (*ResponseJWT, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedShorterServer) Stats(context.Context, *emptypb.Empty) (*ResponseStats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedShorterServer) mustEmbedUnimplementedShorterServer() {}
func (UnimplementedShorterServer) testEmbeddedByValue()                 {}

// UnsafeShorterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShorterServer will
// result in compilation errors.
type UnsafeShorterServer interface {
	mustEmbedUnimplementedShorterServer()
}

func RegisterShorterServer(s grpc.ServiceRegistrar, srv ShorterServer) {
	// If the following call pancis, it indicates UnimplementedShorterServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Shorter_ServiceDesc, srv)
}

func _Shorter_ShorterURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestFullURL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShorterServer).ShorterURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorter_ShorterURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShorterServer).ShorterURL(ctx, req.(*RequestFullURL))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shorter_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestLogin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShorterServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorter_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShorterServer).Login(ctx, req.(*RequestLogin))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shorter_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShorterServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shorter_Stats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShorterServer).Stats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Shorter_ServiceDesc is the grpc.ServiceDesc for Shorter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shorter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shorter.Shorter",
	HandlerType: (*ShorterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShorterURL",
			Handler:    _Shorter_ShorterURL_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Shorter_Login_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _Shorter_Stats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
