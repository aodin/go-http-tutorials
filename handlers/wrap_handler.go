package main

import (
	"log"
	"net/http"
)

func WrapHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am written by the wrap handler.\n"))
		f(w, r)
	}
}

func BasicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am written by the basic handler.\n"))
}

func main() {
	http.HandleFunc("/", WrapHandler(BasicHandler))

	address := ":8081"
	log.Println("Starting server on address", address)
	http.ListenAndServe(address, nil)
}
