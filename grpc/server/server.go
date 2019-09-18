package server

import (
	"context"
	"crypto/tls"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc/pkg/util"
	"log"
	"net"
	"net/http"
	pb "grpc/proto"
	"path"
	"strings"
	"github.com/elazarl/go-bindata-assetfs"
)

var (
	ServerPort string
	CertServerName string
	CertPemPath string
	CertKeyPath string
	SwaggerDir string
	EndPoint string

	tlsConfig *tls.Config
)

func Run() (err error){
	EndPoint = ":"+ServerPort
	tlsConfig = util.GetTLSConfig(CertPemPath,CertKeyPath)

	conn,err :=net.Listen("tcp",EndPoint)
	if err!=nil{
		log.Println("tcp listen err:%v\n",err)
	}
	srv := newServer()
}

func newServer(conn net.Listener)(*http.Server){
	grpcServer := newGrpc()
	gwmux,err := newGateway()
	if err!=nil{
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/",gwmux)
	mux.Handle("/swagger/",serveSwaggerFile)
	serveSwaggerUI(mux)
	return &http.Server{
		Addr: EndPoint,
		Handler:util.GrpcHandleFunc(grpcServer,mux),
		TLSConfig:tlsConfig,
	}
}

func newGrpc()*grpc.Server{
	creds,err:=credentials.NewServerTLSFromFile(CertPemPath,CertKeyPath)
	if err!=nil{
		panic(err)
	}
	opts:=[]grpc.ServerOption{
		grpc.Creds(creds),
	}
	server:=grpc.NewServer(opts...)
	pb.RegisterHelloWorldServer(server,NewHelloService())
	return server
}

func newGateway()(http.Handler,error){
	ctx :=context.Background()
	dcreds,err := credentials.NewClientTLSFromFile(CertPemPath,CertKeyPath)
	if err!=nil{
		return nil,err
	}
	dopts:=[]grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux:=runtime.NewServeMux()
	if err:=pb.RegisterHelloWorldHandlerFromEndpoint(ctx,gwmux,EndPoint,dopts);err!=nil{
		return nil,err
	}
	return gwmux,nil
}
func serveSwaggerFile(w http.ResponseWriter,r *http.Request){
	if !strings.HasSuffix(r.URL.Path,"swagger.json"){
		log.Println("Not Found:%s",r.URL.Path)
		http.NotFound(w,r)
		return
	}
	p:= strings.TrimPrefix(r.URL.Path,"/swagger/")
	p = path.Join(SwaggerDir,p)
	log.Println("serving swagger-file:%s",p)
	http.ServeFile(w,r,p)
}

func serveSwaggerUI(mux *http.ServeMux){
	//fileServer := http.FileServer(&assetfs.AssetFS{
	//	Asset: swagger.
	//})
}
