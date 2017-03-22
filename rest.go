package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"rest/integ"
	"fmt"
	"os"
	"strings"
	"rest/core"
)

var log = integ.Log

func main() {
	var config *core.Config
	args := os.Args[1:]
	if len(args) > 0 {
		config = core.Load(strings.Join(args, " "))
	} else {
		config = core.Load("rest.yml")
	}

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
		integ.Debug = false
	}
	config.Print()

	controller := core.NewController(config)

	defer controller.Close()

	router := gin.Default()

	serviceGroup := router.Group("/service")
	{
		serviceGroup.GET("/:service/ping", func(c *gin.Context) {
			response(controller.Ping(c.Param("service")), c)
		})

		serviceGroup.GET("/:service/validate", func(c *gin.Context) {
			response(controller.Validate(serviceQuery(c)), c)
		})
	}

	servicesGroup := router.Group("/services")
	{
		servicesGroup.GET("/any/ping", func(c *gin.Context) {
			response(controller.PingAny(services(c)), c)
		})

		servicesGroup.GET("/all/ping", func(c *gin.Context) {
			response(controller.PingAll(services(c)), c)
		})

		servicesGroup.GET("/any/validate", func(c *gin.Context) {
			response(controller.ValidateAny(servicesQuery(c)), c)
		})

		servicesGroup.GET("/running/validate", func(c *gin.Context) {
			response(controller.ValidateRunning(servicesQuery(c)), c)
		})

		servicesGroup.GET("/all/validate", func(c *gin.Context) {
			response(controller.ValidateAll(servicesQuery(c)), c)
		})

		servicesGroup.GET("/any/compare", func(c *gin.Context) {
			response(controller.CompareAny(servicesCompare(c)), c)
		})

		servicesGroup.GET("/running/compare", func(c *gin.Context) {
			response(controller.CompareRunning(servicesCompare(c)), c)
		})

		servicesGroup.GET("/all/compare", func(c *gin.Context) {
			response(controller.CompareAll(servicesCompare(c)), c)
		})
	}

	checkGroup := router.Group("/check")
	{
		checkGroup.GET("/:check", func(c *gin.Context) {
			response(controller.Check(c.Param("check")), c)
		})
	}

	router.Run(fmt.Sprintf(":%d", config.Port))
}

func services(c *gin.Context) ([]string) {
	return strings.Split(c.DefaultQuery("services", ""), ",")
}

func servicesQuery(c *gin.Context) ([]string, *core.QueryRequest) {
	return strings.Split(c.DefaultQuery("services", ""), ","), queryReq(c)
}

func servicesCompare(c *gin.Context) ([]string, *core.CompareRequest) {
	return strings.Split(c.DefaultQuery("services", ""), ","),
		&core.CompareRequest{QueryRequest: queryReq(c), Operator: c.DefaultQuery("op", "equals")}
}

func serviceQuery(c *gin.Context) (string, *core.QueryRequest) {
	return c.Param("service"), queryReq(c)
}

func queryReq(c *gin.Context) *core.QueryRequest {
	return &core.QueryRequest{
		Query: c.DefaultQuery("query", ""),
		Expr:  c.DefaultQuery("expr", "")}
}

func response(err error, c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if err == nil {
		c.String(http.StatusOK, "{ \"ok\": true }")
	} else {
		c.String(http.StatusNotAcceptable, fmt.Sprintf("{ \"ok\": false, \"desc:\": %s }", err))
	}
}
