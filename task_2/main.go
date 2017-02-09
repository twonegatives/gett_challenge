package main

import (
  "log"
  _ "fmt"
  _ "database/sql"
  _ "github.com/lib/pq"
)

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func main() {
  db := &Database{}
  db.obtainConnection()
  db.insertRows("drivers", getDrivers())

  var metrics []Metric
  const window int = 500
  hasMoreLines := true
  var pos int64

  for hasMoreLines {
    metrics, pos, hasMoreLines = getMetricsBunch(pos, window)
    if len(metrics) > 0 {
      db.insertRows("metrics", metrics)
    }
  }
}
