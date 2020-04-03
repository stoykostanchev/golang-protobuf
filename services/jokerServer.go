package services

import (
	"context"

	pb "github.com/backend/generated/proto"
)

/*JokerServer implements the protobuf scaffolds for a server*/
type JokerServer struct {
	pb.JokerServer
}

/*GetJokes returns whatever is present in the cache in terms of a JokesReply*/
func (*JokerServer) GetJokes(ctx context.Context, req *pb.JokesRequest) (*pb.JokesReply, error) {

	joke := pb.Joke{Id: 1, Type: "programming", Setup: "Knock Knock?", Punchline: "Yo mamma"}
	jokes := []*pb.Joke{&joke}
	reply := pb.JokesReply{Jokes: jokes}
	return &reply, nil
}
