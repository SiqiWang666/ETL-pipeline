package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"log"
	"net/http"
)

type BrowserCountRequest struct {
	FName string `json:"fname"`
}

//parseBrowser parses user agent
func parseBrowser(ua string) string {
	browsers := []string{"Firefox", "Chrome", "Opera", "Safari", "MSIE"}

	for _, b := range browsers {
		if strings.Contains(ua, b) {
			return b
		}
	}
	return "Other"
}

//format time appropriately
func parseTime(t string) time.Time {
	//January 2nd, 3:04:05 PM of 2006, UTC-0700.
	tm, err := time.Parse("02/Jan/2006:15:04:05 -0700", t)
	if err != nil {
		log.Fatal("Failed to parse time", err)
	}
	return tm
}

func handleCountBrowsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Browser Counts")
	//object for body
	var BCR BrowserCountRequest

	//decode data in body
	err := json.NewDecoder(r.Body).Decode(&BCR)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//fetch data for file name
	lf, err := LogStore.fetchData(BCR.FName)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	//if no lines then file not found, throw error
	if len(lf.Logs) < 1 {
		e := NewError(http.StatusNotFound, BCR.FName+" not found!")
		http.Error(w, e.json, http.StatusNotFound)
		return
	}
	// fetch browser data
	visitTimes := []time.Time{}
	browsers := []string{}

	lf, e := LogStore.fetchData(BCR.FName)
	if e != nil {
		log.Println("Error bad request...")
		return
	}
	for _, v := range lf.Logs {
		browsers = append(browsers, parseBrowser(v.HTTPUserAgent))
		v.TimeLocal = strings.Replace(v.TimeLocal, "[", "", -1)
		v.TimeLocal = strings.Replace(v.TimeLocal, "]", "", -1)
		visitTimes = append(visitTimes, parseTime(v.TimeLocal))
	}
	//place to store IP's
	browserCounts := make(map[string]int)

	//iterate and store IP addresses uniquely.
	for index, b := range browsers {
		key := b + "_" + visitTimes[index].Format("02-01-2006")
		if _, ok := browserCounts[key]; !ok {
			browserCounts[key] = 0
		}
		browserCounts[key]++
	}

	for b, v := range browserCounts {
		keyPieces := strings.Split(b, "_")
		if LogStore.storeBrowserCount(b, keyPieces[1], keyPieces[0], v) != nil {
			log.Println("Failed to store ", b)
		}
	}

	//create response and get json
	res := Response{201, "Success"}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jOut)
}

//handleLineCount
func handleBrowserCount(w http.ResponseWriter, r *http.Request) {

	//fetch parameters from url
	// params := mux.Vars(r)

	// fname := params["fname"]

	browserData, e := LogStore.fetchBrowserData()
	if e != nil {
		log.Println("Error fetching data...")
		return
	}

	browserStats := make(map[string]int)

	for _, v := range browserData {
		if _, ok := browserStats[v.Browser]; !ok {
			browserStats[v.Browser] = 0
		}
		browserStats[v.Browser] += v.Count

	}

	res := ResponseInt{201, "Success", browserStats}
	jOut, err := res.JSON()
	if err != nil {
		log.Println("Error parsing data", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
}
