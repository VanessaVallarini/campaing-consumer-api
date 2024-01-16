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
	"time"

	"github.com/google/uuid"
	"github.com/lockp111/go-easyzap"
)

func main() {
	ctx := context.Background()

	cfg := config.GetConfig()

	db := db.NewDb(cfg)
	defer db.Close()

	awsClient := client.NewAwsClient(cfg.AwsConfig.Url, cfg.AwsConfig.Region)
	if awsClient == nil {
		easyzap.Panic("failed creating aws client")
	}
	campaingService := service.NewCampaignService(db)

	c := model.Campaing{}
	c.Id = uuid.New()
	c.UserId = uuid.New()
	c.SlugId = uuid.New()
	c.MerchantId = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Active = true
	c.Lat = 45.6085
	c.Long = -73.5493
	c.Clicks = 0
	c.Impressions = 0

	queue := config.GetConfig().AwsConfig.QueueCampaing
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for range nums {
		awsClient.SendMessage(ctx, &c, &queue)
	}

	listener.EventTrackingListener(ctx, awsClient, campaingService, queue)

	fmt.Println(cfg)
}
