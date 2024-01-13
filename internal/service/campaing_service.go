package service

import (
	"campaing-comsumer-service/internal/client"
	"campaing-comsumer-service/internal/model"
	"context"

	"github.com/lockp111/go-easyzap"
)

type ICampaignService interface {
	Create(ctx context.Context, campaing *model.Campaing) (int, *string, error)
}

type CampaignService struct {
	awsClient client.IAwsClient
}

func NewCampaignService(awsClient client.IAwsClient) *CampaignService {
	return &CampaignService{
		awsClient: awsClient,
	}
}

func (c CampaignService) Create(ctx context.Context, campaing *model.Campaing, queue *string) error {
	err := c.awsClient.SendMessage(ctx, campaing, queue)
	if err != nil {
		easyzap.Error(ctx, err, "Error to send campaing to queue for merchant_id: %v", campaing.MerchantId)
	}
	return nil
}
