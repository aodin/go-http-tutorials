package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSONResponse struct {
	Name string `json:"name"`
	Code int `json:"code"`
}

func HandleJSONResponse(w http.ResponseWriter, r *http.Request) {
	// By default, the response will return success (status code 200)
	response := &JSONResponse{Name: "Success", Code: 200}
	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), 500)
		return
	}
	// Do not write the content before setting the header or the content type
	// will default to "text/plain"
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/", HandleJSONResponse)

	address := ":8081"
	log.Printf("JSON response example running at address %s\n", address)
	http.ListenAndServe(address, nil)
}