package controllers

import (
	"beegett/models"
	"encoding/json"
  "strconv"
  _ "fmt"
	"github.com/astaxie/beego"
)

type DriverController struct {
	beego.Controller
}

func generateJsonError(err error) map[string]interface{} {
  return map[string]interface{}{"description": err.Error()}
}

func (c *DriverController) Post() {
	var dr models.Driver
  err := json.Unmarshal(c.Ctx.Input.RequestBody, &dr)
  if err == nil {
	  driverId, err := models.AddDriver(dr)
    if err != nil {
      c.Ctx.Output.SetStatus(400)
      c.Data["json"] = generateJsonError(err)
    } else {
      c.Ctx.Output.SetStatus(201)
	    c.Data["json"] = map[string]int{"DriverId": driverId}
    }
  } else {
    c.Ctx.Output.SetStatus(400)
		c.Data["json"] = generateJsonError(err)
  }
	c.ServeJSON()
}

func (c *DriverController) Get() {
	driverId, err := strconv.Atoi(c.Ctx.Input.Param(":driverId"))
	if err == nil {
		ob, err := models.GetDriver(driverId)
		if err != nil {
      c.Ctx.Output.SetStatus(400)
      c.Data["json"] = generateJsonError(err)
		} else {
			c.Data["json"] = ob
		}
	} else {
    c.Ctx.Output.SetStatus(500)
    c.Data["json"] = generateJsonError(err)
  }
	c.ServeJSON()
}

func (c *DriverController) GetAll() {
  offset, _ := c.GetInt64("offset")

	drs, err := models.GetAllDrivers(20, int(offset))
  if err == nil {
	  c.Data["json"] = drs
  } else {
    c.Ctx.Output.SetStatus(500)
	  c.Data["json"] = generateJsonError(err)
  }

  c.ServeJSON()
}

func (c *DriverController) Put() {
	var dr models.Driver
  err := json.Unmarshal(c.Ctx.Input.RequestBody, &dr)
  if err == nil {
	  driverId, err := strconv.Atoi(c.Ctx.Input.Param(":driverId"))
    if err == nil {
	    err = models.UpdateDriver(driverId, dr)
	    if err != nil {
        c.Ctx.Output.SetStatus(400)
		    c.Data["json"] = generateJsonError(err)
	    } else {
        c.Data["json"] = map[string]string{"status": "update success!"}
	    }
    } else {
      c.Ctx.Output.SetStatus(500)
      c.Data["json"] = generateJsonError(err)
    }
  } else {
    c.Ctx.Output.SetStatus(400)
		c.Data["json"] = generateJsonError(err)
  }

	c.ServeJSON()
}

func (c *DriverController) Delete() {
	driverId, err := strconv.Atoi(c.Ctx.Input.Param(":driverId"))
  if err == nil {
    err = models.DeleteDriver(driverId)
    if err != nil {
      c.Ctx.Output.SetStatus(400)
      c.Data["json"] = generateJsonError(err)
    } else {
      c.Data["json"] = map[string]string{"status":" delete success!"}
    }
  } else {
    c.Ctx.Output.SetStatus(500)
    c.Data["json"] = generateJsonError(err)
  }
	c.ServeJSON()
}
