package main

import (
	"log"
	"net/http"
)

func WrapHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("I am written by the wrap handler.\n"))
		handler(w, req)
	}
}

func BasicHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I am written by the basic handler.\n"))
}

func main() {
	http.HandleFunc("/", WrapHandler(BasicHandler))

	address := ":8081"
	log.Println("Starting server on address", address)
	http.ListenAndServe(address, nil)
}
