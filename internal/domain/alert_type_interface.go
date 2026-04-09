package domain

import "context"

type AlertTypeRepository interface {
	Create(ctx context.Context, alertType *AlertType) error
	GetAll(ctx context.Context) ([]AlertType, error)
	GetByID(ctx context.Context, id string) (*AlertType, error)
	Update(ctx context.Context, alertType *AlertType) error
	Delete(ctx context.Context, id string) error
}

type AlertTypeUseCase interface {
	CreateAlertType(ctx context.Context, alertType *AlertType) error
	GetAllAlertTypes(ctx context.Context) ([]AlertType, error)
	GetAlertType(ctx context.Context, id string) (*AlertType, error)
	UpdateAlertType(ctx context.Context, alertType *AlertType) error
	DeleteAlertType(ctx context.Context, id string) error
}
