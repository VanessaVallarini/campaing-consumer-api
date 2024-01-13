package main

import (
	"campaing-comsumer-service/internal/client"
	"campaing-comsumer-service/internal/config"
	"campaing-comsumer-service/internal/db"
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

	db, err := db.NewDb(ctx, cfg)
	if err != nil {
		easyzap.Panic(err)

	}
	defer db.Close()

	awsClient := client.NewAwsClient(ctx, cfg.AwsConfig.Url, cfg.AwsConfig.Region)
	if awsClient == nil {
		easyzap.Panic("failed creating aws client")
	}

	CampaignService := service.NewCampaignService(awsClient)

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

	queue := "queue_campaing"
	//queue := config.GetConfig().AwsConfig.QueueCampaing
	//queueUrl := awsClient.GetQueueURL(ctx, queue)
	CampaignService.Create(ctx, &c, &queue)

	fmt.Println(cfg)
}
