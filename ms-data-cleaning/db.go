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
func (d *Database) StoreValue(key string, value int) error {

	sqlStmt := `
	INSERT INTO SampleTable (key, value)
	VALUES (?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = statement.Exec(key, value)
	if err != nil {
		log.Println("Failed to execute sql", err)
		return err
	}

	return nil

}

//fetchData allows you to fetch data from db.
func (d *Database) fetchValues(fname string) ([]Row, error) {
	rows, err := d.db.Query("SELECT * FROM SampleTable ")
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

//dbinit function will create a table for use for this microservice.
//Change this to include the table you need for your service
func (d *Database) dbInit() {

	//create browser table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS SampleTable (
		key TEXT,
		VALUE int
		)
	`

	statement, _ := d.db.Prepare(sqlStmt)
	statement.Exec()
}

//Store stores a logfile in a database
func (d *Database) StoreLogLine(lf LogFile) {

	sqlStmt := `
	INSERT INTO logs (name, raw_log, 
		remote_addr,
		time_local,
		request_type,
		request_path,
		status,
		body_bytes_sent,
		http_referer,
		http_user_agent,
		created
		) VALUES (?,?, ?,?,?,?,?,?,?,?,?)
	`
	statement, err := d.db.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range lf.Logs {

		statement.Exec(v.Name, v.RawLog, v.RemoteAddr, v.TimeLocal, v.RequestType, v.RequestPath, v.Status, v.BodyBytesSent, v.HTTPReferer, v.HTTPUserAgent, v.Created)
	}
}
