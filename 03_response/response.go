package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var html = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Example HTML</title>
  </head>
  <body>
    <p>%s</p>
  </body>
</html>`

func autoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Example text"))
}

func manualHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Example text"))
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html, "Example HTML")
}

func htmlAsTextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, html, "Example HTML")
}

func htmlAsAttachmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="example.html"`)
	fmt.Fprintf(w, html, "Example HTML")
}

func chunkedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "%d\n", i)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Refresh", "2")
	fmt.Fprintf(w, time.Now().String())
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/html", http.StatusFound)
}

func authenticatedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Basic")
	auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Basic ")
	msg := "Invalid username or password"
	raw, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	credentials := strings.SplitN(string(raw), ":", 2)
	if len(credentials) < 2 {
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	if credentials[0] != "admin" || credentials[1] != "password" {
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Success!"))
}

func main() {
	http.HandleFunc("/auto", autoHandler)
	http.HandleFunc("/manual", manualHandler)
	http.HandleFunc("/html", htmlHandler)
	http.HandleFunc("/html-as-text", htmlAsTextHandler)
	http.HandleFunc("/html-as-attachment", htmlAsAttachmentHandler)
	http.HandleFunc("/chunked", chunkedHandler)
	http.HandleFunc("/refresh", refreshHandler)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/redirect", redirectHandler)
	http.HandleFunc("/auth", authenticatedHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
