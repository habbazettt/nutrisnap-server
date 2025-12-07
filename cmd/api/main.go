package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/habbazettt/nutrisnap-server/config"
	"github.com/habbazettt/nutrisnap-server/internal/bootstrap"
	"github.com/habbazettt/nutrisnap-server/pkg/database"
	"github.com/habbazettt/nutrisnap-server/pkg/logger"
)

func main() {
	cfg := loadConfig()

	bootstrap.InitLogger(cfg)
	bootstrap.InitDatabase(cfg)

	app := bootstrap.NewApp()

	go startServer(app, cfg.Server.Port)

	waitForShutdown()
	shutdown(app)
}

func loadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		logger.Init(logger.Config{Level: "error", Format: "text", Environment: "development"})
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}
	return cfg
}

func startServer(app interface{ Listen(addr string) error }, port string) {
	logger.Info("server listening", "port", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func shutdown(app interface{ Shutdown() error }) {
	logger.Info("shutting down server...")

	if err := app.Shutdown(); err != nil {
		logger.Error("error during server shutdown", "error", err)
	}

	if err := database.Close(); err != nil {
		logger.Error("error closing database connection", "error", err)
	}

	logger.Info("server stopped")
}
