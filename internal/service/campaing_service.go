package service

import (
	"campaing-comsumer-service/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	//"github.com/Vanessa.Vallarini/address-api/pkg/api/proto/v1"
)

type CampaingRepository interface {
	Create(params model.Campaing) error
	Update(params model.Campaing) error
	Delete(param uuid.UUID) error
	GetByMerchantId(param uuid.UUID) (model.Campaing, error)
}

type UserRepository interface {
	GetById(param uuid.UUID) (model.User, error)
}

type SlugRepository interface {
	GetById(param uuid.UUID) (model.Slug, error)
}

type MerchantRepository interface {
	GetById(param uuid.UUID) (model.Merchant, error)
}

type Campaing struct {
	campaingRepository CampaingRepository
	userRepository     UserRepository
	slugRepository     SlugRepository
	merchantRepository MerchantRepository
}

func NewCampaignService(campaingRepository CampaingRepository, userRepository UserRepository, slugRepository SlugRepository, merchantRepository MerchantRepository) *Campaing {
	return &Campaing{
		campaingRepository: campaingRepository,
		userRepository:     userRepository,
		slugRepository:     slugRepository,
		merchantRepository: merchantRepository,
	}
}

func (c Campaing) Handler(ctx context.Context, campaing *model.Event) error {
	if campaing.Action == model.EVENT_ACTION_CREATE {
		return c.create(ctx, campaing)
	}
	if campaing.Action == model.EVENT_ACTION_UPDATE {
		return c.update(ctx, campaing)
	}
	if campaing.Action == model.EVENT_ACTION_DELETE {
		return c.delete(ctx, campaing.Id)
	}
	return fmt.Errorf(fmt.Sprintf("invalid action:%v", campaing.Action))
}

func (c Campaing) create(ctx context.Context, campaing *model.Event) error {
	user, err := c.userRepository.GetById(campaing.UserId)
	if err != nil {
		return err
	}

	slug, err := c.slugRepository.GetById(campaing.SlugId)
	if err != nil {
		return err
	}

	merchant, err := c.merchantRepository.GetById(campaing.MerchantId)
	if err != nil {
		return err
	}

	hasCampaing, err := c.campaingRepository.GetByMerchantId(campaing.MerchantId)
	if err != nil {
		return err
	}

	if user.Id != uuid.Nil && slug.Id != uuid.Nil && merchant.Id != uuid.Nil && !hasCampaing.Active {
		return c.campaingRepository.Create(model.Campaing{
			Id:          uuid.New(),
			UserId:      user.Id,
			SlugId:      slug.Id,
			MerchantId:  merchant.Id,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Active:      true,
			Lat:         campaing.Lat,
			Long:        campaing.Long,
			Clicks:      0,
			Impressions: 0,
		})
	}

	return fmt.Errorf(fmt.Sprintf("campaign create failure. user_id:%v slug_id:%v merchant_id:%v: has_campaing:%v", user.Id, slug.Id, merchant.Id, hasCampaing.Id))
}

func (c Campaing) update(ctx context.Context, campaing *model.Event) error {
	slug, err := c.slugRepository.GetById(campaing.SlugId)
	if err != nil {
		return err
	}

	user, err := c.userRepository.GetById(campaing.UserId)
	if err != nil {
		return err
	}

	hasCampaing, err := c.campaingRepository.GetByMerchantId(campaing.MerchantId)
	if err != nil {
		return err
	}

	merchant, err := c.merchantRepository.GetById(campaing.MerchantId)
	if err != nil {
		return err
	}

	if user.Id != uuid.Nil && slug.Id != uuid.Nil && merchant.Id != uuid.Nil && hasCampaing.Active {
		return c.campaingRepository.Update(model.Campaing{
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

	return fmt.Errorf(fmt.Sprintf("campaign update failure. user_id:%v slug_id:%v merchant_id:%v: has_campaing:%v", user.Id, slug.Id, merchant.Id, hasCampaing.Id))
}

func (c Campaing) delete(ctx context.Context, id uuid.UUID) error {
	return c.campaingRepository.Delete(id)
}
