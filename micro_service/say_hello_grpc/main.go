package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	zlog "github.com/vearne/zaplog"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"log"
	"net"

	_ "github.com/apache/skywalking-go"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // 会完成gzip Compressor的注册
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var rdb *redis.Client

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
	val, err := rdb.Incr(context.Background(), "svc-sayHello-grpc").Result()
	zlog.Info("test hello", zap.Int64("val", val), zap.Error(err))
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	zlog.InitLogger("/tmp/sayHelloGrpc.log", "debug")
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

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
