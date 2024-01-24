package service

import (
	"campaing-comsumer-service/internal/model"
	"context"
	"time"

	"github.com/google/uuid"
	//"github.com/Vanessa.Vallarini/address-api/pkg/api/proto/v1"
)

type CampaingRepository interface {
	Create(params model.Campaing) error
	Update(params model.Campaing) error
	Delete(param uuid.UUID) error
}

type Campaing struct {
	repository CampaingRepository
}

func NewCampaignService(repository CampaingRepository) *Campaing {
	return &Campaing{
		repository: repository,
	}
}

func (c Campaing) CampaingHandler(ctx context.Context, campaing *model.Event) error {
	if campaing.Action == model.EVENT_ACTION_CREATE {
		return c.repository.Create(model.Campaing{
			Id:          uuid.New(),
			UserId:      uuid.New(),
			SlugId:      uuid.New(),
			MerchantId:  uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Active:      true,
			Lat:         campaing.Lat,
			Long:        campaing.Long,
			Clicks:      0,
			Impressions: 0,
		})
	}
	if campaing.Action == model.EVENT_ACTION_UPDATE {
		return c.repository.Create(model.Campaing{
			Id:          campaing.Id,
			UserId:      campaing.UserId,
			SlugId:      campaing.SlugId,
			UpdatedAt:   time.Now(),
			Active:      campaing.Active,
			Lat:         campaing.Lat,
			Long:        campaing.Long,
			Clicks:      campaing.Clicks,
			Impressions: campaing.Impressions,
		})
	}

	return c.repository.Delete(campaing.Id)
}
