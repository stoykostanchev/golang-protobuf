package services

import (
	"io/ioutil"
	"sync"
	"time"

	pb "github.com/backend/generated/proto"
	"github.com/golang/protobuf/proto"
)

const cachePath = "images.cache"
const capacity = 1000

/*JokesReplyProvider stores jokes in memory*/
type JokesReplyProvider struct {
	reply            *pb.JokesReply
	newJokesCh       chan *pb.Joke
	cachePath        string
	dummyJokeCounter int32
	jokesAPI         *JokesAPI
}

var singleton *JokesReplyProvider
var once sync.Once

/*GetJokesReplyProvider returns a singleton that caches jokes in-memory and can fetch new jokes from the web*/
func GetJokesReplyProvider() *JokesReplyProvider {
	once.Do(func() {
		api := GetJokesAPI()
		reply := &pb.JokesReply{Jokes: []*pb.Joke{}}
		e := loadPersistantReply(reply, cachePath)
		if e != nil {
			reply = api.GetRandomJokes()
		}
		singleton = &JokesReplyProvider{
			reply:            reply,
			newJokesCh:       make(chan *pb.Joke, capacity),
			cachePath:        cachePath,
			dummyJokeCounter: 0,
			jokesAPI:         api}
	})
	return singleton
}

/*ProvideJokesReply updates the memory cache and returns it*/
func (rp *JokesReplyProvider) ProvideJokesReply() *pb.JokesReply {
forLoop:
	for {
		select {
		case joke := <-rp.newJokesCh:
			rp.reply.Jokes = append(rp.reply.GetJokes(), joke)
		default:
			break forLoop
		}
	}
	return rp.reply
}

/*StartRegularJokeUpdates spawns a new thread and stars re-fetching every SECS seconds*/
func (rp *JokesReplyProvider) StartRegularJokeUpdates(secs time.Duration) {
	go func(c chan *pb.Joke) {
		for range time.Tick(time.Second * secs) {
			joke, err := rp.jokesAPI.GetJoke()
			if err != nil {
				continue
			}
			jr := pb.JokesReply{}
			e := loadPersistantReply(&jr, rp.cachePath)

			if e != nil || !rp.isJokePersistant(joke, &jr) {
				rp.persistJoke(joke, &jr)
			}
			c <- joke
		}
	}(rp.newJokesCh)
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
		return err
	}
	if err := ioutil.WriteFile(rp.cachePath, out, 0644); err != nil {
		return err
	}
	return nil

}

func loadPersistantReply(message proto.Message, cachePath string) error {
	data, err := ioutil.ReadFile(cachePath)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return err
	}
	return nil
}
