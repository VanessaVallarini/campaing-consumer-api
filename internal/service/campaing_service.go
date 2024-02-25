package service

import (
	"campaing-comsumer-service/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
	//"github.com/Vanessa.Vallarini/address-api/pkg/api/proto/v1"
)

type CampaingRepository interface {
	Create(ctx context.Context, tx pgx.Tx, params model.Campaing) error
	Update(ctx context.Context, tx pgx.Tx, params model.Campaing) error
	Delete(ctx context.Context, tx pgx.Tx, param uuid.UUID) error
	GetByMerchantId(ctx context.Context, param uuid.UUID) (model.Campaing, error)
	GetById(ctx context.Context, param uuid.UUID) (model.Campaing, error)
}

type UserRepository interface {
	GetById(ctx context.Context, param uuid.UUID) (model.User, error)
}

type SlugRepository interface {
	GetById(ctx context.Context, param uuid.UUID) (model.Slug, error)
}

type Transactios interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}

type MerchantRepository interface {
	GetById(ctx context.Context, param uuid.UUID) (model.Merchant, error)
}

type Campaing struct {
	transactios        Transactios
	campaingRepository CampaingRepository
	userRepository     UserRepository
	slugRepository     SlugRepository
	merchantRepository MerchantRepository
}

func NewCampaignService(transactios Transactios, campaingRepository CampaingRepository, userRepository UserRepository, slugRepository SlugRepository, merchantRepository MerchantRepository) *Campaing {
	return &Campaing{
		transactios:        transactios,
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
	user, err := c.userRepository.GetById(ctx, campaing.UserId)
	if err != nil {
		return err
	}

	slug, err := c.slugRepository.GetById(ctx, campaing.SlugId)
	if err != nil {
		return err
	}

	merchant, err := c.merchantRepository.GetById(ctx, campaing.MerchantId)
	if err != nil {
		return err
	}

	hasCampaing, err := c.campaingRepository.GetByMerchantId(ctx, campaing.MerchantId)
	if err != nil {
		return err
	}

	tx, err := c.transactios.Begin(ctx)
	if err != nil {
		return err
	}

	if user.Id == uuid.Nil || slug.Id == uuid.Nil || merchant.Id == uuid.Nil || hasCampaing.Active {
		errWrap := errors.Wrapf(err, "invalid params: userId %v slugId %v merchantId %v hasCampaingActive %v", user.Id, slug.Id, merchant.Id, !hasCampaing.Active)
		easyzap.Warnf("invalid params: userId %v slugId %v merchantId %v hasCampaing %v", user.Id, slug.Id, merchant.Id, !hasCampaing.Active)
		return errWrap
	}

	errCreating := c.campaingRepository.Create(ctx, tx, model.Campaing{
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
	if errCreating != nil {
		errRollback := c.transactios.Rollback(ctx, tx)
		if errRollback != nil {
			return errRollback
		}
		return errCreating
	}

	errCommit := c.transactios.Commit(ctx, tx)
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (c Campaing) update(ctx context.Context, campaing *model.Event) error {
	user, err := c.userRepository.GetById(ctx, campaing.UserId)
	if err != nil {
		return err
	}

	slug, err := c.slugRepository.GetById(ctx, campaing.SlugId)
	if err != nil {
		return err
	}

	merchant, err := c.merchantRepository.GetById(ctx, campaing.MerchantId)
	if err != nil {
		return err
	}

	hasCampaing, err := c.campaingRepository.GetById(ctx, campaing.Id)
	if err != nil {
		return err
	}

	tx, err := c.transactios.Begin(ctx)
	if err != nil {
		return err
	}

	if user.Id == uuid.Nil || slug.Id == uuid.Nil || merchant.Id == uuid.Nil || len(hasCampaing.Id) == 0 {
		errWrap := errors.Wrapf(err, "invalid params: userId %v slugId %v merchantId %v hasCampaing %v", user.Id, slug.Id, merchant.Id, !hasCampaing.Active)
		easyzap.Warnf("invalid params: userId %v slugId %v merchantId %v hasCampaing %v", user.Id, slug.Id, merchant.Id, !hasCampaing.Active)
		return errWrap
	}

	errUpdating := c.campaingRepository.Update(ctx, tx, model.Campaing{
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
	if errUpdating != nil {
		errRollback := c.transactios.Rollback(ctx, tx)
		if errRollback != nil {
			return errRollback
		}
		return errUpdating
	}

	errCommit := c.transactios.Commit(ctx, tx)
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (c Campaing) delete(ctx context.Context, id uuid.UUID) error {
	tx, err := c.transactios.Begin(ctx)
	if err != nil {
		return err
	}

	errDeleting := c.campaingRepository.Delete(ctx, tx, id)
	if errDeleting != nil {
		errRollback := c.transactios.Rollback(ctx, tx)
		if errRollback != nil {
			return errRollback
		}
		return errDeleting
	}

	errCommit := c.transactios.Commit(ctx, tx)
	if errCommit != nil {
		return errCommit
	}

	return nil
}
