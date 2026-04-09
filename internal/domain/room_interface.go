package domain

import "context"

type RoomRepository interface {
	Create(ctx context.Context, room *Room) error
	GetAll(ctx context.Context) ([]Room, error)
	GetByID(ctx context.Context, id string) (*Room, error)
	Update(ctx context.Context, room *Room) error
	Delete(ctx context.Context, id string) error
	GetByNumberAndPavilion(ctx context.Context, number, pavilion string) (*Room, error)
}

type RoomUseCase interface {
	CreateRoom(ctx context.Context, room *Room) error
	GetAllRooms(ctx context.Context) ([]Room, error)
	GetRoom(ctx context.Context, id string) (*Room, error)
	UpdateRoom(ctx context.Context, room *Room) error
	DeleteRoom(ctx context.Context, id string) error
}
