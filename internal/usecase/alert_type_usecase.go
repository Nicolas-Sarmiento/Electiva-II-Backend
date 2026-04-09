package usecase

import (
	"context"
	"errors"
	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/validation"
	"github.com/google/uuid"
)

type alertTypeUseCase struct {
	alertTypeRepo domain.AlertTypeRepository
}

func NewAlertTypeUseCase(repo domain.AlertTypeRepository) domain.AlertTypeUseCase {
	return &alertTypeUseCase{alertTypeRepo: repo}
}

func (u *alertTypeUseCase) CreateAlertType(ctx context.Context, alertType *domain.AlertType) error {
	if err := validation.Validate.Struct(alertType); err != nil {
		return err
	}

	if alertType.ID == "" {
		alertType.ID = uuid.New().String()
	}

	return u.alertTypeRepo.Create(ctx, alertType)
}

func (u *alertTypeUseCase) GetAllAlertTypes(ctx context.Context) ([]domain.AlertType, error) {
	return u.alertTypeRepo.GetAll(ctx)
}

func (u *alertTypeUseCase) GetAlertType(ctx context.Context, id string) (*domain.AlertType, error) {
	return u.alertTypeRepo.GetByID(ctx, id)
}

func (u *alertTypeUseCase) UpdateAlertType(ctx context.Context, alertType *domain.AlertType) error {
	_, err := u.alertTypeRepo.GetByID(ctx, alertType.ID)
	if err != nil {
		return errors.New("tipo de alerta no encontrado")
	}
	return u.alertTypeRepo.Update(ctx, alertType)
}

func (u *alertTypeUseCase) DeleteAlertType(ctx context.Context, id string) error {
	_, err := u.alertTypeRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("tipo de alerta no encontrado")
	}
	return u.alertTypeRepo.Delete(ctx, id)
}
