package usecase

import (
	"context"
	"errors"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/validation"
	"github.com/google/uuid"
)

type deviceUseCase struct {
	deviceRepo domain.DeviceRepository
}

func NewDeviceUseCase(repo domain.DeviceRepository) domain.DeviceUseCase {
	return &deviceUseCase{deviceRepo: repo}
}

func (u *deviceUseCase) CreateDevice(ctx context.Context, device *domain.Wearable) error {
	// 1. Validaciones estructurales (Requeridos, formato MAC address, etc)
	if err := validation.Validate.Struct(device); err != nil {
		return err
	}

	if device.BatteryLevel < 0 || device.BatteryLevel > 100 {
		return errors.New("el nivel de batería debe estar entre 0 y 100")
	}

	// 1.5 Validación de negocio: MAC única
	_, err := u.deviceRepo.GetByMacAddress(ctx, device.MacAddress)
	if err == nil {
		// No debe existir, si err == nil, significa que la encontró y la MAC ya está en uso.
		return errors.New("la dirección MAC ya está registrada en otro dispositivo")
	}

	if device.ID == "" {
		device.ID = uuid.New().String()
	}

	return u.deviceRepo.Create(ctx, device)
}

func (u *deviceUseCase) GetDevice(ctx context.Context, id string) (*domain.Wearable, error) {
	return u.deviceRepo.GetByID(ctx, id)
}

func (u *deviceUseCase) GetAllDevices(ctx context.Context) ([]domain.Wearable, error) {
	return u.deviceRepo.GetAll(ctx)
}

func (u *deviceUseCase) UpdateDevice(ctx context.Context, device *domain.Wearable) error {
	_, err := u.deviceRepo.GetByID(ctx, device.ID)
	if err != nil {
		return errors.New("dispositivo no encontrado")
	}
	return u.deviceRepo.Update(ctx, device)
}

func (u *deviceUseCase) DeleteDevice(ctx context.Context, id string) error {
	_, err := u.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("dispositivo no encontrado")
	}
	return u.deviceRepo.Delete(ctx, id)
}
