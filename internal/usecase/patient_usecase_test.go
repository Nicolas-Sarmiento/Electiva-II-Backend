package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/validation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// PatientRepositoryMock Mock de PatientRepository
type PatientRepositoryMock struct {
	mock.Mock
}

func (m *PatientRepositoryMock) Create(ctx context.Context, patient *domain.Patient) error {
	args := m.Called(ctx, patient)
	return args.Error(0)
}

func (m *PatientRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Patient), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *PatientRepositoryMock) GetAll(ctx context.Context) ([]domain.Patient, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Patient), args.Error(1)
}

func (m *PatientRepositoryMock) Update(ctx context.Context, patient *domain.Patient) error {
	args := m.Called(ctx, patient)
	return args.Error(0)
}

func (m *PatientRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// RoomRepositoryMock Mock de RoomRepository
type RoomRepositoryMock struct {
	mock.Mock
}

func (m *RoomRepositoryMock) Create(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *RoomRepositoryMock) GetAll(ctx context.Context) ([]domain.Room, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Room), args.Error(1)
}

func (m *RoomRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Room), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RoomRepositoryMock) Update(ctx context.Context, room *domain.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *RoomRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RoomRepositoryMock) GetByNumberAndPavilion(ctx context.Context, number, pavilion string) (*domain.Room, error) {
	args := m.Called(ctx, number, pavilion)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Room), args.Error(1)
	}
	return nil, args.Error(1)
}

// DeviceRepositoryMock
type DeviceRepositoryMock struct {
	mock.Mock
}

func (m *DeviceRepositoryMock) Create(ctx context.Context, device *domain.Wearable) error {
	return m.Called(ctx, device).Error(0)
}
func (m *DeviceRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Wearable, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Wearable), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *DeviceRepositoryMock) GetByMacAddress(ctx context.Context, mac string) (*domain.Wearable, error) {
	args := m.Called(ctx, mac)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Wearable), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *DeviceRepositoryMock) GetAll(ctx context.Context) ([]domain.Wearable, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Wearable), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *DeviceRepositoryMock) Update(ctx context.Context, device *domain.Wearable) error {
	return m.Called(ctx, device).Error(0)
}
func (m *DeviceRepositoryMock) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func init() {
	validation.InitValidator() // Necesario para que funcione validation.Validate
}

func TestCreatePatient_Success(t *testing.T) {
	// 1. Arrange Configuración
	patientRepo := new(PatientRepositoryMock)
	roomRepo := new(RoomRepositoryMock)
	deviceRepo := new(DeviceRepositoryMock)

	useCase := NewPatientUseCase(patientRepo, roomRepo, deviceRepo)

	// Datos de prueba correctos
	patient := &domain.Patient{
		FirstName:   "Juan",
		LastName:    "Perez",
		DateOfBirth: time.Now().Add(-time.Hour * 8760 * 70), // 70 años aprox
		RoomID:      "room-123",
		EmergencyContact: &domain.EmergencyContact{
			FirstName: "Maria",
			LastName:  "Perez",
			Phone:     "12345678",
			Mail:      "maria@test.com",
			Relationship: "Madre",
		},
		WearableDevices: []domain.Wearable{
			{ID: "wearable-1"},
		},
	}

	// 2. Arrange Expecativas (Mocks config)
	// Cuando se pregunte si la habitación "room-123" existe, devolver una válida, error=nil
	roomRepo.On("GetByID", mock.Anything, "room-123").Return(&domain.Room{ID: "room-123"}, nil)
	
	// Predecir la búsqueda del wearable
	deviceRepo.On("GetByID", mock.Anything, "wearable-1").Return(&domain.Wearable{ID: "wearable-1"}, nil)

	// Esperamos que llame al guardado sin error
	patientRepo.On("Create", mock.Anything, patient).Return(nil)

	// 3. Act
	err := useCase.CreatePatient(context.Background(), patient)

	// 4. Assert
	assert.NoError(t, err)
	patientRepo.AssertExpectations(t)
	roomRepo.AssertExpectations(t)
	deviceRepo.AssertExpectations(t)
}

func TestCreatePatient_InvalidRoom(t *testing.T) {
	patientRepo := new(PatientRepositoryMock)
	roomRepo := new(RoomRepositoryMock)
	deviceRepo := new(DeviceRepositoryMock)
	useCase := NewPatientUseCase(patientRepo, roomRepo, deviceRepo)

	patient := &domain.Patient{
		FirstName:   "Juan",
		LastName:    "Perez",
		DateOfBirth: time.Now().Add(-time.Hour * 8760 * 70),
		RoomID:      "invalid-room", // Room invalida
		EmergencyContact: &domain.EmergencyContact{
			FirstName: "Maria", LastName: "Perez", Phone: "12345678", Mail: "maria@test.com", Relationship: "Madre",
		},
	}

	// Mock para cuando busquen la habitacion (Devuelve error)
	roomRepo.On("GetByID", mock.Anything, "invalid-room").Return((*domain.Room)(nil), errors.New("not found"))

	err := useCase.CreatePatient(context.Background(), patient)

	assert.Error(t, err)
	assert.Equal(t, "la habitación con ID 'invalid-room' no existe", err.Error())

	// Verificamos que paciente NO fue guardado (Assert no Create)
	patientRepo.AssertNotCalled(t, "Create")
}
