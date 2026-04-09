package domain

import "context"

type ShiftRepository interface {
	Create(ctx context.Context, shift *Shift) error
	GetAll(ctx context.Context) ([]Shift, error)
	GetByID(ctx context.Context, id string) (*Shift, error)
	Update(ctx context.Context, shift *Shift) error
	Delete(ctx context.Context, id string) error
}

type ShiftUseCase interface {
	CreateShift(ctx context.Context, shift *Shift) error
	GetAllShifts(ctx context.Context) ([]Shift, error)
	GetShift(ctx context.Context, id string) (*Shift, error)
	UpdateShift(ctx context.Context, shift *Shift) error
	DeleteShift(ctx context.Context, id string) error
}
