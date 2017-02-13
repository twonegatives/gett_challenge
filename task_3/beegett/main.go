package main

import (
	"beegett/db"
	_ "beegett/models"
	_ "beegett/routers"
	_ "fmt"
	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"log"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := db.Connect()
	handleError(err)
	defer conn.Close()
	//conn.LogMode(true)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()

	//newDriver := models.Driver{Name: "Stephen", LicenseNumber: "10-131-254"}
	//dr, erry := models.AddDriver(newDriver)
	//handleError(erry)
	//fmt.Println(dr)
	//models.UpdateDriver(dr, models.Driver{Name: "Yuriy"})

	//driver, erry2 := models.GetDriver(dr)
	//handleError(erry2)
	//fmt.Println(driver)

	//models.DeleteDriver(dr)

	//drivers, _ := models.GetAllDrivers(20, 0)
	//fmt.Println(drivers)

}

//func getDbConnection() (*gorm.DB) {
//  config    := beego.AppConfig
//  dbHost    := config.String("dbHost")
//  dbName    := config.String("dbName")
//  dbUser    := config.String("dbUser")
//  dbPass    := config.String("dbPass")
//  conn, err := db.Connect(dbHost, dbName, dbUser, dbPass)
//  handleError(err)
//  return conn
//}
