package services

import (
	"io/ioutil"
	"log"
	"sync"
	"time"

	pb "github.com/backend/generated/proto"
	"github.com/golang/protobuf/proto"
)

const cachePath = "temp/images.cache"
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
		log.Printf("Instantiating a jokes provider")
		api := GetJokesAPI()
		reply := &pb.JokesReply{Jokes: []*pb.Joke{}}
		persistInitialSet := false
		e := loadPersistantReply(reply, cachePath)
		if e != nil {
			persistInitialSet = true
			log.Printf("Warning : %s", e)
			reply = api.GetRandomJokes()
		}
		log.Printf("Available jokes count: %d", len(reply.Jokes))
		singleton = &JokesReplyProvider{
			reply:            reply,
			newJokesCh:       make(chan *pb.Joke, capacity),
			cachePath:        cachePath,
			dummyJokeCounter: 0,
			jokesAPI:         api}

		if persistInitialSet {
			singleton.persistJokesResponse(reply)
		}
	})
	return singleton
}

/*ProvideJokesReply updates the memory cache and returns it*/
func (rp *JokesReplyProvider) ProvideJokesReply() *pb.JokesReply {
	jc := 0
forLoop:
	for {
		select {
		case joke := <-rp.newJokesCh:
			rp.reply.Jokes = append(rp.reply.GetJokes(), joke)
			jc++
		default:
			log.Printf("New jokes consumed while providing the response: %d", jc)
			break forLoop
		}
	}
	return rp.reply
}

/*StartRegularJokeUpdates spawns a new thread and stars re-fetching every SECS seconds*/
func (rp *JokesReplyProvider) StartRegularJokeUpdates(secs time.Duration) {
	log.Printf("Starting the process of fetching new jokes every %d secs...", secs)

	go func(c chan *pb.Joke) {
		for range time.Tick(time.Second * secs) {
			log.Printf("- Tick - ")
			joke, err := rp.jokesAPI.GetJoke()
			if err != nil {
				log.Printf("Warning - error while fething the joke: %s", err)
				continue
			}
			rp.addJokeToPersisted(joke)
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

func (rp *JokesReplyProvider) persistJokesResponse(resp *pb.JokesReply) error {
	log.Printf("Number of jokes that are getting persisted into memory now: %d", len(resp.Jokes))
	out, err := proto.Marshal(resp)

	if err != nil {
		log.Printf("Error marshalling: %s", err)
		return err
	}
	if err := ioutil.WriteFile(rp.cachePath, out, 0644); err != nil {
		log.Printf("Error writing file: %s", err)
		return err
	}
	return nil
}

func (rp *JokesReplyProvider) addJokeToPersisted(joke *pb.Joke) error {
	log.Printf("Storing joke to memory: %v", joke)
	jr := pb.JokesReply{}
	e := loadPersistantReply(&jr, rp.cachePath)

	if e != nil || !rp.isJokePersistant(joke, &jr) {
		if e != nil {
			log.Printf("Warning - no cache found %s", e)
		} else {
			log.Printf("Joke already saved in the past")
			return nil
		}
	}
	jr.Jokes = append(jr.Jokes, joke)
	return rp.persistJokesResponse(&jr)
}

func loadPersistantReply(message proto.Message, cachePath string) error {
	log.Printf("Loading jokes from the file system")
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
