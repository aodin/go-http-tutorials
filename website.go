package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Website struct {
	Attrs     map[string]interface{}
	StaticDir string
	homepage  *template.Template
	article   []byte
}

// Serve a single file
func (web *Website) IconHandler(w http.ResponseWriter, r *http.Request) {
	// TODO No reason to perform this join on every request
	http.ServeFile(w, r, filepath.Join(web.StaticDir, "favicon.ico"))
}

// Execute the template for each request
func (web *Website) HomepageHandler(w http.ResponseWriter, r *http.Request) {
	// Make a copy of the attrs map to prevent a race condition
	// TODO Is making a copy of attrs necessary?
	attrs := make(map[string]interface{})
	for k, v := range web.Attrs {
		attrs[k] = v
	}
	// Add the attrs for this request to the map
	attrs["IP"] = strings.SplitN(r.RemoteAddr, ":", 2)[0]
	// Execute the template
	web.homepage.ExecuteTemplate(w, "foundation", attrs)
	// Log the request
	Log(r)
}

// Write the pre-executed template as a response
func (web *Website) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(web.article)
	// Log the request
	Log(r)
}

type Resource struct {
	IP    string `json:"ip"`
	Agent string `json:"agent"`
}

func (web *Website) JSONHandler(w http.ResponseWriter, r *http.Request) {
	// Always write the content type before any content!
	w.Header().Set("Content-Type", "application/json")
	resource := &Resource{
		IP:    strings.SplitN(r.RemoteAddr, ":", 2)[0],
		Agent: r.Header.Get("User-Agent"),
	}
	output, err := json.Marshal(resource)
	if err != nil {
		// You may not want to expose your errors to the public
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)
}

func (web *Website) Status404Handler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func Create(staticDir, templateDir string) (*Website, error) {
	// Parse the nested templates
	homepage, err := template.New("homepage").ParseFiles(filepath.Join(templateDir, "homepage.html"), filepath.Join(templateDir, "foundation.html"))
	if err != nil {
		return nil, err
	}
	// Create the attrs that will be passed to templates during execution
	attrs := map[string]interface{}{
		"Title":     "A Website in Go",
		"StaticURL": "/static/",
		"Time":      time.Now(),
	}
	// Execute a template and save its result
	articleTemplate, err := template.ParseFiles(filepath.Join(templateDir, "article.html"))
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	articleTemplate.Execute(buffer, attrs)
	website := &Website{
		Attrs:     attrs,
		StaticDir: staticDir,
		homepage:  homepage,
		article:   buffer.Bytes(),
	}
	// TODO named return types?
	return website, nil
}

func Log(r *http.Request) {
	// By default, a timestamp will be written by the logger
	log.Printf(`"%s %s" %s "%s"`, r.Method, r.URL, strings.SplitN(r.RemoteAddr, ":", 2)[0], r.Header.Get("User-Agent"))
}

func main() {
	// Declare the flags and parse
	port := flag.Int("port", 8000, "Server Port")
	// If no log is specified it will use stdout
	path := flag.String("log", "", "Log File")
	staticDir := flag.String("static", "./_static", "Static Files")
	templateDir := flag.String("templates", ".", "Templates")
	flag.Parse()

	// If a path was set, tell the logger to write to a file rather than stdout
	if len(*path) > 0 {
		flags := os.O_APPEND | os.O_WRONLY | os.O_CREATE
		logf, err := os.OpenFile(*path, flags, 0644)
		if err != nil {
			panic(err)
		}
		defer logf.Close()
		log.SetOutput(logf)
	}

	// Create the website, which loads the templates
	website, err := Create(*staticDir, *templateDir)
	if err != nil {
		panic(err)
	}

	// Add the pages to the multiplexer (this could be done in Create)
	http.HandleFunc("/", website.HomepageHandler)
	http.HandleFunc("/api", website.JSONHandler)
	http.HandleFunc("/article", website.ArticleHandler)
	http.HandleFunc("/404", website.Status404Handler)
	http.HandleFunc("/favicon.ico", website.IconHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*staticDir))))

	// Start the site and server forever
	address := fmt.Sprintf(":%d", *port)
	log.Println("Running on address:", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
