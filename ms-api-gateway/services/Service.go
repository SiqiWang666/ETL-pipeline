package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type BasicService interface {
	sender() bool
}

type Service struct {
	url string
}

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
