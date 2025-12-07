package database

import (
	"fmt"
	"time"

	"github.com/habbazettt/nutrisnap-server/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Connect(cfg Config) (*gorm.DB, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	logger.Info("database connected successfully",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName,
	)

	return db, nil
}

func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	logger.Info("running database auto-migration")

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	logger.Info("database migration completed successfully")
	return nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
