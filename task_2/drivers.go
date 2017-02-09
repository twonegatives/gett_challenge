package main

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
)

type Driver struct {
  Id              int     `json:"id"`
  Name            string  `json:"name"`
  LicenseNumber   string  `json:"license_number"`
}

func getDrivers() []Driver {
  raw, err := ioutil.ReadFile("./data/drivers.json")
  handleError(err)
  var drivers []Driver
  err = json.Unmarshal(raw, &drivers)
  handleError(err)
  return drivers
}

func (d Driver) ToSqlParams() string {
  return fmt.Sprintf("(%d,'%s','%s')", d.Id, d.Name, d.LicenseNumber)
}
