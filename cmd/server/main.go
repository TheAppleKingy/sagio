package main

import (
	"log/slog"
	"os"
	"sagio/internal/infra/configs"
	// "github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := configs.ParseConfigs(); err != nil {
		slog.Error("Unable to parse configs", "error", err)
		os.Exit(1)
	}
	slog.Info("debug is ", "val", configs.AppConfig.Debug)
	// router := gin.Default()
	// if configs.AppConfig.Debug {
	// 	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }
}
