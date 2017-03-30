package main

import (
	"encoding/json"
	"eye/core"
	"eye/integ"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/urfave/cli"
)

var l = integ.Log

func main() {
	app := cli.NewApp()
	app.Name = "eye"
	app.Usage = "diagnosis application for different resource types like database, HTTP, file system."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "server-vault",
			Aliases: []string{"sv"},
			Usage:   "Start Eye server - access data provided from Vault",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "vault_token, vt",
					Usage:  "Connect to Vault with the `TOKEN`",
					EnvVar: "VAULT_TOKEN",
				},
				cli.StringFlag{
					Name:   "vault_address, va",
					Usage:  "Connect to Vault server with the `ADDRESS`",
					EnvVar: "VAULT_ADDR",
				},
			},
			Action: func(c *cli.Context) error {
				config, err := core.LoadConfig(c.GlobalString("c"))
				if err != nil {
					return err
				}
				accessFinder, err := core.BuildAccessFinderFromVault(
					c.String("vault_token"), c.String("vault_address"), config)
				if err != nil {
					return err
				}
				return start(config, accessFinder)
			},

		}, {
			Name:    "server-console",
			Aliases: []string{"sc"},
			Usage:   "Start Eye server - access data provided from console input",
			Action: func(c *cli.Context) error {
				config, err := core.LoadConfig(c.GlobalString("c"))
				if err != nil {
					return err
				}
				accessFinder, err := core.BuildAccessFinderFromConsole(config)
				if err != nil {
					return err
				}
				return start(config, accessFinder)
			},

		}, {
			Name:    "server-file",
			Aliases: []string{"sc"},
			Usage:   "Start Eye server - access data provided from file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "security, s",
					Usage: "Load security configuration from `FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				config, err := core.LoadConfig(c.GlobalString("c"))
				if err != nil {
					return err
				}
				accessFinder, err := core.BuildAccessFinderFromFile(c.String("security"))
				if err != nil {
					return err
				}
				return start(config, accessFinder)
			},

		},
	}
	err := app.Run(os.Args)
	if err != nil {
		l.Err("%s", err)
	}
}

func prepareDebug(config *core.Config) {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
		integ.Debug = false
	}
	return
}

func start(config *core.Config, accessFinder core.AccessFinder) (err error) {
	prepareDebug(config)
	controller := core.NewEye(config, accessFinder)
	defer controller.Close()
	engine := gin.Default()
	defineRoutes(engine, controller, config)
	return engine.Run(fmt.Sprintf(":%d", config.Port))
}

func defineRoutes(engine *gin.Engine, controller *core.Eye, config *core.Config) {
	engine.StaticFile("/help", "html/doc.html")
	engine.StaticFile("/", "html/doc.html")

	serviceGroup := engine.Group("/service")
	{
		serviceGroup.GET("/:service/ping", func(c *gin.Context) {
			response(controller.Ping(c.Param("service")), c)
		})

		serviceGroup.GET("/:service/validate", func(c *gin.Context) {
			response(controller.Validate(serviceQuery(c)), c)
		})
	}

	servicesGroup := engine.Group("/services")
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
	checkGroup := engine.Group("/check")
	{
		checkGroup.GET("/:check", func(c *gin.Context) {
			response(controller.Check(c.Param("check")), c)
		})
	}
	adminGroup := engine.Group("/admin")
	{
		adminGroup.GET("/reload", func(c *gin.Context) {
			config, err := config.Reload()
			if err == nil {
				controller.UpdateConfig(config)
			}
			response(err, c)
		})

		adminGroup.GET("/config", func(c *gin.Context) {
			c.Header("Content-Type", "application/json; charset=UTF-8")
			c.IndentedJSON(http.StatusOK, config)
		})
	}
}

func services(c *gin.Context) []string {
	return strings.Split(c.DefaultQuery("services", ""), ",")
}

func servicesQuery(c *gin.Context) ([]string, *core.QueryRequest) {
	return strings.Split(c.DefaultQuery("services", ""), ","), queryReq(c)
}

func servicesCompare(c *gin.Context) ([]string, *core.CompareRequest) {

	return strings.Split(c.DefaultQuery("services", ""), ","),
		&core.CompareRequest{
			QueryRequest: queryReq(c), Tolerance: queryInt("tolerance", c),
			Not:          queryFlag("not", c)}
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
		c.String(http.StatusConflict, fmt.Sprintf("{ \"ok\": false, \"desc:\": %v }", jsonDesc))
	}
}
