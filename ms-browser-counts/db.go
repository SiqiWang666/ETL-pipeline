package main

import (
	"database/sql"
	"log"
)

//BrowserCountRow represents a row in the database for browser counts
type BrowserCountRow struct {
	Key     string
	Date    string
	Browser string
	Count   int
}

//Database controls database functionality
type Database struct {
	db *sql.DB
}

//storeBrowserCount stores the count for browsers in logs by date
func (d *Database) storeBrowserCount(key string, dt string, b string, c int) error {
	sqlStmt := `
	INSERT INTO browsers (
		key,
		date,
		browser,
		count
		) VALUES (?,?,?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = statement.Exec(key, dt, b, c)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (d *Database) fetchBrowserData(fname string) ([]BrowserCountRow, error) {

	sqlStmt := "SELECT * from browsers where key ='" + fname + "'"
	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		log.Println("Failed to fetch browser data: ", err)
		return nil, err
	}

	browserStats := []BrowserCountRow{}

	for rows.Next() {
		bcr := BrowserCountRow{}
		rows.Scan(&bcr.Key, &bcr.Date, &bcr.Browser, &bcr.Count)
		browserStats = append(browserStats, bcr)
	}
	return browserStats, nil
}

//fetchData allows you to fetch log data from db.
func (d *Database) fetchData(fname string) (LogFile, error) {
	lf := LogFile{}
	rows, err := d.db.Query("SELECT * FROM logs where name='" + fname + "'")
	if err != nil {
		log.Println(err)
		return lf, err
	}
	for rows.Next() {
		logLine := LogLine{}
		err := rows.Scan(&logLine.Name,
			&logLine.RawLog,
			&logLine.RemoteAddr,
			&logLine.TimeLocal,
			&logLine.RequestType,
			&logLine.RequestPath,
			&logLine.Status,
			&logLine.BodyBytesSent,
			&logLine.HTTPReferer,
			&logLine.HTTPUserAgent,
			&logLine.Created)
		if err != nil {
			log.Println("Failed to fetch data from db: ", err)
			return lf, err
		}
		lf.Logs = append(lf.Logs, logLine)
	}
	return lf, nil
}

func (d *Database) dbInit() {

	//create browser table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS browsers (
		key TEXT,
		date TEXT,
		browser TEXT,
		count int
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()
}
