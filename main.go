package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/backend/generated/proto"

	"google.golang.org/grpc"
)

type jokerServer struct {
	pb.UnimplementedJokerServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterJokerServer(grpcServer, &jokerServer{})
	grpcServer.Serve(lis)
}
