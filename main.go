package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"bytes"
)

const (
	FORWARDURL string = "https://requestb.in/1fnd5sx1"
	PORT       int    = 7070
)

func main() {
	log.Printf("Listening on port %d ... \n\n", PORT)
	http.HandleFunc("/", requestHander)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}

func requestHander(response http.ResponseWriter, request *http.Request) {
	log.Printf("Request Method %s", request.Method)

	forward(request, response)
}

func forward(request *http.Request, response http.ResponseWriter) {

	// Logging request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal("Error in reading request body")
	}
	log.Printf("Request body %s", body)

	// Preparing request
	forwardRequest, err := http.NewRequest(request.Method, FORWARDURL, bytes.NewBufferString(string(body)))

	// Adding all
	forwardRequest.Header = request.Header
	if err != nil {
		panic("Error in preparing forward request")
	}

	defer request.Body.Close()

	client := &http.Client{}

	forwardResponse, err := client.Do(forwardRequest)

	if err != nil {
		panic("Error in sending request using client")
	}

	defer forwardResponse.Body.Close()

	// TODO: check why its receiving encoded response
	body2, err := ioutil.ReadAll(forwardResponse.Body)
	log.Printf("Response from server %s", body2)
	response.Write([]byte(body2))
}
