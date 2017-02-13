package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Metric struct {
	MetricName string  `json:"metric_name"`
	Value      string  `json:"value"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Timestamp  int     `json:"timestamp"`
	DriverId   string  `json:"driver_id"`
}

func (m Metric) ToSqlParams() string {
	return fmt.Sprintf("(default,'%s',%s,%f,%f,%d,%s)", m.MetricName, m.Value, m.Lat, m.Lon, m.Timestamp, m.DriverId)
}

func (m Metric) HasDriverId() bool {
	_, err := strconv.Atoi(m.DriverId)
	return err == nil
}

func openFileAt(startPos int64) *os.File {
	file, err := os.Open("./data/metrics.json")
	handleError(err)
	_, err = file.Seek(startPos, 0)
	handleError(err)
	return file
}

func getScannerWithRewind(file *os.File, currPos *int64) *bufio.Scanner {
	scanner := bufio.NewScanner(file)

	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		*currPos = *currPos + int64(advance)
		return
	}

	scanner.Split(scanLines)
	return scanner
}

func getMetricsBunch(startPos int64, window int) ([]Metric, int64, bool) {
	// NOTE: seems like the amount of lines inside metrics.json is a hint
	// to stay away from loading all lines in memory at once.
	// so what we're going to do here is to load metrics in batches!
	currPos := startPos
	hasMoreLines := true
	file := openFileAt(startPos)
	scanner := getScannerWithRewind(file, &currPos)
	metrics := make([]Metric, window)
	rollback_j := 0
	defer file.Close()

	for i := 0; i < window; i++ {
		if !scanner.Scan() {
			hasMoreLines = false
			metrics = metrics[0:i]
			break
		}

		var metric Metric
		err := json.Unmarshal([]byte(scanner.Text()), &metric)
		handleError(err)

		if metric.HasDriverId() {
			metrics[i-rollback_j] = metric
		} else {
			rollback_j += 1
		}
	}

	handleError(scanner.Err())

	if rollback_j > 0 {
		metrics = metrics[0 : len(metrics)-rollback_j]
	}

	return metrics, currPos, hasMoreLines
}
