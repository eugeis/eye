package main

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"rest/integ"
	"fmt"
	"os"
	"strings"
)

var log = integ.New("[REST] ")

func main() {
	var config *Config
	args := os.Args[1:]
	if len(args) > 0 {
		config = load(strings.Join(args, " "))
	} else {
		config = load("rest.yml")
	}

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
		integ.Debug = false
	}
	config.print()

	serviceFactory := ServiceFactoryByConfig{Config: config}
	defer serviceFactory.Close()

	router := gin.Default()

	router.GET("/:service/ping", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		serviceName := c.Param("service")
		s := serviceFactory.FindService(serviceName)
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
		s := serviceFactory.FindService(serviceName)

		query := c.DefaultQuery("query", "")
		expr := c.DefaultQuery("expr", "")

		if query == "" || expr == "" {
			c.String(http.StatusBadRequest,
				"{ \"ok\": false, \"desc\": \"Please define 'query' and 'expr' parameters\" }")
		} else {
			pattern, err := regexp.Compile(expr)
			if err == nil {
				ok, err := s.Check(query, pattern)
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
	})
	router.Run(fmt.Sprintf(":%d", config.Port))
}
