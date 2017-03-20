package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"rest/integ"
	"fmt"
	"os"
	"strings"
	"rest/conf"
	"rest/service"
)

var log = integ.Log

func main() {
	var config *conf.RestConfig
	args := os.Args[1:]
	if len(args) > 0 {
		config = conf.Load(strings.Join(args, " "))
	} else {
		config = conf.Load("rest.yml")
	}

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
		integ.Debug = false
	}
	config.Print()

	serviceFactory := config.ServiceFactory()
	defer serviceFactory.Close()

	var commandCache = integ.NewCache()

	router := gin.Default()

	router.GET("/:service/ping", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		serviceName := c.Param("service")
		s := serviceFactory.Find(serviceName)
		if s != nil {
			err := s.Ping()
			if err == nil {
				c.String(http.StatusOK, "{ \"ok\": true }")
			} else {
				c.String(http.StatusNotAcceptable, fmt.Sprintf("{ \"ok\": false, \"desc:\": %s }", err))
			}
		} else {
			log.Info("/ping: can't found the service '%s'", serviceName)
			c.String(http.StatusBadRequest,
				fmt.Sprintf("{ \"ok\": false, \"desc:\": can't found the service '%s' }", serviceName))
		}
	})

	router.GET("/:service/validate", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		serviceName := c.Param("service")
		s := serviceFactory.Find(serviceName)
		if s != nil {
			query := c.DefaultQuery("query", "")
			expr := c.DefaultQuery("expr", "")

			if query == "" || expr == "" {
				c.String(http.StatusBadRequest,
					"{ \"ok\": false, \"desc\": \"Please define 'query' and 'expr' parameters\" }")
			} else {
				value, err := commandCache.GetOrBuild(buildCommandKey(serviceName, query, expr), func() (interface{}, error) {
					return s.NewСheck(query, expr)
				})
				if err == nil {
					command, _ := value.(service.Check)
					ok, err := command.Check()
					if err != nil {
						c.String(http.StatusBadRequest, fmt.Sprintf("{ \"ok\": false, \"desc:\": %s }", err))
					} else if ok {
						c.String(http.StatusOK, "{ \"ok\": true }")
					} else {
						c.String(http.StatusNotAcceptable, "{ \"ok\": false }")
					}
				} else {
					log.Info("/validate: error parsing '%s' of '%s'", err, expr)
					c.String(http.StatusBadRequest, fmt.Sprintf("{ \"ok\": false, \"desc:\": %s }", err))
				}
			}
		} else {
			log.Info("/validate: can't found the service '%s'", serviceName)
			c.String(http.StatusBadRequest,
				fmt.Sprintf("{ \"ok\": false, \"desc:\": can't found the service '%s' }", serviceName))
		}
	})
	router.Run(fmt.Sprintf(":%d", config.Port))
}

func buildCommandKey(serviceName string, query string, expr string) string {
	return fmt.Sprintf("%s.q(%s).e(%s)", serviceName, query, expr)
}
