package main

import (
	"log"
	"net/http"
)

func HandleSuccess(w http.ResponseWriter, r *http.Request) {
	// By default, the response will return success (status code 200)
	w.Write([]byte("Success"))
}

func HandleNoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 301)
}

func HandleBadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", 400)
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	// If any content is written to the response before the call to NotFound
	// the request will return success (status code 200).
	// Do no write any content before calling NotFound
	http.NotFound(w, r)
}

func HandleServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Error", 500)
}



func main() {
	http.HandleFunc("/", HandleSuccess)
	http.HandleFunc("/204", HandleNoContent)
	http.HandleFunc("/301", HandleRedirect)
	http.HandleFunc("/400", HandleBadRequest)
	http.HandleFunc("/404", HandleNotFound)
	http.HandleFunc("/500", HandleServerError)

	address := ":8081"
	log.Printf("Status code example running at address %s\n", address)
	http.ListenAndServe(address, nil)
}