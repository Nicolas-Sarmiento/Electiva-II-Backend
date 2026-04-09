package usecase

import (
	"context"
	"errors"
	"testing"

	"ancianato-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRoom_Success(t *testing.T) {
	roomRepo := new(RoomRepositoryMock)
	useCase := NewRoomUseCase(roomRepo)

	room := &domain.Room{
		Floor:        1,
		RoomNumber:   "101",
		RoomPavilion: "Norte",
	}

	// Simular que no existe ninguno con ese número y pabellón
	roomRepo.On("GetByNumberAndPavilion", mock.Anything, "101", "Norte").Return((*domain.Room)(nil), errors.New("not found"))

	// Permitir creación
	roomRepo.On("Create", mock.Anything, room).Return(nil)

	err := useCase.CreateRoom(context.Background(), room)

	assert.NoError(t, err)
	assert.NotEmpty(t, room.ID) // Generado uuid
	roomRepo.AssertExpectations(t)
}

func TestCreateRoom_Duplicated(t *testing.T) {
	roomRepo := new(RoomRepositoryMock)
	useCase := NewRoomUseCase(roomRepo)

	room := &domain.Room{
		Floor:        1,
		RoomNumber:   "101",
		RoomPavilion: "Norte",
	}

	// Simular que YA existe!
	roomRepo.On("GetByNumberAndPavilion", mock.Anything, "101", "Norte").Return(&domain.Room{ID: "some-id"}, nil)

	err := useCase.CreateRoom(context.Background(), room)

	assert.Error(t, err)
	assert.Equal(t, "la habitación ya existe en este pabellón", err.Error())
	roomRepo.AssertNotCalled(t, "Create")
}
