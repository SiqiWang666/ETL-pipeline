package main

import (
	"database/sql"
	"log"
)

//Database controls database functionality
type Database struct {
	db *sql.DB
}

//Row represents a row in the database for line counts
type Row struct {
	Key   string
	Value int
}

//StoreValue stores a key and value in a database
func (d *Database) StoreValue(key string, value int) bool {

	sqlStmt := `
	INSERT INTO Visitors (Date, Counts)
	VALUES (?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = statement.Exec(key, value)
	if err != nil {
		log.Println("Failed to execute sql", err)
		return false
	}

	return true
}

//fetchData allows you to fetch visitor counts of all log files.
func (d *Database) fetchValues() ([]Row, error) {
	// fetch values from database
	rows, err := d.db.Query("SELECT * FROM Visitors")

	if err != nil {
		log.Println(err)
		return nil, err
	}
	rs := []Row{}
	for rows.Next() {
		r := Row{}
		err := rows.Scan(&r.Key,
			&r.Value)
		if err != nil {
			log.Println("Failed to fetch data from db: ", err)
			return rs, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}

//fetchLog allows you to fetch list of log lines of a file
func (d *Database) fetchLog() (LogFile, error) {
	lf := LogFile{}
	rows, err := d.db.Query("SELECT * FROM logs")
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
			log.Println("Failed to fetch data from table logs: ", err)
			return lf, err
		}
		lf.Logs = append(lf.Logs, logLine)
	}
	return lf, nil
}

//dbinit function will create a table for use for this microservice.
func (d *Database) dbInit() {
	//create visitors table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS Visitors (
		Date text,
		Counts int
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()
}
