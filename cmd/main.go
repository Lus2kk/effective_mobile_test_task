package main

import (
	"fmt"
	"log/slog"
	"test_effective_mobile_task/internal/config"

	"github.com/gin-gonic/gin"
)



func main () {


	router := gin.Default()

	//config
	cfg, err := config.MustLoad()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}
	fmt.Printf("Config loaded successfully: %+v\n", cfg)

	//TODO: logger init

	//TODO: init repo

	//TODO: init service

	//TODO: init handler

	//TODO: register routes


	router.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
}