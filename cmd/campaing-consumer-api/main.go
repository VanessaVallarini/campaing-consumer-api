package main

import (
	"campaing-comsumer-service/cmd/campaing-consumer-api/health"
	"campaing-comsumer-service/internal/client"
	"campaing-comsumer-service/internal/config"
	"campaing-comsumer-service/internal/listener"
	"campaing-comsumer-service/internal/model"
	"campaing-comsumer-service/internal/repository"
	"campaing-comsumer-service/internal/service"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lockp111/go-easyzap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	ctx := context.Background()

	tracer.Start(tracer.WithRuntimeMetrics())
	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.GoroutineProfile,
		),
	)
	if err != nil {
		easyzap.Fatal(ctx, err, "failed to start profiler")
	}
	defer tracer.Stop()
	defer profiler.Stop()

	//configs
	config.Init()
	cfg := config.GetConfig()
	awsCfg := config.GetAwsConfig()
	dbCfg := config.GetDatabaseConfig()

	//clients
	db := repository.NewPostgresClient(dbCfg)
	defer db.Close()

	awsClient := client.NewAwsClient(awsCfg.Url, awsCfg.Region)
	if awsClient == nil {
		easyzap.Panic("failed creating aws client")
	}

	//repositories
	campaingRepository := repository.NewCampaingRepository(db)
	userRepository := repository.NewUserRepository(db)
	slugRepository := repository.NewSlugRepository(db)
	merchantRepository := repository.NewMerchantRepository(db)

	//services
	campaingService := service.NewCampaignService(campaingRepository, userRepository, slugRepository, merchantRepository)

	//testes
	queue := awsCfg.QueueCampaing
	c := model.Event{}
	c.UserId = uuid.MustParse("c3eeb9b0-051c-4803-b0a4-f6060bcb40d9")
	c.SlugId = uuid.MustParse("f43e580b-ffb2-490d-aea1-b2f0435d624b")
	c.MerchantId = uuid.MustParse("2ed8b772-1714-46de-98ab-c2653bb03d78")
	c.Lat = 45.6085
	c.Long = -73.5493
	c.Action = model.EVENT_ACTION_CREATE
	awsClient.SendMessage(ctx, &c, &queue)

	//listener
	//go
	listener.EventTrackingListener(ctx, awsClient, campaingService, awsCfg.QueueCampaing)

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
