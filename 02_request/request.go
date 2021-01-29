package main

import (
	"fmt"
	"net/http"
	"log"
)

var name string

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL, r.Header.Get("User-Agent"))
	if r.Method == http.MethodPost {
		name = r.PostFormValue("name")
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
