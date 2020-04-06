package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"log"
	"net/http"
)

type WebsiteCountRequest struct {
	FName string `json:"fname"`
}

func parseWebsite(url string) string {
	temp := url[7:]
	result := strings.Split(temp, "/")
	return result[0]
}

func handleCountWebsites(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Website Counts")
	//object for body
	var WCR WebsiteCountRequest

	//decode data in body
	err := json.NewDecoder(r.Body).Decode(&WCR)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//fetch data for file name
	lf, err := LogStore.fetchData(WCR.FName)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	//if no lines then file not found, throw error
	if len(lf.Logs) < 1 {
		e := NewError(http.StatusNotFound, WCR.FName+" not found!")
		http.Error(w, e.json, http.StatusNotFound)
		return
	}

	//place to store website url
	websiteCounts := make(map[string]int)

	//iterate and store website url uniquely.
	for _, v := range lf.Logs {
		// log.Printf(v.HTTPReferer)
		website := parseWebsite(v.HTTPReferer)
		if _, ok := websiteCounts[website]; !ok {
			websiteCounts[website] = 1
		}
		websiteCounts[website]++
	}

	for b, v := range websiteCounts {
		key := WCR.FName + "_" + b
		if LogStore.storeWebsiteCount(key, b, v) != nil {
			log.Println("Failed to store ", b)
		}
	}

	//create response and get json
	res := Response{201, "Success"}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jOut)
}

func handleWebsiteCount(w http.ResponseWriter, r *http.Request) {

	websiteData, e := LogStore.fetchWebsiteData()
	if e != nil {
		log.Println("Error fetching data...")
		return
	}

	websiteStats := make(map[string]int)

	for _, v := range websiteData {
		if _, ok := websiteStats[v.Website]; !ok {
			websiteStats[v.Website] = 0
		}
		websiteStats[v.Website] += v.Count
	}

	res := ResponseInt{201, "Success", websiteStats}
	jOut, err := res.JSON()
	if err != nil {
		log.Println("Error parsing data", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
}
