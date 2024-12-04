package main

import (
	"backend/config"
	"backend/logger"
	golog "log"

	"github.com/sanity-io/litter"
	"go.uber.org/zap"
)

func main() {
	// _ := context.Background()
	cfg := config.InitConfigurations()
	golog.Print(litter.Sdump(config.Configuration))

	logger, err := logger.ConfigureLogging(&cfg.LogConfig)
	if err != nil {
		golog.Fatalf("error initializing logger - %v", err)
	}
	logger.Info("Starting Application...", zap.String("asd", "asd"))

}
