package service

import (
	"campaing-comsumer-service/internal/db"
	"campaing-comsumer-service/internal/model"
	"context"
	"fmt"
	//"github.com/Vanessa.Vallarini/address-api/pkg/api/proto/v1"
)

type ICampaignService interface {
	CreateHandler(ctx context.Context, campaing *model.Event) error
}

type CampaignService struct {
	db db.IDb
}

func NewCampaignService(db db.IDb) *CampaignService {
	return &CampaignService{
		db: db,
	}
}

func (c CampaignService) CreateHandler(ctx context.Context, campaing *model.Event) error {
	fmt.Println(campaing)
	return nil
}
