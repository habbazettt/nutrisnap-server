package bootstrap

import (
	"github.com/habbazettt/nutrisnap-server/config"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/pkg/database"
	"github.com/habbazettt/nutrisnap-server/pkg/logger"
	"gorm.io/gorm"
)

func InitDatabase(cfg *config.Config) *gorm.DB {
	db, err := database.Connect(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		panic(err)
	}

	if cfg.IsDevelopment() {
		runMigrations(db)
	}

	return db
}

func runMigrations(db *gorm.DB) {
	if err := database.AutoMigrate(db,
		&models.User{},
		&models.OAuthAccount{},
		&models.Product{},
		&models.Scan{},
		&models.Correction{},
	); err != nil {
		logger.Error("failed to run migrations", "error", err)
		panic(err)
	}
}
