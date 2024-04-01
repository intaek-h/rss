package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type response struct {
		Error string // `json:"message"`
	}

	if code > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	resp := &response{Error: msg}

	respondWithJSON(w, code, resp)
}
