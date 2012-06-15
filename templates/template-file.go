package main

// The template file server will use html/template to display the file
// example.html with the request's URI.

import (
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	URI string
}

// ParseFiles will parse the given file, while the Must operator will
// panic if the parse returned any non-nil errors
var templateFile = template.Must(template.ParseFiles("uri.html"))

func templateHandler(w http.ResponseWriter, req *http.Request) {
	templateFile.Execute(w, &Page{req.URL.String()})
}

func main() {

	serverAddress := ":8080"

	http.HandleFunc("/", templateHandler)

	log.Println("Starting server on address:", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Panic("Listen and serve failed")
	}
}