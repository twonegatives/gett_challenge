package routers

import (
	"beegett/controllers"

	"github.com/astaxie/beego"
)

func init() {

  beego.Router("/v1/driver/:driverId:int",&controllers.DriverController{}, "get:Get;put:Put;delete:Delete")
  beego.Router("/v1/driver",&controllers.DriverController{}, "get:GetAll;post:Post")

  beego.Router("/v1/driver/:driverId:int/metric/:metricId:int",&controllers.MetricController{}, "get:Get;put:Put;delete:Delete")
  beego.Router("/v1/driver/:driverId:int/metric",&controllers.MetricController{}, "get:GetAll;post:Post")

}
