package main

import (
	"database/sql"
	"log"
)

//BrowserCountRow represents a row in the database for browser counts
type WebsiteCountRow struct {
	Key     string
	Website string
	Count   int
}

//Database controls database functionality
type Database struct {
	db *sql.DB
}

//storeWebsiteCount stores the count for browsers in logs by date
func (d *Database) storeWebsiteCount(key string, w string, c int) error {
	// check if the website had already been stored
	sqlStmt := `
		UPDATE websites
		SET count = count + ?
		WHERE website = ?
		`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("Error cannot prepare update code", err)
		return err
	}
	res, err := statement.Exec(c, w)
	if err != nil {
		log.Println("Error cannot run update code", err)
		return err
	}
	// if nothing
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println("Error checking effect", err)
		return err
	}
	if affect != 0 {
		return nil
	}

	sqlStmt = `
	INSERT INTO websites (
		key,
		website,
		count
		) VALUES (?,?,?)
	`
	statement, err = d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println("Error cannot prepare insert code", err)
		return err
	}

	_, err = statement.Exec(key, w, c)

	if err != nil {
		log.Println("Error cannot run insert code", err)
		return err
	}
	return nil
}

func (d *Database) fetchWebsiteData() ([]WebsiteCountRow, error) {

	sqlStmt := "SELECT * from websites"
	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		log.Println("Failed to fetch website data: ", err)
		return nil, err
	}

	websiteStats := []WebsiteCountRow{}

	for rows.Next() {
		wcr := WebsiteCountRow{}
		rows.Scan(&wcr.Key, &wcr.Website, &wcr.Count)
		websiteStats = append(websiteStats, wcr)
	}
	return websiteStats, nil
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
	CREATE TABLE IF NOT EXISTS websites (
		key TEXT,
		website TEXT,
		count int
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()
}
