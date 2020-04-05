package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// LogStore is a database conenction
var LogStore Database

// ServiceName is the name of this service.
var ServiceName string

// ReadConfig reads the config from a file
func ReadConfig() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath("../")
	//Set the config type
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}

func main() {
	ServiceName = "ms-browser-counts"

	LogStore = Database{}

	ReadConfig()

	var err error

	LogStore.db, err = sql.Open("sqlite3", "../ETL.db")

	if err != nil {
		log.Fatal(err)
	}
	defer LogStore.db.Close()
	//init database
	LogStore.dbInit()

	r := mux.NewRouter()
	r.HandleFunc("/browser/count", handleCountBrowsers).Methods("POST")
	r.HandleFunc("/browser/count", handleBrowserCount).Methods("GET")
	log.Println("Listening on: ", viper.GetString("services."+ServiceName))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("services."+ServiceName), r))
}
