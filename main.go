package main

import (
	"io/ioutil"
	"meteor/asset/index"
	"meteor/asset/static"
	"meteor/configs"
	"meteor/db"
	"meteor/middlewares"
	route "meteor/routers"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var router *gin.Engine
var config *configs.Config

func init() {
	configs.FlagInit()
	config = getConfig()
	db.InitMongo(config.Mongo)
	router = gin.New()
	staticFs := assetfs.AssetFS{
		Asset:    static.Asset,
		AssetDir: static.AssetDir,
	}
	router.StaticFS("/static", &staticFs)
	router.GET("/", bindDataIndexHandler)
	router.GET("/favicon.ico", bindDataFaviconHandler)
	api := router.Group("/api")
	route.InitRoutes(api)
	router.Use(middlewares.GlobalRecover)
	router.Use(middlewares.GlobalResponse)
}

func bindDataIndexHandler(c *gin.Context) {
	if data, err := index.Asset("index.html"); err != nil {
	} else {
		if _, err := c.Writer.Write(data); err != nil {

		}
	}
}

func bindDataFaviconHandler(c *gin.Context) {
	if data, err := index.Asset("favicon.ico"); err != nil {
	} else {
		if _, err := c.Writer.Write(data); err != nil {

		}
	}
}

func getConfig() *configs.Config {
	yamlFile, err := ioutil.ReadFile(*configs.ConfigPath)
	if err != nil {
		panic(err)
	}
	var c configs.Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	return &c
}

func main() {
	_ = router.Run("0.0.0.0:9090")
}
