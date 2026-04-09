package usecase

import (
	"context"
	"errors"
	"testing"

	"ancianato-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AlertTypeRepositoryMock struct {
	mock.Mock
}

func (m *AlertTypeRepositoryMock) Create(ctx context.Context, alertType *domain.AlertType) error {
	return m.Called(ctx, alertType).Error(0)
}
func (m *AlertTypeRepositoryMock) GetByID(ctx context.Context, id string) (*domain.AlertType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.AlertType), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *AlertTypeRepositoryMock) GetAll(ctx context.Context) ([]domain.AlertType, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.AlertType), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *AlertTypeRepositoryMock) Update(ctx context.Context, alertType *domain.AlertType) error {
	return m.Called(ctx, alertType).Error(0)
}
func (m *AlertTypeRepositoryMock) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreateAlertType_Success(t *testing.T) {
	repo := new(AlertTypeRepositoryMock)
	useCase := NewAlertTypeUseCase(repo)

	at := &domain.AlertType{
		Name: "Caída",
		Code: "FALL",
	}

	repo.On("Create", mock.Anything, at).Return(nil)

	err := useCase.CreateAlertType(context.Background(), at)

	assert.NoError(t, err)
	assert.NotEmpty(t, at.ID)
	repo.AssertExpectations(t)
}

func TestUpdateAlertType_NotFound(t *testing.T) {
	repo := new(AlertTypeRepositoryMock)
	useCase := NewAlertTypeUseCase(repo)

	at := &domain.AlertType{
		ID:          "invalid-id",
		Name:        "Fuga",
	}

	repo.On("GetByID", mock.Anything, "invalid-id").Return((*domain.AlertType)(nil), errors.New("not found"))

	err := useCase.UpdateAlertType(context.Background(), at)

	assert.Error(t, err)
	assert.Equal(t, "tipo de alerta no encontrado", err.Error())
	repo.AssertNotCalled(t, "Update")
}
