package main

import (
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/backend/generated/proto"
	js "github.com/backend/services"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func main() {
	log.SetOutput(os.Stdout)

	port := "8080"
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal(err)
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
			} else {
				log.Printf("Request not recognised as a grpc-web request")
			}
		}), &http2.Server{}),
	}
	log.Printf("Http server starting on port %s...", port)
	httpServer.Serve(lis)
}
