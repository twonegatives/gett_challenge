package models


import (
  "beegett/db"
  _ "database/sql"
  _ "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Metric struct {
  ID              int      `gorm:"not null"`
  MetricName      *string
  Value           *int
  Lat             *float64
  Lon             *float64
  Timestamp       *int
  DriverId        *int
}

func AddMetric(metric Metric) (int, error) {
  queryResult := db.Conn.Create(&metric)
  return metric.ID, queryResult.Error
}

func GetMetric(driverId int, metricId int) (Metric, error) {
  var metric Metric
  queryResult := db.Conn.Where(&Metric{ID: metricId, DriverId: &driverId}).First(&metric)
  return metric, queryResult.Error
}

func GetAllMetrics(driverId int, limit int, offset int) ([]Metric, error) {
  var metrics []Metric
  queryResult := db.Conn.Limit(limit).Offset(offset).Where(&Metric{DriverId: &driverId}).Find(&metrics)
  return metrics, queryResult.Error
}

func UpdateMetric(driverId int, metricId int, updatedMetric Metric) (error) {
  var metric Metric
  queryResult := db.Conn.Where(&Metric{ID: metricId, DriverId: &driverId}).First(&metric)
  if queryResult.Error != nil {
    return queryResult.Error
  } else {
    queryResult = db.Conn.Model(&metric).Updates(updatedMetric)
    return queryResult.Error
  }
}

func DeleteMetric(driverId int, metricId int) error {
  // NOTE: when delete a record, you need to ensure it's primary field has value,
  // and GORM will use the primary key to delete the record, if primary field's blank,
  // GORM will delete all records for the model
  var metric Metric
  var queryResult *gorm.DB

  queryResult = db.Conn.Where(&Metric{ID: metricId, DriverId: &driverId}).First(&metric)
  if queryResult.Error == nil {
    queryResult = db.Conn.Delete(&metric)
  }
  return queryResult.Error
}
