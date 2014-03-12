package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Write some bytes
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(`Content-Type`, `text/html; charset=utf-8`)
	w.Write([]byte(`<!doctype html><html><body><a href="/a">Click me</a></body></html>`))
}

func FavIconHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(``))
}

// Log writes request details to the standard logger
// * Remote address, first checking if there is a proxied IP header
// * Method and URL
// * User agent
// * Referer (i.e. the referring address)
// Attributes are quote escaped when necessary
func LogRequest(r *http.Request) {
	ip := strings.SplitN(r.Header.Get("X-Real-IP"), ":", 2)[0]
	if ip == "" {
		ip = strings.SplitN(r.RemoteAddr, ":", 2)[0]
	}
	log.Printf(
		`%q %s %q %q`,
		fmt.Sprintf(`%s %s`, r.Method, r.URL),
		ip,
		r.Header.Get("User-Agent"),
		r.Header.Get("Referer"),
	)
}

// We can also turn the Log function into a "wrapper"
func LogWrapper(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)

		// Log the request after handling
		// This could be performed in a goroutine if it contained any
		// long-running operations
		LogRequest(r)
	}
}

func main() {
	// Declare the flags and parse
	port := flag.Int("port", 8000, "Server Port")
	// If no log is specified it will use stdout
	path := flag.String("log", "", "Log File")
	flag.Parse()

	// If a path was set, tell the logger to write to a file rather than stdout
	if len(*path) > 0 {
		flags := os.O_APPEND | os.O_WRONLY | os.O_CREATE
		logf, logErr := os.OpenFile(*path, flags, 0644)
		if logErr != nil {
			panic(logErr)
		}
		defer logf.Close()
		log.SetOutput(logf)
	}

	http.HandleFunc("/", LogWrapper(HomeHandler))

	// Catch the favicon request
	http.HandleFunc("/favicon.ico", FavIconHandler)

	// Start the site and serve forever
	address := fmt.Sprintf(":%d", *port)
	log.Println("Running on address:", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
