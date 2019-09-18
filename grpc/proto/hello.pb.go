package proto

import (
	"fmt"
	proto1 "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"math"
)
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

const _ = proto1.ProtoPackageIsVersion2

type HelloWorldRequest struct {
	Referer string `protobuf:"bytes,1,opt,name=referer" json:"referer"`
}

func (m *HelloWorldRequest)Reset()  {
	*m = HelloWorldRequest{}
}
func (m *HelloWorldRequest)String()string{
	return proto1.CompactTextString(m)
}

func (*HelloWorldRequest)ProtoMessage(){}
func (*HelloWorldRequest)Descriptor()([]byte,[]int){
	return fileDescriptor0,[]int{0}
}
func (m *HelloWorldRequest)GetReferer()string{
	if m!=nil{
		return m.Referer
	}
	return ""
}

type HelloWorldResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *HelloWorldResponse)Reset(){*m = HelloWorldResponse{}}
func (m *HelloWorldResponse) String() string            { return proto1.CompactTextString(m) }
func (*HelloWorldResponse) ProtoMessage()               {}
func (*HelloWorldResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HelloWorldResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init(){
	proto1.RegisterType((*HelloWorldRequest)(nil),"proto.HelloWorldRequest")
	proto1.RegisterType((*HelloWorldResponse)(nil),"proto.HelloWorldResponse")
}

var _ context.Context
var _ grpc.ClientConn

const _ = grpc.SupportPackageIsVersion4

type HelloWorldClient interface {
	SayHelloWorld(ctx context.Context, in *HelloWorldRequest, opts ...grpc.CallOption) (*HelloWorldResponse, error)
}

type helloWorldClient struct {
	cc *grpc.ClientConn
}

func NewHelloWorldClient(cc *grpc.ClientConn) HelloWorldClient {
	return &helloWorldClient{cc}
}

func (c *helloWorldClient)SayHelloWorld(ctx context.Context,in *HelloWorldClient, opts ...grpc.CallOption) (*HelloWorldResponse, error){
	out := new(HelloWorldResponse)
	err := grpc.Invoke(ctx, "/proto.HelloWorld/SayHelloWorld", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type HelloWorldServer interface {
	SayHelloWorld(context.Context, *HelloWorldRequest) (*HelloWorldResponse, error)
}

func RegisterHelloWorldServer(s *grpc.Server,srv HelloWorldServer){
	s.RegisterService(&_HelloWorld_serviceDesc,srv)
}

func _HelloWorld_SayHelloWorld_Handler(srv interface{},ctx context.Context,dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error){
	in :=new(HelloWorldRequest)
	if err:=dec(in);err!=nil{
		return nil,err
	}
	if interceptor == nil{
		return srv.(HelloWorldServer).SayHelloWorld(ctx,in)
	}
	info := &grpc.UnaryServerInfo{
		Server:srv,
		FullMethod:"/proto.HelloWorld/SayHelloWorld",
	}
	handler := func(ctx context.Context,req interface{})(interface{},error){
		return srv.(HelloWorldServer).SayHelloWorld(ctx,req.(*HelloWorldRequest))
	}
	return interceptor(ctx,in,info,handler)
}
var _HelloWorld_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.HelloWorld",
	HandlerType: (*HelloWorldServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHelloWorld",
			Handler:    _HelloWorld_SayHelloWorld_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hello.proto",
}

func init() { proto1.RegisterFile("hello.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 184 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0x32, 0xe9, 0xf9, 0xf9,
	0xe9, 0x39, 0xa9, 0xfa, 0x89, 0x05, 0x99, 0xfa, 0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x99,
	0xf9, 0x79, 0xc5, 0x10, 0x45, 0x4a, 0xba, 0x5c, 0x82, 0x1e, 0x20, 0x3d, 0xe1, 0xf9, 0x45, 0x39,
	0x29, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x12, 0x5c, 0xec, 0x45, 0xa9, 0x69, 0xa9,
	0x45, 0xa9, 0x45, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x30, 0xae, 0x92, 0x1e, 0x97, 0x10,
	0xb2, 0xf2, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x90, 0xfa, 0xdc, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4,
	0x54, 0x98, 0x7a, 0x28, 0xd7, 0x28, 0x9b, 0x8b, 0x0b, 0xa1, 0x5e, 0x28, 0x96, 0x8b, 0x37, 0x38,
	0xb1, 0x12, 0x49, 0x40, 0x02, 0xe2, 0x0a, 0x3d, 0x0c, 0x27, 0x48, 0x49, 0x62, 0x91, 0x81, 0xd8,
	0xa6, 0x24, 0xde, 0x74, 0xf9, 0xc9, 0x64, 0x26, 0x41, 0x25, 0x1e, 0x7d, 0xb0, 0x6f, 0xe3, 0xcb,
	0x41, 0xb2, 0x56, 0x8c, 0x5a, 0x49, 0x6c, 0x60, 0x2d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x4c, 0xe7, 0x98, 0xc4, 0x06, 0x01, 0x00, 0x00,
}