package domain

import "context"

type DeviceRepository interface {
	Create(ctx context.Context, device *Wearable) error
	GetByID(ctx context.Context, id string) (*Wearable, error)
	GetByMacAddress(ctx context.Context, mac string) (*Wearable, error)
	GetAll(ctx context.Context) ([]Wearable, error)
	Update(ctx context.Context, device *Wearable) error
	Delete(ctx context.Context, id string) error
}

type DeviceUseCase interface {
	CreateDevice(ctx context.Context, device *Wearable) error
	GetDevice(ctx context.Context, id string) (*Wearable, error)
	GetAllDevices(ctx context.Context) ([]Wearable, error)
	UpdateDevice(ctx context.Context, device *Wearable) error
	DeleteDevice(ctx context.Context, id string) error
}
