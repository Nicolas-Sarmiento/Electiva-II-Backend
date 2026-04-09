package repository

import (
	"context"

	"ancianato-backend/internal/domain"
	"gorm.io/gorm"
)

type shiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) domain.ShiftRepository {
	return &shiftRepository{db: db}
}

func (r *shiftRepository) Create(ctx context.Context, shift *domain.Shift) error {
	return r.db.WithContext(ctx).Create(shift).Error
}

func (r *shiftRepository) GetAll(ctx context.Context) ([]domain.Shift, error) {
	var shifts []domain.Shift
	err := r.db.WithContext(ctx).Find(&shifts).Error
	return shifts, err
}

func (r *shiftRepository) GetByID(ctx context.Context, id string) (*domain.Shift, error) {
	var shift domain.Shift
	err := r.db.WithContext(ctx).First(&shift, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepository) Update(ctx context.Context, shift *domain.Shift) error {
	return r.db.WithContext(ctx).Save(shift).Error
}

func (r *shiftRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Shift{}).Error
}
