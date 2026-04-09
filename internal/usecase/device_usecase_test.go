package usecase

import (
	"context"
	"errors"
	"testing"

	"ancianato-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateDevice_Success(t *testing.T) {
	deviceRepo := new(DeviceRepositoryMock)
	useCase := NewDeviceUseCase(deviceRepo)

	device := &domain.Wearable{
		MacAddress:   "AA:BB:CC:DD:EE:FF",
		BatteryLevel: 100,
		IsActive:     true,
	}

	// 1. Mock de verificación MAC duplicada (retorna error si no existe, lo cual es lo que espera el UseCase para dejar seguir)
	deviceRepo.On("GetByMacAddress", mock.Anything, "AA:BB:CC:DD:EE:FF").Return((*domain.Wearable)(nil), errors.New("not found"))

	// 2. Mock de guardado
	deviceRepo.On("Create", mock.Anything, device).Return(nil)

	err := useCase.CreateDevice(context.Background(), device)

	assert.NoError(t, err)
	deviceRepo.AssertExpectations(t)
	// Verificar que haya generado un ID automáticamente
	assert.NotEmpty(t, device.ID)
}

func TestCreateDevice_DuplicatedMac(t *testing.T) {
	deviceRepo := new(DeviceRepositoryMock)
	useCase := NewDeviceUseCase(deviceRepo)

	device := &domain.Wearable{
		MacAddress:   "AA:BB:CC:DD:EE:FF",
		BatteryLevel: 100,
		IsActive:     true,
	}

	// Mock para simular que YA existe un dispositivo con esa MAC (Retorna un dispositivo sin error)
	deviceRepo.On("GetByMacAddress", mock.Anything, "AA:BB:CC:DD:EE:FF").Return(&domain.Wearable{ID: "existing-id"}, nil)

	err := useCase.CreateDevice(context.Background(), device)

	assert.Error(t, err)
	assert.Equal(t, "la dirección MAC ya está registrada en otro dispositivo", err.Error())
	deviceRepo.AssertNotCalled(t, "Create")
}
