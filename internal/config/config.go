package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.Logger
	Port   string
	DB     *database.Queries
}

func New() *Config {
	logger := GetLogger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal("failed to load environment variables", zap.Error(err))
	}

	port := os.Getenv("PORT")

	// Database initialisation

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassord := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s", dbUser, dbPassord, dbHost, dbName)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Fatal("failed to open database connection", zap.Error(err))
	}

	return &Config{
		Port:   port,
		Logger: logger,
		DB:     database.New(db),
	}
}
