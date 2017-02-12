package models


import (
	"errors"
  "beegett/db"
  _ "database/sql"
  _ "fmt"
  _ "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Driver struct {
  ID              int      `gorm:"not null"`
  Name            *string  `gorm:"size:255;not null"`
  LicenseNumber   *string  `gorm:"size:50;not null"`
}

func AddDriver(driver Driver) (int, error) {
  if db.Conn.NewRecord(driver) {
    queryResult := db.Conn.Create(&driver)
    return driver.ID, queryResult.Error
  } else {
    return 0, errors.New("driver record with this id already exists")
  }
}

func GetDriver(driverId int) (Driver, error) {
  var driver Driver
  queryResult := db.Conn.First(&driver, driverId)
  return driver, queryResult.Error
}

func GetAllDrivers(limit int, offset int) ([]Driver, error) {
  var drivers []Driver
  queryResult := db.Conn.Limit(limit).Offset(offset).Find(&drivers)
  return drivers, queryResult.Error
}

func UpdateDriver(driverId int, updatedDriver Driver) (error) {
  var driver Driver
  queryResult := db.Conn.First(&driver, driverId)
  if queryResult.Error != nil {
    return queryResult.Error
  } else
  {
    queryResult = db.Conn.Model(&driver).Updates(updatedDriver)
    return queryResult.Error
  }
}

func DeleteDriver(driverId int) error {
  // NOTE: when delete a record, you need to ensure it's primary field has value,
  // and GORM will use the primary key to delete the record, if primary field's blank,
  // GORM will delete all records for the model
  queryResult := db.Conn.Delete(Driver{ID: driverId})
  return queryResult.Error
}
