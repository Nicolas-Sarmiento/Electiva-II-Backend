package repository

import (
	"context"

	"ancianato-backend/internal/domain"
	"gorm.io/gorm"
)

type alertTypeRepository struct {
	db *gorm.DB
}

func NewAlertTypeRepository(db *gorm.DB) domain.AlertTypeRepository {
	return &alertTypeRepository{db: db}
}

func (r *alertTypeRepository) Create(ctx context.Context, alertType *domain.AlertType) error {
	return r.db.WithContext(ctx).Create(alertType).Error
}

func (r *alertTypeRepository) GetAll(ctx context.Context) ([]domain.AlertType, error) {
	var types []domain.AlertType
	err := r.db.WithContext(ctx).Find(&types).Error
	return types, err
}

func (r *alertTypeRepository) GetByID(ctx context.Context, id string) (*domain.AlertType, error) {
	var alertType domain.AlertType
	err := r.db.WithContext(ctx).First(&alertType, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &alertType, nil
}

func (r *alertTypeRepository) Update(ctx context.Context, alertType *domain.AlertType) error {
	return r.db.WithContext(ctx).Save(alertType).Error
}

func (r *alertTypeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.AlertType{}).Error
}
