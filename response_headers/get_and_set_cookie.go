package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	COOKIE_NAME = "example-cookie"
	COOKIE_VALUE = "I-am-an-example-cookie"
)

// Sets a cookie with name COOKIE_NAME and value COOKIE_VALUE.
// If the cookie is present on the request, reads it back.

func HandleCookie(w http.ResponseWriter, r *http.Request) {
	// Check for the cookie
	cookie, err := r.Cookie(COOKIE_NAME)

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
	http.HandleFunc("/", HandleCookie)

	address := ":8081"
	log.Println("Starting cookie server on address:", address)
	http.ListenAndServe(address, nil)
}