package test

import (
  "bytes"
  "beegett/db"
  "beegett/models"
	"encoding/json"
  "log"
  "io"
  "os"
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "beegett/routers"
  "fmt"
	"github.com/astaxie/beego"
  "github.com/tborisova/clean_like_gopher"
	. "github.com/smartystreets/goconvey/convey"
)

var (
  dbCleaner clean_like_gopher.Generic
)

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestMain(m *testing.M) {
  conn, err := db.Connect()
  handleError(err)
  cleaner := getDatabaseCleaner()
  defer conn.Close()
  defer cleaner.Close()
  os.Exit(m.Run())
}

func getDatabaseCleaner() clean_like_gopher.Generic {
  host, database, user, pass := db.GetConnectionOptions()
  options := map[string]string{"host": host, "dbName": database, "username": user, "password": pass}
  cleaner, err := clean_like_gopher.NewCleaningGopher("postgres", options)
  handleError(err)
  dbCleaner = cleaner
  return cleaner
}

func sendRequest(verb string, url string, body io.Reader, verbose bool) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(verb, url, body)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

  if verbose {
	  beego.Trace(fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))
  }

  return w
}

func getBasicDriver() models.Driver {
  driverName := "John Doe"
  driverLicense := "1-02-345-67"
  driver := models.Driver{Name: &driverName, LicenseNumber: &driverLicense}
  return driver
}

func hashToJsonParam(obj interface{}) io.Reader {
  marshalled, _ := json.Marshal(&obj)
  return bytes.NewBuffer(marshalled)
}

// [post] Create
func TestDriverCreateCorrect(t *testing.T) {
  driver := getBasicDriver()
  marshalled, _ := json.Marshal(&driver)

  w := sendRequest("POST", "/v1/driver", bytes.NewBuffer(marshalled), false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Creating a new driver\n", t, func() {
	        Convey("Status Code Should Be 201", func() {
	          So(w.Code, ShouldEqual, 201)
	        })
          Convey("We should get new driver's id", func() {
            So(response.(map[string]interface{})["DriverId"], ShouldNotBeNil)
          })
  })
  dbCleaner.Clean(nil)
}

func TestDriverCreateIncorrectParam(t *testing.T) {
  driverName := "I have no serial number"
  driver := models.Driver{Name: &driverName}

  w := sendRequest("POST", "/v1/driver", hashToJsonParam(driver), false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Creating a driver without needed attribute\n", t, func() {
	        Convey("Status Code Should Be 400", func() {
	          So(w.Code, ShouldEqual, 400)
	        })
          Convey("We should get error description", func() {
            So(response.(map[string]interface{})["description"], ShouldContainSubstring, "violates not-null constraint")
          })
  })
  dbCleaner.Clean(nil)
}

// [get] Retrieve (show)
func TestDriverShowCorrect(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)

  w := sendRequest("GET", fmt.Sprintf("/v1/driver/%d", driver.ID), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting information about an existing driver\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
          Convey("We should get an expected driver name in response struct", func() {
            So(response.(map[string]interface{})["Name"], ShouldEqual, *driver.Name)
          })
	})
  dbCleaner.Clean(nil)
}

func TestDriverShowUnexistant(t *testing.T) {
  w := sendRequest("GET", fmt.Sprintf("/v1/driver/%d", 100500), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting information about unexistant driver\n", t, func() {
	        Convey("Status Code Should Be 400", func() {
	          So(w.Code, ShouldEqual, 400)
	        })
          Convey("Should state that record not found explicitly", func() {
            So(response.(map[string]interface{})["description"], ShouldEqual, "record not found")
          })
	})
  dbCleaner.Clean(nil)
}

func TestDriverShowIncorrectId(t *testing.T) {
  w := sendRequest("GET", fmt.Sprintf("/v1/driver/%s", "something"), nil, false)

	Convey("Subject: Trying to access driver with unappropriate id format\n", t, func() {
	        Convey("Status Code Should Be 404", func() {
	          So(w.Code, ShouldEqual, 404)
	        })
	})
  dbCleaner.Clean(nil)
}

// [get] Retrieve (index)

func generateDriversSequence() {
  drivers := map[string]string{"John Doe": "1-02-345-67", "Jane Doe": "8-76-543-21", "Michael Doe": "9-11-234-25"}
  for name, license := range drivers {
    dr := models.Driver{Name: &name, LicenseNumber: &license}
    db.Conn.Create(&dr)
  }
}

func TestDriverIndexWithOffset(t *testing.T) {
  generateDriversSequence()

  w := sendRequest("GET", fmt.Sprintf("/v1/driver?offset=%d",1), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting drivers with offset\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
          Convey("Should return only two records", func() {
            So(len(response.([]interface{})), ShouldEqual, 2)
          })
	})
  dbCleaner.Clean(nil)
}

func TestDriverIndexWithoutOffset(t *testing.T) {
  generateDriversSequence()

  w := sendRequest("GET", "/v1/driver", nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting drivers without offset\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
          Convey("Should return all three records", func() {
            So(len(response.([]interface{})), ShouldEqual, 3)
          })
	})
  dbCleaner.Clean(nil)
}

func TestDriverIndexWithIncorrectOffset(t *testing.T) {
  generateDriversSequence()

  w := sendRequest("GET", fmt.Sprintf("/v1/driver?offset=%s","wrongValue"), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Getting drivers with incorrect offset value\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
          Convey("Should return all three records", func() {
            So(len(response.([]interface{})), ShouldEqual, 3)
          })
	})
  dbCleaner.Clean(nil)
}

// [put] Update
func TestDriverUpdateCorrect(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)

  w := sendRequest("PUT", fmt.Sprintf("/v1/driver/%d", driver.ID), hashToJsonParam(map[string]string{"Name": "Wendy Sweet"}), false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Updating existing driver\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
	        Convey("Driver should have changed attributes", func() {
            var updatedDriver models.Driver
            db.Conn.First(&updatedDriver, driver.ID)
	          So(*updatedDriver.Name, ShouldEqual, "Wendy Sweet")
	        })
  })
  dbCleaner.Clean(nil)
}

func TestDriverUpdateUnexistantId(t *testing.T) {
  w := sendRequest("PUT", fmt.Sprintf("/v1/driver/%d", 100500), hashToJsonParam(map[string]string{"Name": "Wendy Sweet"}), false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Updating unexisting driver\n", t, func() {
	        Convey("Status Code Should Be 400", func() {
	          So(w.Code, ShouldEqual, 400)
	        })
          Convey("Should include an error description", func() {
            So(response.(map[string]interface{})["description"], ShouldNotBeNil)
          })
  })
  dbCleaner.Clean(nil)
}

// [delete] Delete
func TestDriverDeleteCorrect(t *testing.T) {
  driver := getBasicDriver()
  db.Conn.Create(&driver)

  w := sendRequest("DELETE", fmt.Sprintf("/v1/driver/%d", driver.ID), nil, false)
  var response interface{}
  json.Unmarshal(w.Body.Bytes(), &response)

  Convey("Subject: Removing existing driver\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	          So(w.Code, ShouldEqual, 200)
	        })
	        Convey("Driver should not exist anymore", func() {
            var count int
            db.Conn.Model(&models.Driver{}).Where("id = ?", driver.ID).Count(&count)
	          So(count, ShouldEqual, 0)
	        })
  })
  dbCleaner.Clean(nil)
}
