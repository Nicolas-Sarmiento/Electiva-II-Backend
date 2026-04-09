package usecase

import (
	"context"
	"errors"
	"testing"

	"ancianato-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AlertRepositoryMock struct {
	mock.Mock
}

func (m *AlertRepositoryMock) Create(ctx context.Context, alert *domain.Alert) error {
	return m.Called(ctx, alert).Error(0)
}
func (m *AlertRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Alert, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Alert), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *AlertRepositoryMock) GetAll(ctx context.Context) ([]domain.Alert, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Alert), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *AlertRepositoryMock) Update(ctx context.Context, alert *domain.Alert) error {
	return m.Called(ctx, alert).Error(0)
}
func (m *AlertRepositoryMock) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreateAlert_Success(t *testing.T) {
	alertRepo := new(AlertRepositoryMock)
	patientRepo := new(PatientRepositoryMock)
	deviceRepo := new(DeviceRepositoryMock)
	alertTypeRepo := new(AlertTypeRepositoryMock)

	useCase := NewAlertUseCase(alertRepo, patientRepo, deviceRepo, alertTypeRepo)

	alert := &domain.Alert{
		PatientID:   "patient-1",
		WearableID:  "wearable-1",
		AlertTypeID: "type-1",
		AlertStatus: "ACTIVE",
		AlertLevel:  "HIGH",
	}

	patientRepo.On("GetByID", mock.Anything, "patient-1").Return(&domain.Patient{ID: "patient-1"}, nil)
	deviceRepo.On("GetByID", mock.Anything, "wearable-1").Return(&domain.Wearable{ID: "wearable-1"}, nil)
	alertTypeRepo.On("GetByID", mock.Anything, "type-1").Return(&domain.AlertType{ID: "type-1"}, nil)

	alertRepo.On("Create", mock.Anything, alert).Return(nil)

	err := useCase.CreateAlert(context.Background(), alert)

	assert.NoError(t, err)
	assert.NotEmpty(t, alert.ID)
	assert.False(t, alert.CreatedAt.IsZero())
	
	alertRepo.AssertExpectations(t)
	patientRepo.AssertExpectations(t)
	deviceRepo.AssertExpectations(t)
	alertTypeRepo.AssertExpectations(t)
}

func TestCreateAlert_MissingPatient(t *testing.T) {
	alertRepo := new(AlertRepositoryMock)
	patientRepo := new(PatientRepositoryMock)
	deviceRepo := new(DeviceRepositoryMock)
	alertTypeRepo := new(AlertTypeRepositoryMock)

	useCase := NewAlertUseCase(alertRepo, patientRepo, deviceRepo, alertTypeRepo)

	alert := &domain.Alert{
		PatientID:   "missing-patient",
		WearableID:  "wearable-1",
		AlertTypeID: "type-1",
		AlertStatus: "ACTIVE",
		AlertLevel:  "HIGH",
	}

	patientRepo.On("GetByID", mock.Anything, "missing-patient").Return((*domain.Patient)(nil), errors.New("not found"))

	err := useCase.CreateAlert(context.Background(), alert)

	assert.Error(t, err)
	assert.Equal(t, "el paciente con ID 'missing-patient' no existe", err.Error())
	alertRepo.AssertNotCalled(t, "Create")
}
