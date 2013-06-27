package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func HandleFormData(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	w.Write([]byte(fmt.Sprintf("Username: %s\nPassword: %s\n", username, password)))
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
	response, postErr := http.PostForm(address, url.Values{"username": {"idiot"}, "password": {"1234"}})
	if postErr != nil {
		panic(postErr.Error())
	}
	defer response.Body.Close()

	body, ioErr := ioutil.ReadAll(response.Body)
	if ioErr != nil {
		panic(ioErr.Error())
	}

	// ioutil will read the response as a []byte
	log.Printf("Response:\n%s\n", string(body))
}