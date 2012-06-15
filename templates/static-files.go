package main

// Creates a http server that will serve static files accessed by an URI
// beginning with the string specified by STATIC_URI.
// The static files reside in the directory specified by STATIC_DIR.

// When the server is running, an example static file can be accessed through
// your browser at http://localhost:8080/static/example.js

import (
	"log"
	"net/http"
)

const (
	STATIC_URI = "/static/"
	STATIC_DIR = "./static-examples/"
)

func main() {

	serverAddress := ":8080"

	// Returns an instance of http.Handler that will serve the contents
	// of the filesystem starting at and below the given directory.
	fileHandler := http.FileServer(http.Dir(STATIC_DIR))

	// Strip prefix removes the given string from the URI, and passes it to
	// the given handler.
	prefixlessHandler := http.StripPrefix(STATIC_URI, fileHandler)

	// Handle will pass any URI beginning with the supplied string to the
	// given handler.
	http.Handle(STATIC_URI, prefixlessHandler)

	log.Println("Starting server on address:", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Panic("Listen and serve failed")
	}
}