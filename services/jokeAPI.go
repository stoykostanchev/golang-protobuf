package services

import (
	"io/ioutil"
	"log"
	"net/http"

	pb "github.com/backend/generated/proto"
	jsonpb "github.com/golang/protobuf/jsonpb"
)

/*JokesAPI is a struct responsible for communicating with the external API that provides jokes*/
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
		return nil, err2
	}
	return joke, nil
}

/*GetRandomJokes fetches a list of random jokes*/
func (api *JokesAPI) GetRandomJokes() *pb.JokesReply {
	log.Printf("Fetching new jokes")
	resp, err := http.Get(api.baseURL + "/programming/ten")
	reply := &pb.JokesReply{Jokes: []*pb.Joke{}}

	if err != nil {
		log.Printf("Warning - fetching failed: %s", err)
		return reply
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Warning - failed to read response body: %s", err)
		return reply
	}
	// Response body is an invalid json - an array of jokes
	stringifiedBody := "{ \"jokes\": " + string(body) + "}"

	err = jsonpb.UnmarshalString(stringifiedBody, reply)

	if err != nil {
		log.Printf("Warning - unable to parse the response into a message: %s", err)
		return reply
	}
	return reply
}
