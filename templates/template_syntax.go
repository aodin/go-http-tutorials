package main

import (
	"html/template"
	"log"
	"net/http"
)

// If the template cannot be parsed, Must() will panic
var syntax = template.Must(template.ParseFiles("examples/syntax.html"))

type Item struct {
	Name string
	Value int64
	Status bool
	private int64
}

var dict = map[int64] *Item {
	10: &Item{Name: "Super Bass-O-Matic 1976", Value: 124, private: 2},
	20: &Item{Name: "Swill", Value: 27, private: 7},
	30: &Item{Name: "HiberNol", Value: 8234, private: 3},
}

var list = []*Item {
	&Item{Name: "Super Bass-O-Matic 1976", Value: 124},
	&Item{Name: "Swill", Value: 27},
	&Item{Name: "HiberNol", Value: 8234},
}

var attrs = map[string] interface{} { 	
	"Title": "Syntax Examples",
	"Dict": dict,
	"List": list,
}

// An example http.Handler that serves the homepage template
func HomepageHandler(w http.ResponseWriter, r *http.Request) {

	// The name of the base template must be provided using ExecuteTemplate
	// rather than a vanilla call to Execute(w, attrs).
	// The name "foundation" comes from the first call to define in the
	// foundation.html template: {{ define "foundation" }}
	syntax.Execute(w, attrs)
}

func main() {
	// Provide the http.Handler for the homepage
	http.HandleFunc("/", HomepageHandler)

	// Start the server
	address := ":8081"
	log.Printf("Template server running on address %s\n", address)
	http.ListenAndServe(address, nil)
}