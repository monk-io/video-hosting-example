package logger

import (
	"os"

	"go.uber.org/zap"
)

func New() *zap.Logger {
	env := os.Getenv("GO_ENV")

	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return logger
}
