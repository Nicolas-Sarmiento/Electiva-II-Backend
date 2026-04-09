package repository

import (
	"context"
	"time"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/cache"
	"gorm.io/gorm"
)

type alertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) domain.AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) Create(ctx context.Context, alert *domain.Alert) error {
	err := r.db.WithContext(ctx).Create(alert).Error
	if err == nil {
		cache.AppCache.Delete("alert:ALL")
	}
	return err
}

func (r *alertRepository) GetByID(ctx context.Context, id string) (*domain.Alert, error) {
	cacheKey := "alert:" + id
	if cachedData, found := cache.AppCache.Get(cacheKey); found {
		return cachedData.(*domain.Alert), nil
	}

	var alert domain.Alert
	err := r.db.WithContext(ctx).First(&alert, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheKey, &alert, 2*time.Second) // TTL Corto
	return &alert, nil
}

func (r *alertRepository) GetAll(ctx context.Context) ([]domain.Alert, error) {
	cacheKey := "alert:ALL"
	if cachedData, found := cache.AppCache.Get(cacheKey); found {
		return cachedData.([]domain.Alert), nil
	}

	var alerts []domain.Alert
	err := r.db.WithContext(ctx).Find(&alerts).Error
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheKey, alerts, 2*time.Second) // TTL Corto
	return alerts, nil
}

func (r *alertRepository) Update(ctx context.Context, alert *domain.Alert) error {
	err := r.db.WithContext(ctx).Save(alert).Error
	if err == nil {
		cache.AppCache.Delete("alert:" + alert.ID)
		cache.AppCache.Delete("alert:ALL")
	}
	return err
}

func (r *alertRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Delete(&domain.Alert{}, "id = ?", id).Error
	if err == nil {
		cache.AppCache.Delete("alert:" + id)
		cache.AppCache.Delete("alert:ALL")
	}
	return err
}
