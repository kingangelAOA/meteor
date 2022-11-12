package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"
	"web/asset/index"
	"web/asset/static"
	"web/configs"
	"web/db"
	"web/middlewares"
	route "web/routers"
	"web/service"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	configs.FlagInit()
	configs.InitConfig()
	db.InitMongo(configs.Conf.Mongo)
	go service.InitPlugin()
	router = gin.New()
	// router.Use(middlewares.LoggerToFile())
	staticFs := assetfs.AssetFS{
		Asset:    static.Asset,
		AssetDir: static.AssetDir,
	}
	router.StaticFS("/static", &staticFs)
	router.GET("/", bindDataIndexHandler)
	router.GET("/favicon.ico", bindDataFaviconHandler)
	router.Use(middlewares.GlobalRecover)
	router.Use(middlewares.GlobalResponse)
	api := router.Group("/api")
	route.InitRoutes(api)
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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{
		Addr:    ":9090",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	go func() {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}()
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	err := service.CloseAllPlugin()
	log.Printf("Server exiting: %s", err.Error())
}
