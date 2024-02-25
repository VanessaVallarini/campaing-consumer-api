package main

import (
	"campaing-comsumer-service/cmd/campaing-consumer-api/health"
	"campaing-comsumer-service/internal/client"
	"campaing-comsumer-service/internal/config"
	"campaing-comsumer-service/internal/listener"
	"campaing-comsumer-service/internal/metrics"
	"campaing-comsumer-service/internal/repository"
	"campaing-comsumer-service/internal/service"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/lockp111/go-easyzap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()

	//configs
	config.Init()
	cfg := config.GetConfig()
	awsCfg := config.GetAwsConfig()
	dbCfg := config.GetDatabaseConfig()

	//metrics
	metrics := metrics.NewMetrics()

	//database
	easyzap.Info(ctx, "Creating pool")
	db, err := repository.CreatePool(ctx, &dbCfg)
	if err != nil {
		easyzap.Fatal(ctx, err, "error when start pool")
	}
	err = db.Ping(ctx)
	if err != nil {
		easyzap.Fatal(ctx, err, "failed to ping pool")
	}

	//clients
	awsClient := client.NewAwsClient(metrics, awsCfg.Url, awsCfg.Region)
	if awsClient == nil {
		easyzap.Panic("failed creating aws client")
	}

	//repositories
	transactios := repository.NewTransactionDao(metrics, db)
	campaingRepository := repository.NewCampaingRepository(metrics, db)
	userRepository := repository.NewUserRepository(metrics, db)
	slugRepository := repository.NewSlugRepository(metrics, db)
	merchantRepository := repository.NewMerchantRepository(metrics, db)

	//services
	campaingService := service.NewCampaignService(transactios, campaingRepository, userRepository, slugRepository, merchantRepository)

	//listener
	go listener.EventTrackingListener(ctx, metrics, awsClient, campaingService, awsCfg.QueueCampaing)

	meta := echo.New()
	meta.HideBanner = true
	meta.HidePort = true

	healthChecker := health.HealthMonitor{
		Checked: map[string]health.HealthCheck{},
	}

	healthController := health.HealthController{
		Monitor: healthChecker,
	}

	health.EchoRegister(meta, healthController, "/health", "/env")
	meta.GET("/prometheus", echo.WrapHandler(promhttp.Handler()))

	go func() {
		easyzap.Info(ctx, "Starting metadata server at "+cfg.HealthPort)
		err := meta.Start(cfg.HealthPort)
		easyzap.Fatal(ctx, err, "failed to start server")
	}()

	// Listen for system signals to gracefully stop the application
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	switch <-signalChannel {
	case os.Interrupt:
		easyzap.Info(ctx, "Received SIGINT, stopping...")
	case syscall.SIGTERM:
		easyzap.Info(ctx, "Received SIGTERM, stopping...")
	}

}
