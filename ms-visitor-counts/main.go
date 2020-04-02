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

// LogStore is a database connection
var LogStore Database

// ServiceName is the name of the service read from config file
var ServiceName string

// ReadConfig from config file
func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}

func main() {
	ServiceName = "ms-visitor-counts"
	LogStore = Database{}

	ReadConfig()

	var err error

	// For test only
	//LogStore.db, err = sql.Open("sqlite3", "../monolith/ETL.db")

	LogStore.db, err = sql.Open("sqlite3", "../ETL.db")
	if err != nil {
		log.Fatal(err)
	}
	defer LogStore.db.Close()

	// init database
	LogStore.dbInit()

	// Define routes
	r := mux.NewRouter()

	r.HandleFunc("/visitor/counts", VisitorHandler).Methods("POST")
	r.HandleFunc("/visitor/counts/{fname}", FetchVisitorHandler).Methods("GET")

	log.Println("Server is listening on: ", viper.GetString("services."+ServiceName))

	log.Fatal(http.ListenAndServe(":"+viper.GetString("services."+ServiceName), r))
}
