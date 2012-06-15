package main

// Sets a cookie with name COOKIE_NAME and value COOKIE_VALUE
// If the cookie is present on the request, reads it back

import (
	"fmt"
	"log"
	"net/http"
)

const (
	COOKIE_NAME = "example-cookie"
	COOKIE_VALUE = "I-am-an-example-cookie"
)

func cookieHandler(w http.ResponseWriter, req *http.Request) {

	// Check for the cookie and get its value
	cookie, err := req.Cookie(COOKIE_NAME)

	// Cookies MUST be set before any output is written to the ResponseWriter
	// MaxAge is in seconds, set to 0 for a session-only cookie
	http.SetCookie(w, &http.Cookie{Name: COOKIE_NAME, Value: COOKIE_VALUE, MaxAge: 30})

	if err != nil {
		fmt.Fprintf(w, "No cookie with the name '%s' was found", COOKIE_NAME)
	} else {
		fmt.Fprintf(w, "Found cookie '%s' with value '%s'", COOKIE_NAME, cookie.Value)
	}
}


func main() {

	serverAddress := ":8080"

	http.HandleFunc("/", cookieHandler)

	log.Println("Starting server on address:", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Panic("Listen and serve failed")
	}
}