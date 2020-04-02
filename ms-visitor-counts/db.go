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
func (d *Database) StoreValue(fname string, key string, value int) bool {

	sqlStmt := `
	INSERT INTO Visitors (Fname, key, Counts)
	VALUES (?,?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = statement.Exec(fname, key, value)
	if err != nil {
		log.Println("Failed to execute sql", err)
		return false
	}

	return true
}

//fetchData allows you to fetch visitor counts of a file.
func (d *Database) fetchValues(fname string) ([]Row, error) {
	// fetch values from database
	rows, err := d.db.Query("SELECT key, Counts FROM Visitors WHERE Fname='" + fname + "'")

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

func (d *Database) fetchLog(fname string) (LogFile, error) {
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
			log.Println("Failed to fetch data from table logs: ", err)
			return lf, err
		}
		lf.Logs = append(lf.Logs, logLine)
	}
	return lf, nil
}

//dbinit function will create a table for use for this microservice.
func (d *Database) dbInit() {

	//create browser table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS Visitors (
		Fname text,
		key text,
		Counts int
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()
}
