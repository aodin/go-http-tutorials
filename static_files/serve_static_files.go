package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// The static directory is the name of the local filesystem directory
	// that stores static files you wish to serve.
	// The trailing slash is optional as long as a directory is provided.
	staticDirectory := "examples/"

	// Instead of specifying an absolute path name, you can use the 'os'
	// package to determine the current working directory and concatenate
	// the path to the static files.
	currentDirectory, osErr := os.Getwd()
	if osErr != nil {
		panic("Could not get the working directory:" + osErr.Error())
	}
	absStaticDirectory := filepath.Join(currentDirectory, staticDirectory)

	// The static URL is the root URL that will be used to serve the static
	// files, e.g. localhost:8081/static/. It is specified as an absolute URL
	// and the trailing slash is required (if missing, you can still see a
	// list of resources at /static, but will not be able to browse to a
	// specific resource).
	// Do not forget the trailing slash!
	staticURL := "/static/"

	// The address of the server to be started. This will only serve 
	// localhost. Do not forget the colon before the port number!
	address := ":8081"

	http.Handle(staticURL, http.StripPrefix(staticURL, http.FileServer(http.Dir(absStaticDirectory))))
	log.Printf("Static file server running at address %s\n", address)
	http.ListenAndServe(address, nil)
	// You can now browse the resources at <addr>/static/
}