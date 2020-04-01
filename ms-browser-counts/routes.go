package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func handleBrowserCount(w http.ResponseWriter, r *http.Request) {
	browserData := LogStore.fetchBrowserData()

	browserStats := make(map[string]int)

	for _, v := range browserData {
		if _, ok := browserStats[v.Browser]; !ok {
			browserStats[v.Browser] = 0
		}
		browserStats[v.Browser] += v.Count

	}

	jOut, err := json.Marshal(browserStats)
	if err != nil {
		log.Println("Error Unmarshalling data", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
}