package main

import (
	"log"
	"net/http"
	"encoding/json"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499{
		log.Println("Responding with 5XX error:", msg)
	}

	// this struct will marshal the json
	type errResponse struct{
		Error string `json:"Error"`
	}
 
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err!= nil{
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}