package controllers

import (
	"beegett/models"
	"encoding/json"
	_ "fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type MetricController struct {
	beego.Controller
}

func (c *MetricController) Post() {
	var mr models.Metric
	driverId, _ := strconv.Atoi(c.Ctx.Input.Param(":driverId"))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &mr)

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = generateJsonError(err)
	} else {
		mr.DriverId = &driverId
		metricId, err := models.AddMetric(mr)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = generateJsonError(err)
		} else {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]int{"MetricId": metricId}
		}
	}
	c.ServeJSON()
}

func (c *MetricController) Get() {
	driverId, _ := strconv.Atoi(c.Ctx.Input.Param(":driverId"))
	metricId, _ := strconv.Atoi(c.Ctx.Input.Param(":metricId"))

	ob, err := models.GetMetric(driverId, metricId)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = generateJsonError(err)
	} else {
		c.Data["json"] = ob
	}

	c.ServeJSON()
}

func (c *MetricController) GetAll() {
	offset, _ := c.GetInt64("offset")
	driverId, _ := strconv.Atoi(c.Ctx.Input.Param(":driverId"))

	mrs, err := models.GetAllMetrics(driverId, 20, int(offset))
	if err == nil {
		c.Data["json"] = mrs
	} else {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = generateJsonError(err)
	}

	c.ServeJSON()
}

func (c *MetricController) Put() {
	var mr models.Metric
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &mr)

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = generateJsonError(err)
	} else {
		metricId, _ := strconv.Atoi(c.Ctx.Input.Param(":metricId"))
		driverId, _ := strconv.Atoi(c.Ctx.Input.Param(":driverId"))
		err = models.UpdateMetric(driverId, metricId, mr)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = generateJsonError(err)
		} else {
			c.Data["json"] = map[string]string{"status": "update success!"}
		}
	}

	c.ServeJSON()
}

func (c *MetricController) Delete() {
	metricId, _ := strconv.Atoi(c.Ctx.Input.Param(":metricId"))
	driverId, _ := strconv.Atoi(c.Ctx.Input.Param(":driverId"))
	err := models.DeleteMetric(driverId, metricId)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = generateJsonError(err)
	} else {
		c.Data["json"] = map[string]string{"status": " delete success!"}
	}
	c.ServeJSON()
}
