package main

import (
	"log"
	"os"

	db "github.com/ln0rd/tech_challenge_12soat/internal/infrastructure/db"
	routes "github.com/ln0rd/tech_challenge_12soat/internal/interface/http"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	err    error
	logger *zap.Logger
)

func main() {
	config := zap.NewProductionConfig()
	if os.Getenv("ENVIRONMENT_LEVEL") == "development" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Tenta carregar o .env e loga o resultado
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file",
			zap.Error(err),
			zap.String("current_dir", getCurrentDir()))
	} else {
		logger.Info("Successfully loaded .env file")
		// Log some env vars to verify they were loaded
		logger.Debug("Environment variables",
			zap.String("DB_HOST", os.Getenv("DATABASE_HOST")),
			zap.String("DB_NAME", os.Getenv("DATABASE_NAME")),
			zap.String("DB_PORT", os.Getenv("DATABASE_PORT")))
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		logger.Warn("HTTP_PORT not set, using default 8080")
		httpPort = "8080"
	}

	db.InitDB(logger)

	logger.Info("Initializing the application...")
	r := mux.NewRouter()
	rt := routes.NewRouter(logger)

	rt.SetupRouter(r)

	logger.Info("Server starting", zap.String("port", httpPort))

	log.Fatal(http.ListenAndServe(":"+httpPort, handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(r)))
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "error getting current dir"
	}
	return dir
}
