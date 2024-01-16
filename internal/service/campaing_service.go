package service

import (
	"campaing-comsumer-service/internal/db"
	"campaing-comsumer-service/internal/model"
	"context"
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
	return nil
}
