package usecase

import (
	"context"
	"testing"

	"ancianato-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ShiftRepositoryMock struct {
	mock.Mock
}

func (m *ShiftRepositoryMock) Create(ctx context.Context, shift *domain.Shift) error {
	return m.Called(ctx, shift).Error(0)
}
func (m *ShiftRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Shift, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Shift), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *ShiftRepositoryMock) GetAll(ctx context.Context) ([]domain.Shift, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Shift), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *ShiftRepositoryMock) Update(ctx context.Context, shift *domain.Shift) error {
	return m.Called(ctx, shift).Error(0)
}
func (m *ShiftRepositoryMock) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestCreateShift_Success(t *testing.T) {
	shiftRepo := new(ShiftRepositoryMock)
	useCase := NewShiftUseCase(shiftRepo)

	shift := &domain.Shift{
		Name:      "Mañana",
		StartTime: "06:00:00",
		EndTime:   "14:00:00",
	}

	shiftRepo.On("Create", mock.Anything, shift).Return(nil)

	err := useCase.CreateShift(context.Background(), shift)

	assert.NoError(t, err)
	assert.NotEmpty(t, shift.ID)
	shiftRepo.AssertExpectations(t)
}

func TestCreateShift_InvalidDuration(t *testing.T) {
	shiftRepo := new(ShiftRepositoryMock)
	useCase := NewShiftUseCase(shiftRepo)

	// Falla porque debe ser >= 1 hora
	shift := &domain.Shift{
		Name:      "Invalido",
		StartTime: "10:00",
		EndTime:   "10:30",
	}

	err := useCase.CreateShift(context.Background(), shift)

	assert.Error(t, err)
	assert.Equal(t, "el turno debe tener al menos 1 hora de duración", err.Error())
	shiftRepo.AssertNotCalled(t, "Create")
}

func TestCreateShift_InvalidFormat(t *testing.T) {
	shiftRepo := new(ShiftRepositoryMock)
	useCase := NewShiftUseCase(shiftRepo)

	shift := &domain.Shift{
		Name:      "Invalido",
		StartTime: "XX:YY",
		EndTime:   "ZZ:WW",
	}

	err := useCase.CreateShift(context.Background(), shift)

	assert.Error(t, err)
	// Comprobar que en el error diga formato inválido
	assert.Contains(t, err.Error(), "formato inválido")
	shiftRepo.AssertNotCalled(t, "Create")
}
