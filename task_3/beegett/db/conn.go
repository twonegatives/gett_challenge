package db

import (
  "github.com/jinzhu/gorm"
	"github.com/astaxie/beego"
  "fmt"
)

var (
  Conn *gorm.DB
)

func GetConnectionOptions() (string, string, string, string) {
  config := beego.AppConfig
  return config.String("dbHost"), config.String("dbName"), config.String("dbUser"), config.String("dbPass")
}

func getConnectionString() string {
  host, database, user, pass := GetConnectionOptions()
  result := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, database, pass)
  return result
}

func Connect() (*gorm.DB, error) {
  dbConnString := getConnectionString()
  db, err := gorm.Open("postgres", dbConnString)
  Conn = db
  return db, err
}
