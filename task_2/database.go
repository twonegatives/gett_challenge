package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "strings"
  "reflect"
  "database/sql"
  "io/ioutil"
  _ "log"
  _ "github.com/lib/pq"
)

type Config struct {
  Username  string `yaml:"username"`
  Password  string `yaml:"password"`
  Host      string `yaml:"host"`
  Database  string `yaml:"database"`
}

type Database struct {
  config Config
  connection *sql.DB
}

func (d *Database) obtainConnection() {
  username, password, host, dbname := d.obtainCredentials()
  db, err := sql.Open("postgres", "user=" + username + " password=" + password + " host=" + host + " dbname=" + dbname)
  handleError(err)
  d.connection = db
}

func (d *Database) obtainCredentials() (string, string, string, string) {
  if (Config{}) == d.config {
    var config Config
    raw, err := ioutil.ReadFile("database.yml")
    handleError(err)
    err = yaml.Unmarshal(raw, &config)
    handleError(err)
    d.config = config
  }
  return d.config.Username, d.config.Password, d.config.Host, d.config.Database
}

func (d *Database) executeInsert(tableName string, values []string) {
  joined_values := strings.Join(values, ", ")
  query := fmt.Sprintf("INSERT INTO %s VALUES %s;", tableName, joined_values)
  _, err := d.connection.Exec(query)
  handleError(err)
}

func (d *Database) insertRows(tableName string, input interface{}) {
  elements  := reflect.ValueOf(input)
  sqlValues := make([]string, elements.Len())

  for i := 0; i < elements.Len(); i++ {
    element := elements.Index(i)
    string  := element.MethodByName("ToSqlParams").Call(make([]reflect.Value, 0))[0].Interface().(string)
    sqlValues[i] = string
  }

  d.executeInsert(tableName, sqlValues)
}