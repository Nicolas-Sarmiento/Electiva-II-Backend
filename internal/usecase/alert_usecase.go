package usecase

import (
	"context"
	"errors"
	"time"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/validation"
	"github.com/google/uuid"
)

type alertUseCase struct {
	alertRepo     domain.AlertRepository
	patientRepo   domain.PatientRepository
	deviceRepo    domain.DeviceRepository
	alertTypeRepo domain.AlertTypeRepository
}

func NewAlertUseCase(
	repo domain.AlertRepository,
	patientRepo domain.PatientRepository,
	deviceRepo domain.DeviceRepository,
	alertTypeRepo domain.AlertTypeRepository,
) domain.AlertUseCase {
	return &alertUseCase{
		alertRepo:     repo,
		patientRepo:   patientRepo,
		deviceRepo:    deviceRepo,
		alertTypeRepo: alertTypeRepo,
	}
}

func (u *alertUseCase) CreateAlert(ctx context.Context, alert *domain.Alert) error {
	// 1. Validaciones estructurales
	if err := validation.Validate.Struct(alert); err != nil {
		return err
	}

	if alert.PatientID == "" || alert.WearableID == "" {
		return errors.New("patientId y wearableId son requeridos")
	}

	// 2. Validar que el paciente existe
	if _, err := u.patientRepo.GetByID(ctx, alert.PatientID); err != nil {
		return errors.New("el paciente con ID '" + alert.PatientID + "' no existe")
	}

	// 3. Validar que el wearable existe
	if _, err := u.deviceRepo.GetByID(ctx, alert.WearableID); err != nil {
		return errors.New("el wearable con ID '" + alert.WearableID + "' no existe")
	}

	// 4. Validar que el tipo de alerta existe
	if _, err := u.alertTypeRepo.GetByID(ctx, alert.AlertType); err != nil {
		return errors.New("el tipo de alerta con ID '" + alert.AlertType + "' no existe")
	}

	// 5. Asignar UUID y timestamp
	if alert.ID == "" {
		alert.ID = uuid.New().String()
	}

	if alert.CreatedAt.IsZero() {
		alert.CreatedAt = time.Now()
	}

	return u.alertRepo.Create(ctx, alert)
}

func (u *alertUseCase) GetAlert(ctx context.Context, id string) (*domain.Alert, error) {
	return u.alertRepo.GetByID(ctx, id)
}

func (u *alertUseCase) GetAllAlerts(ctx context.Context) ([]domain.Alert, error) {
	return u.alertRepo.GetAll(ctx)
}

func (u *alertUseCase) UpdateAlert(ctx context.Context, alert *domain.Alert) error {
	existingAlert, err := u.alertRepo.GetByID(ctx, alert.ID)
	if err != nil {
		return errors.New("alerta no encontrada")
	}

	// Actualizar solo campos permitidos en este flujo
	existingAlert.AlertStatus = alert.AlertStatus
	existingAlert.AlertLevel = alert.AlertLevel
	//existingAlert.NurseID = alert.NurseID
	existingAlert.ResolvedAt = alert.ResolvedAt

	return u.alertRepo.Update(ctx, existingAlert)
}

func (u *alertUseCase) DeleteAlert(ctx context.Context, id string) error {
	_, err := u.alertRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("alerta no encontrada")
	}
	return u.alertRepo.Delete(ctx, id)
}
