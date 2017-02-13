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

}
