package usecase

import (
	"context"
	"errors"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/validation"
	"github.com/google/uuid"
)

type roomUseCase struct {
	roomRepo domain.RoomRepository
}

func NewRoomUseCase(repo domain.RoomRepository) domain.RoomUseCase {
	return &roomUseCase{roomRepo: repo}
}

func (u *roomUseCase) CreateRoom(ctx context.Context, room *domain.Room) error {
	if err := validation.Validate.Struct(room); err != nil {
		return err
	}

	existingRoom, err := u.roomRepo.GetByNumberAndPavilion(ctx, room.RoomNumber, room.RoomPavilion)
	if err == nil && existingRoom != nil {
		return errors.New("la habitación ya existe en este pabellón")
	}

	if room.ID == "" {
		room.ID = uuid.New().String()
	}

	return u.roomRepo.Create(ctx, room)
}

func (u *roomUseCase) GetAllRooms(ctx context.Context) ([]domain.Room, error) {
	return u.roomRepo.GetAll(ctx)
}

func (u *roomUseCase) GetRoom(ctx context.Context, id string) (*domain.Room, error) {
	return u.roomRepo.GetByID(ctx, id)
}

func (u *roomUseCase) UpdateRoom(ctx context.Context, room *domain.Room) error {
	existing, err := u.roomRepo.GetByID(ctx, room.ID)
	if err != nil {
		return errors.New("habitación no encontrada")
	}

	// Verificar duplicados si cambian el número o pabellón
	if existing.RoomNumber != room.RoomNumber || existing.RoomPavilion != room.RoomPavilion {
		dupRoom, err := u.roomRepo.GetByNumberAndPavilion(ctx, room.RoomNumber, room.RoomPavilion)
		if err == nil && dupRoom != nil && dupRoom.ID != room.ID {
			return errors.New("ya existe otra habitación con ese número en este pabellón")
		}
	}

	return u.roomRepo.Update(ctx, room)
}

func (u *roomUseCase) DeleteRoom(ctx context.Context, id string) error {
	_, err := u.roomRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("habitación no encontrada")
	}
	return u.roomRepo.Delete(ctx, id)
}
