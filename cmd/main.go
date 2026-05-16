package main

import (
	"fmt"
	"log/slog"
	"test_effective_mobile_task/internal/config"
	"test_effective_mobile_task/internal/handler"
	"test_effective_mobile_task/internal/repo"
	"test_effective_mobile_task/internal/routes"
	"test_effective_mobile_task/internal/service"
	"test_effective_mobile_task/storage"
	"github.com/gin-gonic/gin"
)



func main () {



	//config
	cfg, err := config.MustLoad()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}

	router := gin.Default()

	slog.Info("Config loaded successfully")
	//TODO: logger init

	//TODO: db init
	db, err := storage.InitDatabase(cfg.Database)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return
	}
	defer db.Close()
	slog.Info("Database initialized successfully")

	//TODO: init repo
	repo := repo.NewSubscriptionRepo(db)
	
	slog.Info("Repository initialized successfully")

	//TODO: init service
	service := service.NewSubscriptionService(repo)
	
	slog.Info("Service initialized successfully")

	//TODO: init handler
	handler := handler.NewSubscriptionHandler(service)
	
	slog.Info("Handler initialized successfully")

	//TODO: register routes
	routes.SubscriptionRoutes(router, handler)
	slog.Info("Routes registered successfully! Server is running...")

	router.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
}