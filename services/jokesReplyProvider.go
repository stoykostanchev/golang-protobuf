package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	pb "github.com/backend/generated/proto"
	jsonpb "github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const cachePath = "images.cache"
const capacity = 1000

/*JokesAPI is a type that talks to the remote jokes API*/
type JokesAPI struct {
	baseURL string
}

/*GetJokesAPI returns an intance of the jokesAPI, which should be used to talk to the remote server */
func GetJokesAPI() *JokesAPI {
	return &JokesAPI{baseURL: "https://official-joke-api.appspot.com/jokes"}
}

/*GetJoke returns a random joke*/
func (api *JokesAPI) GetJoke() (*pb.Joke, error) {
	resp, err := http.Get(api.baseURL + "/programming/random")
	joke := &pb.Joke{}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	stringifiedBody := string(body)
	// the response is an invalid json - obj wrapped in an arr [{...}]
	r := stringifiedBody[1:(len(stringifiedBody) - 1)]

	err2 := jsonpb.UnmarshalString(r, joke)

	if err2 != nil {
		fmt.Println(err2)
		return nil, err2
	}
	return joke, nil
}

/*GetRandomJokes fetches a list of random jokes*/
func (api *JokesAPI) GetRandomJokes() *pb.JokesReply {
	resp, err := http.Get(api.baseURL + "/programming/ten")
	reply := &pb.JokesReply{Jokes: []*pb.Joke{}}

	if err != nil {
		fmt.Println(err)
		return reply
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return reply
	}
	// Response body is an invalid json - an array of jokes
	stringifiedBody := "{ \"jokes\": " + string(body) + "}"

	err2 := jsonpb.UnmarshalString(stringifiedBody, reply)

	if err2 != nil {
		fmt.Println(err2)
		return reply
	}
	return reply
}

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
			fmt.Println("Adding all jokes to the memory")
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
				fmt.Println("Err when trying to get a joke:", err)
				continue
			}
			fmt.Println("Got a new joke.", joke)
			jr := pb.JokesReply{}
			e := loadPersistantReply(&jr, rp.cachePath)

			if e != nil || !rp.isJokePersistant(joke, &jr) {
				rp.persistJoke(joke, &jr)
			} else {
				fmt.Println("Failed to persist or joke did not need persisting")
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
		log.Fatalln("Failed to encode a reply with jokes:", err)
		return err
	}
	if err := ioutil.WriteFile(rp.cachePath, out, 0644); err != nil {
		log.Fatalln("Failed to write the images cache", err)
		return err
	}
	return nil

}
