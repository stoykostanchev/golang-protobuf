package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	pb "github.com/backend/generated/proto"
	"github.com/golang/protobuf/proto"
)

/*JokesReplyProvider stores jokes in memory*/
type JokesReplyProvider struct {
	reply      pb.JokesReply
	newJokesCh chan *pb.Joke
	cachePath  string
}

var singleton *JokesReplyProvider
var once sync.Once

func getManager(cachePath string, capacity int32) *JokesReplyProvider {
	once.Do(func() {
		singleton = &JokesReplyProvider{
			reply:      pb.JokesReply{Jokes: []*pb.Joke{}},
			newJokesCh: make(chan *pb.Joke, 1000),
			cachePath:  cachePath}
	})
	return singleton
}

func (rp *JokesReplyProvider) startRegularJokeUpdates() {
	go func(c chan *pb.Joke) {
		for range time.Tick(time.Second * 1) {
			joke := rp.getJoke()
			jr := pb.JokesReply{}
			e := rp.getJokesReplyCache(&jr)

			if e != nil || rp.jokeNeedsCaching(joke) {
				rp.cacheJoke(joke, &jr)
			}

			c <- joke
		}
	}(rp.newJokesCh)
}

func (rp *JokesReplyProvider) getJoke() *pb.Joke {
	joke := pb.Joke{Id: 1, Type: "programming", Setup: "Knock Knock?", Punchline: "Yo mamma"}
	return &joke
}

func (rp *JokesReplyProvider) jokeNeedsCaching(joke *pb.Joke) bool {
	for _, j := range rp.reply.Jokes {
		if j.GetId() == joke.GetId() {
			return false
		}
	}
	return true
}

func (rp *JokesReplyProvider) getJokesReplyCache(message proto.Message) error {
	data, err := ioutil.ReadFile(rp.cachePath)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message: %w", err)
	}
	return nil
}

func (rp *JokesReplyProvider) cacheJoke(joke *pb.Joke, reply *pb.JokesReply) error {
	reply.Jokes = append(reply.Jokes, joke)
	out, err := proto.Marshal(reply)

	if err != nil {
		log.Fatalln("Failed to encode a reply with jokes:", err)
		return err
	}
	if err := ioutil.WriteFile(rp.cachePath, out, 0644); err != nil {
		log.Fatalln("Failed to write the images cache", err)
		return err
	}
	return nil

}

// forLoop:
// 	for {
// 		select {
// 		case joke := <-ch: // does this fetch only1?
// 			memJokes = append(memJokes, joke)
// 		default:
// 			break forLoop
// 		}
// 	}
