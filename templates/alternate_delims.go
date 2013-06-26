package main

import (
	"html/template"
	"log"
	"net/http"
)

// If the template cannot be parsed, Must() will panic
// Replace the default template variable delimiters `{{` and `}}` with
// an alternate pattern.
var homepageTemplate = template.Must(template.New("homepage").Delims(`<%`, `%>`).ParseFiles("examples/bracket_homepage.html", "examples/bracket_foundation.html"))

// An example http.Handler that serves the homepage template
func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	attrs := map[string] interface{} { 	
		"Title": "I am a Title from Attributes",
	}
	// The name of the base template must be provided using ExecuteTemplate
	// rather than a vanilla call to Execute(w, attrs).
	// The name "foundation" comes from the first call to define in the
	// foundation.html template: {{ define "foundation" }}
	homepageTemplate.ExecuteTemplate(w, "foundation", attrs)
}

func main() {
	// Provide the http.Handler for the homepage
	http.HandleFunc("/", HomepageHandler)

	// Start the server
	address := ":8081"
	log.Printf("Template server running on address %s\n", address)
	http.ListenAndServe(address, nil)
}