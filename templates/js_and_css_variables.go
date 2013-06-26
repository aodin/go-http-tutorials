package main

import (
	"html/template"
	"log"
	"net/http"
)

// Must() will panic if the template returns a non-nil error
var example = template.Must(template.ParseFiles("examples/js_and_css.html"))

// An example http.Handler that serves a template
func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	// The JS() and CSS() functions will safely escape the given string
	// according to the requested format.
	attrs := map[string] interface{} { 	
		"Title": "I am a Title from Attributes",
		"JS": template.JS("console.log('I am safely injected JavaScript!');"),
		"CSS": template.CSS("p { color: blue; border: 1px solid #aaa }"),
	}
	example.Execute(w, attrs)
}

func main() {
	// Provide the http.Handler for the homepage
	http.HandleFunc("/", TemplateHandler)

	// Start the server
	address := ":8081"
	log.Printf("Template server running on address %s\n", address)
	http.ListenAndServe(address, nil)
}