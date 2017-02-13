package models

import (
	"beegett/db"
	_ "database/sql"
	_ "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Metric struct {
	//NOTE: we use pointers here and in Driver model because gorm forces
	//us to have an object to store its query results in. As soon as we
	//define this object, our strings get "" value, and ints get 0, so
	//we can not refuse to save this data into database (as fields are actually not null...).
	//two options to solve this are pointers and sql.NullString and sql.NullIntegers.

	ID         int `gorm:"not null"`
	MetricName *string
	Value      *int
	Lat        *float64
	Lon        *float64
	Timestamp  *int
	DriverId   *int
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

func UpdateMetric(driverId int, metricId int, updatedMetric Metric) error {
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
	var metric Metric
	var queryResult *gorm.DB

	//NOTE: we double-check existance of metric with particular driverId here,
	//because we dont want to accidentaly delete a metric which belongs to ANOTHER
	//driver

	queryResult = db.Conn.Where(&Metric{ID: metricId, DriverId: &driverId}).First(&metric)
	if queryResult.Error == nil {
		queryResult = db.Conn.Delete(&metric)
	}
	return queryResult.Error
}
