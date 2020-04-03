package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Deprecated

// VisitorCountsReq defines the body of POST request
type VisitorCountsReq struct {
	Fname string `json: "fname"`
}

// VisitorHandler handles POST request, route /visitor/counts
func VisitorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("VisitorHandler route is processing the request")
	var requestBody VisitorCountsReq
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	//Call db.go func
	// Fetch log line from `logs` table
	lf, err := LogStore.fetchLog()
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	if len(lf.Logs) < 1 {
		e := NewError(http.StatusNotFound, requestBody.Fname+" not found!")
		http.Error(w, e.json, http.StatusNotFound)
		return
	}

	// Main functionality: count visitors
	ips := []string{}
	visitTimes := []time.Time{}
	for _, v := range lf.Logs {
		ips = append(ips, v.RemoteAddr)
		v.TimeLocal = strings.Replace(v.TimeLocal, "[", "", -1)
		v.TimeLocal = strings.Replace(v.TimeLocal, "]", "", -1) // what is -1 meaning?
		visitTimes = append(visitTimes, parseTime(v.TimeLocal))
	}

	uniqueIps := make(map[string][]string)
	for index, ip := range ips {
		key := visitTimes[index].Format("02-01-2006")

		if _, ok := uniqueIps[key]; !ok {
			uniqueIps[key] = []string{}
		}
		if !checkExists(uniqueIps[key], ip) {
			uniqueIps[key] = append(uniqueIps[key], ip)
		}
	}

	visitorCounts := make(map[string]int)
	for k, v := range uniqueIps {
		visitorCounts[k] = len(v)
	}

	for k, v := range visitorCounts {
		if !LogStore.StoreValue(k, v) {
			log.Println("Failed to store ", k)
			return
		}
	}

	//create response and get json
	res := Response{201, "Success"}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jOut)
}

// FetchVisitorHandler handles GET request, route /visitor/counts
func FetchVisitorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchVisitorHandler route is processing the request")

	//Get parameter: fname
	// params := mux.Vars(r)
	// fname := params["fname"]
	//Call db.go func
	vc, err := LogStore.fetchValues()

	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, e.StatusCode)
		return
	}

	if len(vc) < 1 {
		e := NewError(http.StatusNotFound, "Log files are not found!")
		http.Error(w, e.json, http.StatusBadRequest)
	}

	visitorCountsMap := make(map[string]int)
	for _, row := range vc {
		visitorCountsMap[row.Key] = row.Value
	}

	//return application/json response
	res := ResponseInt{201, "Success", visitorCountsMap}
	jOut, _ := res.JSON()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jOut))
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

//check if element found in golang string
func checkExists(list []string, v string) bool {
	for _, a := range list {
		if a == v {
			return true
		}
	}
	return false
}
