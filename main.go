package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/backend/generated/proto"
	grpcweb "github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type jokerServer struct {
	pb.JokerServer
}

func (*jokerServer) GetJokes(ctx context.Context, req *pb.JokesRequest) (*pb.JokesReply, error) {
	joke := pb.Joke{Id: 1, Type: "programming", Setup: "Knock Knock?", Punchline: "Yo mamma"}
	jokes := []*pb.Joke{&joke}
	reply := pb.JokesReply{Jokes: jokes}
	return &reply, nil
}

func main() {
	fmt.Println("Hello")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)

	httpServer := &http.Server{
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Sparta?")
			fmt.Println(grpcweb.ListGRPCResources(grpcServer))
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					fmt.Println("Is grpc!")
					grpcWebServer.ServeHTTP(w, r)
				} else {
					fmt.Println("Is not grpc!")
				}
			}
		}), &http2.Server{}),
	}
	pb.RegisterJokerServer(grpcServer, &jokerServer{})
	// grpcServer.Serve(lis)
	httpServer.Serve(lis)
}
