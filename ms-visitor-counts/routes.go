package main

import (
	"log"
	"net/http"
)

// VisitorHandler handles POST request
func VisitorHandler(w http.ResponseWriter, r *http.Request) {
	//TODO call db.go func
	//TODO return json format
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("VisitorHandler works well")
}

// FetchVisitorHandler handles GET request
func FetchVisitorHandler(w http.ResponseWriter, r *http.Request) {
	//TODO call db.go func
	//TODO return json format
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("FetchVisitorHandler works well")
}
