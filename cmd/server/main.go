package main

import (
	"backend/config"
	"backend/handler"
	"backend/httpserver"
	"backend/logger"
	"backend/metrics"
	"context"
	"fmt"
	golog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"

	"github.com/sanity-io/litter"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.InitConfigurations()
	golog.Print(litter.Sdump(config.Configuration))

	logger, err := logger.ConfigureLogging(&cfg.LogConfig)
	if err != nil {
		golog.Fatalf("error initializing logger - %v", err)
	}
	logger.Info("Starting Application...", zap.String("asd", "asd"))

	registerHystrixMetrics()

	// _, err := bigcache.NewInMemCache(config.Configuration.BigCache)
	if err != nil {
		// logger.Fatalw(ctx, "Error creating big cache", zap.Error(err))
	}
	router := httpserver.NewRouter()
	updateHttpRouter(ctx, router)

	hs := httpserver.NewServer(&config.HttpServer{}, router)

	go hs.Start()
	waitForTermination()
	fmt.Println(ctx, "User initiated shutdown::")
}

func registerHystrixMetrics() {
	collector := metrics.InitializePrometheusCollector(metrics.PrometheusCollectorConfig{
		Namespace: "foundation_gateway_service",
	})
	metricCollector.Registry.Register(collector.NewPrometheusCollector)
}

func updateHttpRouter(ctx context.Context, router httpserver.Router) {

	handler.NewLocationAPIHandler(ctx, router)

	//For SLT mode lets add a handler which can gracefully terminate out server.
	if false {
		// logger.Warnw(ctx, "Running app in SLT mode with kill switch at /rest/kill")
		router.AddRoute(httpserver.RouteConfig{
			Path: "/rest/kill",
			Handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				if err != nil {
					// logger.Fatalw(ctx, "error sending signal", zap.Error(err))
				}
			}),
			Methods:    []string{http.MethodGet},
			Instrument: false,
		})
	}
}

func waitForTermination() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt /*os.Kill,*/, syscall.SIGTERM)
	<-ch
}
