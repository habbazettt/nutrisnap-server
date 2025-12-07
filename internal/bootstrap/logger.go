package bootstrap

import (
	"github.com/habbazettt/nutrisnap-server/config"
	"github.com/habbazettt/nutrisnap-server/pkg/logger"
)

func InitLogger(cfg *config.Config) {
	logFormat := "text"
	if cfg.IsProduction() {
		logFormat = "json"
	}

	logger.Init(logger.Config{
		Level:       cfg.Server.LogLevel,
		Format:      logFormat,
		Environment: cfg.Server.Environment,
	})

	logger.Info("logger initialized",
		"environment", cfg.Server.Environment,
		"log_level", cfg.Server.LogLevel,
	)
}
