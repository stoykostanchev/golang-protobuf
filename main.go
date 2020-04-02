package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)
	http.Handle("/", http.FileServer(http.Dir("/backend")))
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Printf("Listened")
	// see docker-machine ip to figure out what ip to hit
}
