package main

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // 会完成gzip Compressor的注册
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("pb.HelloRequest", in.Name)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("get metadata error")
	}
	for key, val := range md {
		fmt.Printf("%v:%v\n", key, val)
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("say_hello_grpc starting...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
