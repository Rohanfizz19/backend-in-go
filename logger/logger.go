package logger

import (
	"backend/config"

	"go.uber.org/zap"
)

func ConfigureLogging(cfg *config.Log) (*zap.Logger, error) {
	var logger *zap.Logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	logger.Info("Logger Initialized!")
	return logger, nil
}
