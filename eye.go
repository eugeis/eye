package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"eye/integ"
	"fmt"
	"os"
	"strings"
	"eye/core"
	"strconv"
	"encoding/json"
)

var log = integ.Log

func main() {
	controller, err := core.NewController(configFiles(), []string{"security.yml"})
	if err != nil {
		log.Err("Exit! %v", err)
	}

	if !controller.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
		integ.Debug = false
	}

	defer controller.Close()

	router := gin.Default()

	router.StaticFile("/help", "./html/doc.html")
	router.StaticFile("/", "./html/doc.html")

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

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/reload", func(c *gin.Context) {
			controller.ReloadConfig()
			response(nil, c)
		})

		adminGroup.GET("/config", func(c *gin.Context) {
			c.Header("Content-Type", "application/json; charset=UTF-8")
			c.IndentedJSON(http.StatusOK, controller.Config)
		})
	}

	router.Run(fmt.Sprintf(":%d", controller.Config.Port))
}

func configFiles() (ret []string) {
	args := os.Args[1:]
	if len(args) > 0 {
		ret = args
	} else {
		ret = []string{"eye.yml"}
	}
	return
}

func services(c *gin.Context) ([]string) {
	return strings.Split(c.DefaultQuery("services", ""), ",")
}

func servicesQuery(c *gin.Context) ([]string, *core.QueryRequest) {
	return strings.Split(c.DefaultQuery("services", ""), ","), queryReq(c)
}

func servicesCompare(c *gin.Context) ([]string, *core.CompareRequest) {

	return strings.Split(c.DefaultQuery("services", ""), ","),
		&core.CompareRequest{
			QueryRequest: queryReq(c), Tolerance:
			queryInt("tolerance", c),
			Not: queryFlag("not", c)}
}

func queryInt(key string, c *gin.Context) (ret int) {
	if value := c.Query(key); value != "" {
		ret, _ = strconv.Atoi(value)
	}
	return ret
}

func queryFlag(key string, c *gin.Context) (ret bool) {
	if value, ok := c.GetQuery(key); ok {
		if value == "" {
			ret = true
		} else {
			ret, _ = strconv.ParseBool(value)
		}
	}
	return ret
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
	c.Header("Content-Type", "application/json; charset=UTF-8")
	if err == nil {
		c.String(http.StatusOK, "{ \"ok\": true }")
	} else {
		jsonDesc, _ := json.Marshal(err.Error())
		c.String(http.StatusExpectationFailed, fmt.Sprintf("{ \"ok\": false, \"desc:\": %s }", jsonDesc))
	}
}
