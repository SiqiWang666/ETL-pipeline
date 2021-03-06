package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

//Data cleaning related
type BasicService interface {
	sender() bool
}

type Service struct {
	url string
}

//Data cleaning related
func (s *Service) Post(requestBody []byte) (map[string]interface{}, error) {
	resp, err := http.Post(s.url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error making post request to ", s.url, ": ", err)
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Println("Error decoding json:", err)
		return nil, err
	}

	resp.Body.Close()
	return result, nil

}

//Data cleaning related
type DataCleaning struct {
	host        string
	serviceName string
	Service
}

//Data cleaning related
func DataCleaningService() *DataCleaning {
	return &DataCleaning{
		host:        "http://localhost:",
		serviceName: viper.GetString("services.ms-data-cleaning"),
	}
}

//Data cleaning related, call "/data/clean" to clean the data.
//It will trigger the microservice - data- cleaning to do that
func (dc *DataCleaning) CleanData(rawLogFile []byte, fname string) bool {
	requestBody, err := json.Marshal(map[string]string{
		"rawLogFile": string(rawLogFile),
		"fname":      fname,
	})

	if err != nil {
		log.Println("Error parsing CleanData request body:", requestBody)
		return false
	}

	dc.url = dc.genUrl("/data/clean")
	result, err := dc.Post(requestBody)
	if err != nil {
		return false
	}

	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		return false
	}

	log.Println("Data Clean completed")
	return true

}

//Data cleaning related
func (dc *DataCleaning) genUrl(path string) string {
	return dc.host + dc.serviceName + path
}

//handleBrowserCount handles fetching line counts
func handleBrowserCount(w http.ResponseWriter, r *http.Request) {

	// //fetch parameters from url
	// params := mux.Vars(r)

	// fname := params["fname"]
	url := "http://localhost:" + viper.GetString("services.ms-browser-counts") + "/browser/count"

	log.Println("Fetching URL: ", url)
	//make request to ms
	resp, err := http.Get(url)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		e := NewError(http.StatusInternalServerError, result["error"].(string))
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//convert response to proper response
	resOut := ResponseInt{int(result["statusCode"].(float64)), result["message"].(string), ConvertMapInterfaceToMapInt(result["data"])}
	jOut, _ := resOut.JSON()

	//return response to client
	w.WriteHeader(int(result["statusCode"].(float64)))
	fmt.Fprintf(w, jOut)
}

//handleWebsiteCount handles fetching line counts
func handleWebsiteCount(w http.ResponseWriter, r *http.Request) {

	// //fetch parameters from url
	// params := mux.Vars(r)

	// fname := params["fname"]
	url := "http://localhost:" + viper.GetString("services.ms-website-counts") + "/website/count"

	log.Println("Fetching URL: ", url)
	//make request to ms
	resp, err := http.Get(url)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		e := NewError(http.StatusInternalServerError, result["error"].(string))
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//convert response to proper response
	resOut := ResponseInt{int(result["statusCode"].(float64)), result["message"].(string), ConvertMapInterfaceToMapInt(result["data"])}
	jOut, _ := resOut.JSON()

	//return response to client
	w.WriteHeader(int(result["statusCode"].(float64)))
	fmt.Fprintf(w, jOut)
}

// handleVisitorCount handles fetching visitor counts
func handleVisitorCount(w http.ResponseWriter, r *http.Request) {
	//fetch parameters from url
	// params := mux.Vars(r)

	// fname := params["fname"]
	url := "http://localhost:" + viper.GetString("services.ms-visitor-counts") + "/visitor/counts"

	log.Println("Fetching URL: ", url)
	//make request to ms
	resp, err := http.Get(url)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		e := NewError(http.StatusInternalServerError, result["error"].(string))
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//convert response to proper response
	resOut := ResponseInt{int(result["statusCode"].(float64)), result["message"].(string), ConvertMapInterfaceToMapInt(result["data"])}
	jOut, _ := resOut.JSON()

	//return response to client
	w.WriteHeader(int(result["statusCode"].(float64)))
	fmt.Fprintf(w, jOut)
}

//handleLinesCount handles fetching line counts
func handleLinesCount(w http.ResponseWriter, r *http.Request) {

	//fetch parameters from url
	params := mux.Vars(r)

	fname := params["fname"]
	url := "http://localhost:" + viper.GetString("services.ms-line-count") + "/lines/count/" + fname

	log.Println("Fetching URL: ", url)
	//make request to ms
	resp, err := http.Get(url)
	if err != nil {
		e := NewError(http.StatusBadRequest, err.Error())
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	//decode reesponse body
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	//if statusCode is above 300 then its an error, parse and return
	if result["statusCode"].(float64) > 300 {
		log.Println("Error", result["error"].(string))
		e := NewError(http.StatusInternalServerError, result["error"].(string))
		http.Error(w, e.json, http.StatusBadRequest)
		return
	}

	//convert response to proper response
	resOut := ResponseInt{int(result["statusCode"].(float64)), result["message"].(string), ConvertMapInterfaceToMapInt(result["data"])}
	jOut, _ := resOut.JSON()

	//return response to client
	w.WriteHeader(int(result["statusCode"].(float64)))
	fmt.Fprintf(w, jOut)
}

//handleServeUploadPage serves the static html file
func handleServeUploadPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/upload.html")
}

//handleUploadLog handles uploading of log file and triggers etl pipeline
func handleUploadLog(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Received Uploaded File: %+v\n", handler.Filename)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	//clean log file and store in db
	// We call data clean before running the pipeline
	//Currently, processLogFile method is in ms-data-cleaning. To trigger processLogFile, we need to call CleanData which would
	//connect to ms-data-cleaning, then the route of ms-data-cleaning would call handleRoute which contains processLogFile method.
	var resDc bool = DataCleaningService().CleanData(fileBytes, handler.Filename)
	// processLogFile(fileBytes, handler.Filename)
	//

	result := runPipeline(handler.Filename)
	result["Data Cleaner"] = resDc

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>Pipeline Status</h1>")
	fmt.Fprintf(w, "Log File Uploaded Successfully<br>")
	//iterate over results
	for k, v := range result {
		var status string
		//check status
		if v {
			status = `<font size="3" color="green">Completed</font>`
		} else {
			status = `<font size="3" color="red">Failed</font>`
		}
		//print status
		fmt.Fprintf(w, "<strong>"+k+"</strong>: "+status+"<br>")
	}

}
