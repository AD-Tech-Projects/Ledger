package database

import (
	"os"

	"github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(buildDatabaseDsn()), &gorm.Config{
		Logger: nil,
	})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	logger.Debug("Database connection established successfully")
	return db
}

func buildDatabaseDsn() string {
	return "host=" + getEnv("DB_HOST", "localhost") +
		" user=" + getEnv("DB_USER", "postgres") +
		" password=" + getEnv("DB_PASSWORD", "password") +
		" dbname=" + getEnv("DB_NAME", "ledger") +
		" port=" + getEnv("DB_PORT", "5432") +
		" sslmode=disable"
}

// TODO - move to libs
func getEnv(key string, defaultValue string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue
	}
	return env
}
