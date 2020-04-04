package main

import (
	"fmt"
	"net"
	"net/http"

	js "github.com/backend/services"

	pb "github.com/backend/generated/proto"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))

	if err != nil {
		return
	}
	jokesProvider := js.GetJokesReplyProvider()
	jokesProvider.StartRegularJokeUpdates(180)
	grpcServer := grpc.NewServer()
	jokerServer := js.JokerServer{}
	pb.RegisterJokerServer(grpcServer, &jokerServer)

	grpcWebServer := grpcweb.WrapServer(grpcServer)

	httpServer := &http.Server{
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
			w.Header().Set("grpc-status", "")
			w.Header().Set("grpc-message", "")

			if grpcWebServer.IsGrpcWebRequest(r) {
				grpcWebServer.ServeHTTP(w, r)
			}
		}), &http2.Server{}),
	}
	httpServer.Serve(lis)
}
