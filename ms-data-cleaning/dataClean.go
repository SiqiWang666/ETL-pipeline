package main

import (
	"bufio"
	"strconv"
	"strings"
	"time"
)

//THIS FUNCTIONALITY NEEDS TO BE MOVED TO ms-data-cleaning
//TO TEST YOUR IMPLEMENTATION, COMMENT THIS OUT AND SEE IF LINE COUNT STILL WORKS.
//processLogFile takes in an uploaded logfile, stores the data, processes stats.
func processLogFile(rawLogFile string, fname string) bool {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(rawLogFile)))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	logFile := parseFile(lines, fname)

	//Store parsed logs
	LogStore.StoreLogLine(logFile)

	return true
}

//THIS FUNCTIONALITY NEEDS TO BE MOVED TO ms-data-cleaning
//TO TEST YOUR IMPLEMENTATION, COMMENT THIS OUT AND SEE IF LINE COUNT STILL WORKS.
//parseFile will take a slice of strings and parse the fields.
func parseFile(lines []string, fname string) LogFile {

	//list to store log lines
	lf := LogFile{}

	for _, line := range lines {

		lineSplit := strings.Split(line, " ")
		userAgent := strings.Join(lineSplit[11:], " ")
		status, _ := strconv.Atoi(lineSplit[8])
		totalBytes, _ := strconv.Atoi(lineSplit[9])
		tempLine := LogLine{
			fname,
			line,
			lineSplit[0],
			lineSplit[3] + " " + lineSplit[4],
			lineSplit[5],
			lineSplit[6],
			status,
			totalBytes,
			lineSplit[10],
			userAgent,
			time.Now(),
		}

		lf.Logs = append(lf.Logs, tempLine)
	}
	return lf
}
