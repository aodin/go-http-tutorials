package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func HandleFormData(w http.ResponseWriter, r *http.Request) {
	output := fmt.Sprintf(
		"\nUsername: %s\nPassword: %s\n",
		r.FormValue("username"),
		r.FormValue("password"),
	)
	w.Write([]byte(output))
}

func main() {
	// Start the server that will receive the POST data
	// A full address with protocol scheme is required for http.PostForm
	port := 8081
	address := fmt.Sprintf("http://localhost:%d/", port)
	// The server requires us to specify a TCP address to listen on, not
	// a full URL
	serverAddress := fmt.Sprintf(":%d", port)
	go func() {
		http.HandleFunc("/", HandleFormData)
		http.ListenAndServe(serverAddress, nil)
	}()

	// If you receive a dial tcp [address]: connection refused error, try
	// using a short sleep here. The http server can not start before the
	// POST request is made. Don't forget to import "time"!
	// time.Sleep(time.Second)

	// Send POST data and wait for the response
	// url.Values is of type map[string] []string
	values := url.Values{"username": {"idiot"}, "password": {"1234"}}
	response, err := http.PostForm(address, values)
	if err != nil {
		log.Fatalf("Error during PostForm: %s", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error during ReadAll: %s", err)
	}
	log.Printf("%s", body)
}
