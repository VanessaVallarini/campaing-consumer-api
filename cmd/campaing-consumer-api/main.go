package main

import (
	"campaing-comsumer-service/internal/client"
	"campaing-comsumer-service/internal/config"
	"campaing-comsumer-service/internal/listener"
	"campaing-comsumer-service/internal/model"
	"campaing-comsumer-service/internal/repository"
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
	db := repository.NewPostgresClient(cfg)
	defer db.Close()

	awsClient := client.NewAwsClient(cfg.AwsConfig.Url, cfg.AwsConfig.Region)
	if awsClient == nil {
		easyzap.Panic("failed creating aws client")
	}

	campaingRepository := repository.NewCampaingRepository(db)

	//services
	campaingService := service.NewCampaignService(campaingRepository)

	//testes
	cc := model.Event{}
	cc.Lat = 45.6085
	cc.Long = -73.5493
	cc.Action = model.EVENT_ACTION_CREATE
	queue := config.GetConfig().AwsConfig.QueueCampaing
	awsClient.SendMessage(ctx, &cc, &queue)

	cu := model.Event{}
	cu.Id = uuid.MustParse("79c1066f-b7e2-48ee-a157-48ec7224e7d1")
	cu.UserId = uuid.New()
	cu.SlugId = uuid.New()
	cu.Active = true
	cu.Lat = 45.6085
	cu.Long = -73.5493
	cu.Clicks = 2
	cu.Impressions = 4
	cu.Action = model.EVENT_ACTION_UPDATE
	awsClient.SendMessage(ctx, &cu, &queue)

	cd := model.Event{}
	cd.Id = uuid.MustParse("79c1066f-b7e2-48ee-a157-48ec7224e7d1")
	cd.Action = model.EVENT_ACTION_DELETE
	awsClient.SendMessage(ctx, &cd, &queue)

	//listener
	//go
	listener.EventTrackingListener(ctx, awsClient, campaingService, queue)

	fmt.Println(cfg)
}
