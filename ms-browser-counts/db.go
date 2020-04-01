package main

import (
	"database/sql"
	"log"
)

//Database controls database functionality
type Database struct {
	db *sql.DB
}

func (d *Database) fetchBrowserData() []BrowserCountRow {

	sqlStmt := `
		SELECT * from browsers
	`
	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		log.Println("Failed to fetch browser data: ", err)
	}

	browserStats := []BrowserCountRow{}

	for rows.Next() {
		bcr := BrowserCountRow{}
		rows.Scan(&bcr.Key, &bcr.Date, &bcr.Browser, &bcr.Count)
		browserStats = append(browserStats, bcr)
	}
	return browserStats
}