package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

//runPipeline calls functions of other microservices to calculate data metrics
func runPipeline(fname string) map[string]bool {
	//store results of the pipeline
	results := make(map[string]bool)

	//parseFile
	//make call to datacleaner microservice here
	//We won't call data clean in the pipeline, instead, we will trigger data clean
	//before this pipeline

	//count lines
	results["Count Lines"] = countLines(fname)

	//count browsers
	//make call to browserCounts microservice here
	results["Count Browsers"] = countBrowser(fname)

	//count visitors
	//make call to visitorCounts microservice here

	//count websites
	//make call to websiteCounter microservice here
	results["Count Websites"] = countWebsite(fname)

	return results
}

type ReqStruct struct {
	FName   string `json:"fname"`
	ANumber int    `json:"anumber"`
}

func countWebsite(fname string) bool {
	requestBody, err := json.Marshal(map[string]string{
		"fname": fname,
	})
	if err != nil {
		log.Println("Error parsing count website request body:", requestBody)
		return false
	}

	url := "http://localhost:" + viper.GetString("services.ms-website-counts") + "/website/count"

	log.Println("Posting URL: ", url, " with ", string(requestBody))
	//make request to ms
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error making post request to ", url, ": ", err)
		return false
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	// var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding json:", err)
		return false
	}
	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		return false
	}
	//otherwise succesful return true
	return true
}

func countBrowser(fname string) bool {
	requestBody, err := json.Marshal(map[string]string{
		"fname": fname,
	})
	if err != nil {
		log.Println("Error parsing count browser request body:", requestBody)
		return false
	}

	url := "http://localhost:" + viper.GetString("services.ms-browser-counts") + "/browser/count"

	log.Println("Posting URL: ", url, " with ", string(requestBody))
	//make request to ms
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error making post request to ", url, ": ", err)
		return false
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	// var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding json:", err)
		return false
	}
	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		return false
	}
	//otherwise succesful return true
	return true
}

func countLines(fname string) bool {

	requestBody, err := json.Marshal(map[string]string{
		"fname": fname,
	})
	if err != nil {
		log.Println("Error parsing countLines request body:", requestBody)
		return false
	}

	url := "http://localhost:" + viper.GetString("services.ms-line-count") + "/lines/count"

	log.Println("Posting URL: ", url, " with ", string(requestBody))
	//make request to ms
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error making post request to ", url, ": ", err)
		return false
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	// var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding json:", err)
		return false
	}
	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		return false
	}
	//otherwise succesful return true
	return true
}
