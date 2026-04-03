package main

import (
	"flag"
	"log/slog"
	"os"
	"sagio/internal/infra/configs"

	"github.com/caarlos0/env/v11"
	// "github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

type SetupConfig struct {
	HTTP bool
	AMQP bool
}

func main() {
	if err := env.Parse(configs.AppConfig); err != nil {
		slog.Error("Unable to parse configs", "error", err)
		os.Exit(1)
	}
	var conf SetupConfig
	var resource string
	flag.BoolVar(&conf.HTTP, "http", true, "Enable HTTP transport")
	flag.BoolVar(&conf.AMQP, "amqp", true, "Enable AMQP transport")
	flag.StringVar(&resource, "resource", "init", "Resource name for saga initialization (endpoint/queue)")
	flag.Parse()
	// router := gin.Default()
	// if configs.AppConfig.Debug {
	// 	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }
}
