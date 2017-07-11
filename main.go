package main

import (
	"encoding/json"
	"eye/core"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/urfave/cli"
	"gee/as"
	"gee/lg"
	"errors"
	"path/filepath"
)

var l = core.Log

func main() {
	app := cli.NewApp()
	app.Name = "eye"
	app.Usage = "diagnosis application for different resource types like database, HTTP, file system."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Configuration files, separated by ','",
			Value: "eye.yml",
		}, cli.StringFlag{
			Name: "environments, e",
			Usage: "Environments (file suffixes) for configurations and properties files." +
				"This files - <file>_<suffix>.<ext> -will be loaded from same directory as a configuration, properties file.",
			Value: "",
		}, cli.StringFlag{
			Name:  "home, b",
			Usage: "Application Home",
			Value: ".",
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
				config, err := loadConfig(c)
				if err != nil {
					return err
				}
				accessFinder, err := as.BuildAccessFinderFromVault("eye",
					c.String("vault_token"), c.String("vault_address"),
					config.Name, config.ExtractAccessKeys())
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
				config, err := loadConfig(c)
				if err != nil {
					return err
				}
				accessFinder, err := as.BuildAccessFinderFromConsole(config.ExtractAccessKeys())
				if err != nil {
					return err
				}
				return start(config, accessFinder)
			},

		}, {
			Name:    "server-file",
			Aliases: []string{"sf"},
			Usage:   "Start Eye server - access data provided from file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "security, s",
					Usage: "Load security configuration from `FILE`",
				},
			},
			Action: func(c *cli.Context) error {
				config, err := loadConfig(c)
				if err != nil {
					return err
				}
				accessFinder, err := as.BuildAccessFinderFromFile(c.String("security"))
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

func loadConfig(c *cli.Context) (ret *core.Config, err error) {
	configFiles := strings.Split(c.GlobalString("c"), ",")
	environments := strings.Split(c.GlobalString("e"), ",")
	if appHome, err := filepath.Abs(c.GlobalString("home")); err == nil {
		ret, err = core.LoadConfig(configFiles, environments, appHome)
	}
	return
}

func prepareDebug(config *core.Config) {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
		lg.Debug = false
	}
	return
}

func start(config *core.Config, accessFinder as.AccessFinder) (err error) {
	prepareDebug(config)
	controller := core.NewEye(config, accessFinder)
	defer controller.Close()
	engine := gin.Default()
	defineRoutes(engine, controller, config)
	return engine.Run(fmt.Sprintf(":%d", config.Port))
}

func defineRoutes(engine *gin.Engine, controller *core.Eye, config *core.Config) {
	currentConfig := config
	engine.Static("/html", "html")
	//engine.StaticFS("/log", LocalFs{})
	engine.Static("/log", fmt.Sprintf("%v/log", currentConfig.AppHome))
	engine.StaticFile("/", "html/doc.html")

	engine.GET("/ping", func(c *gin.Context) {
		response(nil, c)
	})

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

		servicesGroup.GET("/running/compare", func(c *gin.Context) {
			response(controller.CompareRunning(servicesCompare(c)), c)
		})

		servicesGroup.GET("/all/compare", func(c *gin.Context) {
			response(controller.CompareAll(servicesCompare(c)), c)
		})

		servicesGroup.GET("/all/compare/:check", func(c *gin.Context) {
			response(controller.CompareAll(servicesCompare(c)), c)
		})
	}
	checkGroup := engine.Group("/check")
	{
		checkGroup.GET("/:check", func(c *gin.Context) {
			response(controller.Check(c.Param("check")), c)
		})
	}
	exportGroup := engine.Group("/export")
	{
		exportGroup.GET("/:name", func(c *gin.Context) {
			var params = make(map[string]string)
			for k, v := range c.Request.URL.Query() {
				if len(v) > 0 {
					params[k] = v[0]
				}
			}
			response(controller.Export(c.Param("name"), params), c)
		})
	}
	executorGroup := engine.Group("/execute")
	{
		executorGroup.GET("/:name", func(c *gin.Context) {
			var params = make(map[string]string)
			for k, v := range c.Request.URL.Query() {
				if len(v) > 0 {
					params[k] = v[0]
				}
			}
			response(controller.Execute(c.Param("name"), params), c)
		})
	}
	adminGroup := engine.Group("/admin")
	{
		adminGroup.GET("/reload", func(c *gin.Context) {
			config, err := config.Reload()
			if err == nil {
				currentConfig =config
				controller.UpdateConfig(config)
			}
			response(err, c)
		})

		adminGroup.GET("/config", func(c *gin.Context) {
			c.Header("Content-Type", "application/json; charset=UTF-8")
			c.IndentedJSON(http.StatusOK, currentConfig)
		})
	}
}

func services(c *gin.Context) []string {
	return strings.Split(c.DefaultQuery("services", ""), ",")
}

func servicesQuery(c *gin.Context) ([]string, *core.ValidationRequest) {
	return strings.Split(c.DefaultQuery("services", ""), ","), validationReq(c)
}

func servicesCompare(c *gin.Context) ([]string, *core.ValidationRequest) {
	return strings.Split(c.DefaultQuery("services", ""), ","),
		validationReq(c)
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

func serviceQuery(c *gin.Context) (string, *core.ValidationRequest) {
	return c.Param("service"), validationReq(c)
}

func validationReq(c *gin.Context) *core.ValidationRequest {
	return &core.ValidationRequest{
		Query:    c.DefaultQuery("query", ""),
		RegExpr:  c.DefaultQuery("expr", ""),
		EvalExpr: c.Query("eval"), }
}

func response(err error, c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	if err == nil {
		c.String(http.StatusOK, "{ \"ok\": true }")
	} else {
		jsonDesc, _ := json.Marshal(err.Error())
		c.String(http.StatusConflict, fmt.Sprintf("{ \"ok\": false, \"desc:\": %s }", jsonDesc))
	}
}

type LocalFs struct {
}

func (o LocalFs) Open(name string) (http.File, error) {
	return nil, errors.New("Can't find resource" + name)
}
