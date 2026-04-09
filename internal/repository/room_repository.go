package repository

import (
	"context"

	"ancianato-backend/internal/domain"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *roomRepository) GetAll(ctx context.Context) ([]domain.Room, error) {
	var rooms []domain.Room
	err := r.db.WithContext(ctx).Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	var room domain.Room
	err := r.db.WithContext(ctx).First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) Update(ctx context.Context, room *domain.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *roomRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Room{}).Error
}

func (r *roomRepository) GetByNumberAndPavilion(ctx context.Context, number, pavilion string) (*domain.Room, error) {
	var room domain.Room
	err := r.db.WithContext(ctx).Where("room_number = ? AND room_pavilion = ?", number, pavilion).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}
