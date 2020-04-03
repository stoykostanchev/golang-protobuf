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

const cachePath = "images.cache"
const capacity = 1000

func loadPersistantReply(message proto.Message, cachePath string) error {
	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message: %w", err)
	}
	return nil
}

/*JokesReplyProvider stores jokes in memory*/
type JokesReplyProvider struct {
	reply            pb.JokesReply
	newJokesCh       chan *pb.Joke
	cachePath        string
	dummyJokeCounter int32
}

var singleton *JokesReplyProvider
var once sync.Once

/*GetJokesReplyProvider returns a singleton that caches jokes in-memory and can fetch new jokes from the web*/
func GetJokesReplyProvider() *JokesReplyProvider {
	once.Do(func() {
		reply := pb.JokesReply{Jokes: []*pb.Joke{}}
		e := loadPersistantReply(&reply, cachePath)
		if e != nil {
			fmt.Println("Err loading")
		}
		singleton = &JokesReplyProvider{
			reply:            reply,
			newJokesCh:       make(chan *pb.Joke, capacity),
			cachePath:        cachePath,
			dummyJokeCounter: 0}
	})
	return singleton
}

/*ProvideJokesReply updates the memory cache and returns it*/
func (rp *JokesReplyProvider) ProvideJokesReply() *pb.JokesReply {
forLoop:
	for {
		select {
		case joke := <-rp.newJokesCh:
			fmt.Println("Adding all jokes to the memory")
			rp.reply.Jokes = append(rp.reply.GetJokes(), joke)
		default:
			break forLoop
		}
	}
	return &rp.reply
}

/*StartRegularJokeUpdates spawns a new thread and stars re-fetching every SECS seconds*/
func (rp *JokesReplyProvider) StartRegularJokeUpdates(secs time.Duration) {
	fmt.Println("Timer starting")
	go func(c chan *pb.Joke) {
		for range time.Tick(time.Second * secs) {
			joke := rp.getJoke()
			fmt.Println("Tick. Got a new joke.", joke.GetId())
			jr := pb.JokesReply{}
			e := loadPersistantReply(&jr, rp.cachePath)

			if e != nil || !rp.isJokePersistant(joke, &jr) {
				rp.persistJoke(joke, &jr)
				fmt.Println("Persisted")
			} else {
				fmt.Println("Failed to persist or joke did not need persisting")
			}
			c <- joke
		}
	}(rp.newJokesCh)
}

func (rp *JokesReplyProvider) getJoke() *pb.Joke {
	rp.dummyJokeCounter++
	joke := pb.Joke{Id: rp.dummyJokeCounter, Type: "programming", Setup: "Knock Knock?", Punchline: "Yo mamma"}

	return &joke
}

func (rp *JokesReplyProvider) isJokePersistant(joke *pb.Joke, reply *pb.JokesReply) bool {
	for _, j := range reply.Jokes {
		if j.GetId() == joke.GetId() {
			return true
		}
	}
	return false
}

func (rp *JokesReplyProvider) persistJoke(joke *pb.Joke, reply *pb.JokesReply) error {
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
