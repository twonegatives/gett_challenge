package test

import (
  "bytes"
  "beegett/db"
  "beegett/models"
	"encoding/json"
  _ "log"
  _ "io"
  _ "os"
	_ "net/http"
	_ "net/http/httptest"
	"testing"
	_ "runtime"
	_ "path/filepath"
	_ "beegett/routers"
  "fmt"
	_ "github.com/astaxie/beego"
  _ "github.com/tborisova/clean_like_gopher"
	. "github.com/smartystreets/goconvey/convey"
)

func getBasicMetric() models.Metric {
  metricName := "characteristics.distance_from_home"
  metricValue := 15
  lat := 35.123
  lon := 47.258
  timestamp := 3458339
  driverId := 102030

  metric := models.Metric{MetricName: &metricName, Value: &metricValue, Lat: &lat, Lon: &lon, Timestamp: &timestamp, DriverId: &driverId}
  return metric
}

// [post] Create
func TestMetricCreateCorrect(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)
  metric := getBasicMetric()
  marshalled, _ := json.Marshal(&metric)

  w := sendRequest("POST", fmt.Sprintf("/v1/driver/%d/metric",driver.ID), bytes.NewBuffer(marshalled), false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Creating a new metric\n", t, func() {
	        Convey("Status Code Should Be 201", func() {
	          So(w.Code, ShouldEqual, 201)
	        })
          Convey("We should get new metric's id", func() {
            So(response.(map[string]interface{})["MetricId"], ShouldNotBeNil)
          })
          Convey("Created metric should belong to provided driver", func(){
            metricId := response.(map[string]interface{})["MetricId"]
            var savedMetric models.Metric
            db.Conn.First(&savedMetric, metricId)
            So(*savedMetric.DriverId, ShouldEqual, driver.ID)
          })
  })
  dbCleaner.Clean(nil)
}

// [get] Retrieve (show)
func TestMetricShowCorrect(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)
  metric := getBasicMetric()
  metric.DriverId = &driver.ID
  db.Conn.Create(&metric)

  url := fmt.Sprintf("/v1/driver/%d/metric/%d", driver.ID, metric.ID)
  w := sendRequest("GET", url, nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting information about an existing metric\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
	        Convey("Should return expected metric value in return", func() {
            value := response.(map[string]interface{})["Value"]
	          So(value, ShouldEqual, 15)
	        })
	})
  dbCleaner.Clean(nil)
}

func TestMetricShowUnexistant(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)
  metric := getBasicMetric()
  metric.DriverId = &driver.ID
  db.Conn.Create(&metric)

  url := fmt.Sprintf("/v1/driver/%d/metric/%d", driver.ID + 1, metric.ID)
  w := sendRequest("GET", url, nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Accessing metric with wrong driver id\n", t, func() {
	        Convey("Status Code Should Be 400", func() {
	          So(w.Code, ShouldEqual, 400)
	        })
          Convey("Should state that record not found explicitly", func() {
            So(response.(map[string]interface{})["description"], ShouldEqual, "record not found")
          })
	})
  dbCleaner.Clean(nil)
}

// [get] Retrieve (index)

func generateMetricsSequence(driverId int) {
  for i := 0; i< 3; i++ {
    metric := getBasicMetric()
    metric.DriverId = &driverId
    db.Conn.Create(&metric)
  }
}

func TestMetricIndexWithOffset(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)
  generateMetricsSequence(driver.ID)

  w := sendRequest("GET", fmt.Sprintf("/v1/driver/%d/metric?offset=%d", driver.ID, 1), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting metrics with offset\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
          Convey("Should return only two records", func() {
            So(len(response.([]interface{})), ShouldEqual, 2)
          })
	})
  dbCleaner.Clean(nil)
}

func TestMetricIndexWithoutOffset(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)
  generateMetricsSequence(driver.ID)

  w := sendRequest("GET", fmt.Sprintf("/v1/driver/%d/metric", driver.ID), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting metrics without offset\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
          Convey("Should return all three records", func() {
            So(len(response.([]interface{})), ShouldEqual, 3)
          })
	})
  dbCleaner.Clean(nil)
}

func TestMetricIndexForBlankDriver(t *testing.T) {
  w := sendRequest("GET", fmt.Sprintf("/v1/driver/%d/metric", 100500), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting metrics without offset\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
          Convey("Should not contain any records", func() {
            So(len(response.([]interface{})), ShouldEqual, 0)
          })
	})
  dbCleaner.Clean(nil)
}

// [put] Update
func TestMetricUpdateCorrect(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)
  metric := getBasicMetric()
  metric.DriverId = &driver.ID
  db.Conn.Create(&metric)
  newMetricName := "network.reception_strength"

  w := sendRequest("PUT", fmt.Sprintf("/v1/driver/%d/metric/%d", driver.ID, metric.ID), hashToJsonParam(map[string]string{"MetricName": newMetricName}), false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Updating existing metric\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
	        Convey("Metric should have changed attributes", func() {
            var updatedMetric models.Metric
            db.Conn.First(&updatedMetric, metric.ID)
	          So(*updatedMetric.MetricName, ShouldEqual, newMetricName)
	        })
  })
  dbCleaner.Clean(nil)
}

func TestMetricUpdateUnexistantId(t *testing.T) {
  metric := getBasicMetric()
  db.Conn.Create(&metric)

  w := sendRequest("PUT", fmt.Sprintf("/v1/driver/%d/metric/%d", 111, metric.ID), hashToJsonParam(map[string]string{"MetricName": "random.metric"}), false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Updating unexisting metric\n", t, func() {
	        Convey("Status Code Should Be 400", func() {
	          So(w.Code, ShouldEqual, 400)
	        })
          Convey("Should state that record not found explicitly", func() {
            So(response.(map[string]interface{})["description"], ShouldEqual, "record not found")
          })
  })
  dbCleaner.Clean(nil)
}

// [delete] Delete
func TestMetricDeleteCorrect(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)
  metric := getBasicMetric()
  metric.DriverId = &driver.ID
  db.Conn.Create(&metric)

  w := sendRequest("DELETE", fmt.Sprintf("/v1/driver/%d/metric/%d", driver.ID, metric.ID), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Removing existing metric\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
	        Convey("Metric should not exist anymore", func() {
            var count int
            db.Conn.Model(&models.Metric{}).Where("id = ?", metric.ID).Count(&count)
	          So(count, ShouldEqual, 0)
	        })
  })
  dbCleaner.Clean(nil)
}
