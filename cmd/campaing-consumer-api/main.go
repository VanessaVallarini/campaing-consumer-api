package main

import (
	"campaing-comsumer-service/internal/client"
	"campaing-comsumer-service/internal/config"
	"campaing-comsumer-service/internal/db"
	"campaing-comsumer-service/internal/listener"
	"campaing-comsumer-service/internal/model"
	"campaing-comsumer-service/internal/service"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lockp111/go-easyzap"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	//clients
	db := db.NewDbClient(cfg)
	defer db.Close()

	awsClient := client.NewAwsClient(cfg.AwsConfig.Url, cfg.AwsConfig.Region)
	if awsClient == nil {
		easyzap.Panic("failed creating aws client")
	}

	//services
	campaingService := service.NewCampaignService(db)

	//testes
	c := model.Event{}
	c.UserId = uuid.New()
	c.SlugId = uuid.New()
	c.MerchantId = uuid.New()
	c.Lat = 45.6085
	c.Long = -73.5493
	c.Action = model.EVENT_ACTION_CREATE
	queue := config.GetConfig().AwsConfig.QueueCampaing
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for range nums {
		awsClient.SendMessage(ctx, &c, &queue)
	}

	//listener
	//go
	listener.EventTrackingListener(ctx, awsClient, campaingService, queue)

	fmt.Println(cfg)
}
