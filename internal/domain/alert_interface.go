package domain

import "context"

type AlertRepository interface {
	Create(ctx context.Context, alert *Alert) error
	GetByID(ctx context.Context, id string) (*Alert, error)
	GetAll(ctx context.Context) ([]Alert, error)
	Update(ctx context.Context, alert *Alert) error
	Delete(ctx context.Context, id string) error
	HasActiveAlert(ctx context.Context, patientID string, alertType string) (bool, error)
}

type AlertUseCase interface {
	CreateAlert(ctx context.Context, alert *Alert) error
	GetAlert(ctx context.Context, id string) (*Alert, error)
	GetAllAlerts(ctx context.Context) ([]Alert, error)
	UpdateAlert(ctx context.Context, alert *Alert) error
	DeleteAlert(ctx context.Context, id string) error
}
