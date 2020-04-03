package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/backend/generated/proto"
	"github.com/golang/protobuf/proto"
	grpcweb "github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

const CACHE_PATH = "images.cache"

type jokerServer struct {
	pb.JokerServer
}

func (*jokerServer) GetJokes(ctx context.Context, req *pb.JokesRequest) (*pb.JokesReply, error) {
	joke := pb.Joke{Id: 1, Type: "programming", Setup: "Knock Knock?", Punchline: "Yo mamma"}
	jokes := []*pb.Joke{&joke}
	reply := pb.JokesReply{Jokes: jokes}
	return &reply, nil
}

func getJoke() *pb.Joke {
	joke := pb.Joke{Id: 1, Type: "programming", Setup: "Knock Knock?", Punchline: "Yo mamma"}
	return &joke
}

func jokeNeedsCaching(joke *pb.Joke, reply *pb.JokesReply) bool {
	for _, j := range reply.Jokes {
		if j.GetId() == joke.GetId() {
			return false
		}
	}
	return true
}

func getJokesReplyCache(message proto.Message) error {
	data, err := ioutil.ReadFile(CACHE_PATH)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message: %w", err)
	}
	return nil
}

func cacheJoke(joke *pb.Joke, reply *pb.JokesReply) error {
	reply.Jokes = append(reply.Jokes, joke)
	out, err := proto.Marshal(reply)

	if err != nil {
		log.Fatalln("Failed to encode a reply with jokes:", err)
		return err
	}
	if err := ioutil.WriteFile(CACHE_PATH, out, 0644); err != nil {
		log.Fatalln("Failed to write the images cache", err)
		return err
	}
	return nil

}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	ch := make(chan *pb.Joke, 1000)

	go func(c chan *pb.Joke) {
		for range time.Tick(time.Second * 1) {
			joke := getJoke()
			var jr pb.JokesReply
			e := getJokesReplyCache(&jr)

			if e != nil || jokeNeedsCaching(joke, &jr) {
				cacheJoke(joke, &jr)
			}

			fmt.Println("Publishing")
			c <- joke
		}
	}(ch)
	memJokes := []*pb.Joke{}
	httpServer := &http.Server{
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
			w.Header().Set("grpc-status", "")
			w.Header().Set("grpc-message", "")
		forLoop:
			for {
				select {
				case joke := <-ch: // does this fetch only1?
					fmt.Println("Consuming!")
					memJokes = append(memJokes, joke)
				default:
					break forLoop
				}
			}
			fmt.Println("----")
			for _, i := range memJokes {
				fmt.Println(i.GetSetup())
			}
			if grpcWebServer.IsGrpcWebRequest(r) {
				grpcWebServer.ServeHTTP(w, r)
			}
		}), &http2.Server{}),
	}
	pb.RegisterJokerServer(grpcServer, &jokerServer{})
	// grpcServer.Serve(lis)
	httpServer.Serve(lis)
}
