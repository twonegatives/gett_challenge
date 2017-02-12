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
  // NOTE: we'd better reset id sequence after this,
  // as long as we insert ID values by hands
  db.insertRows("drivers", getDrivers())
  db.updateIdSeq("driver_id_seq", "id", "drivers")

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
