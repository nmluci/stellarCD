// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.10
// source: fileop.proto

package fileop

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StellarFileClient is the client API for StellarFile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StellarFileClient interface {
	Search(ctx context.Context, in *FileQuery, opts ...grpc.CallOption) (*Files, error)
	Download(ctx context.Context, in *Files, opts ...grpc.CallOption) (*Empty, error)
	Archive(ctx context.Context, in *FileArchive, opts ...grpc.CallOption) (*Empty, error)
}

type stellarFileClient struct {
	cc grpc.ClientConnInterface
}

func NewStellarFileClient(cc grpc.ClientConnInterface) StellarFileClient {
	return &stellarFileClient{cc}
}

func (c *stellarFileClient) Search(ctx context.Context, in *FileQuery, opts ...grpc.CallOption) (*Files, error) {
	out := new(Files)
	err := c.cc.Invoke(ctx, "/fileop.StellarFile/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stellarFileClient) Download(ctx context.Context, in *Files, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/fileop.StellarFile/Download", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stellarFileClient) Archive(ctx context.Context, in *FileArchive, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/fileop.StellarFile/Archive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StellarFileServer is the server API for StellarFile service.
// All implementations must embed UnimplementedStellarFileServer
// for forward compatibility
type StellarFileServer interface {
	Search(context.Context, *FileQuery) (*Files, error)
	Download(context.Context, *Files) (*Empty, error)
	Archive(context.Context, *FileArchive) (*Empty, error)
	mustEmbedUnimplementedStellarFileServer()
}

// UnimplementedStellarFileServer must be embedded to have forward compatible implementations.
type UnimplementedStellarFileServer struct {
}

func (UnimplementedStellarFileServer) Search(context.Context, *FileQuery) (*Files, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedStellarFileServer) Download(context.Context, *Files) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedStellarFileServer) Archive(context.Context, *FileArchive) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Archive not implemented")
}
func (UnimplementedStellarFileServer) mustEmbedUnimplementedStellarFileServer() {}

// UnsafeStellarFileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StellarFileServer will
// result in compilation errors.
type UnsafeStellarFileServer interface {
	mustEmbedUnimplementedStellarFileServer()
}

func RegisterStellarFileServer(s grpc.ServiceRegistrar, srv StellarFileServer) {
	s.RegisterService(&StellarFile_ServiceDesc, srv)
}

func _StellarFile_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StellarFileServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileop.StellarFile/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StellarFileServer).Search(ctx, req.(*FileQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _StellarFile_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Files)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StellarFileServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileop.StellarFile/Download",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StellarFileServer).Download(ctx, req.(*Files))
	}
	return interceptor(ctx, in, info, handler)
}

func _StellarFile_Archive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileArchive)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StellarFileServer).Archive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fileop.StellarFile/Archive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StellarFileServer).Archive(ctx, req.(*FileArchive))
	}
	return interceptor(ctx, in, info, handler)
}

// StellarFile_ServiceDesc is the grpc.ServiceDesc for StellarFile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StellarFile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fileop.StellarFile",
	HandlerType: (*StellarFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _StellarFile_Search_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _StellarFile_Download_Handler,
		},
		{
			MethodName: "Archive",
			Handler:    _StellarFile_Archive_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fileop.proto",
}
