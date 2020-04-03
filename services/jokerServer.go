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
	jokesProvider := GetJokesReplyProvider()
	reply := jokesProvider.ProvideJokesReply()
	return reply, nil
}
