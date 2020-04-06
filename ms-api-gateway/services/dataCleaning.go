package services

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

type DataCleaning struct {
	host        string
	serviceName string
	Service
}

func DataCleaningService() *DataCleaning {
	return &DataCleaning{
		host:        "http://localhost:",
		serviceName: viper.GetString("services.ms-data-cleaning"),
	}
}

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

func (dc *DataCleaning) genUrl(path string) string {
	return dc.host + dc.serviceName + path
}
