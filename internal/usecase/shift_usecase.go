package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/validation"
	"github.com/google/uuid"
)

type shiftUseCase struct {
	shiftRepo domain.ShiftRepository
}

func NewShiftUseCase(repo domain.ShiftRepository) domain.ShiftUseCase {
	return &shiftUseCase{shiftRepo: repo}
}

// parseShiftTime acepta "HH:MM" o "HH:MM:SS"
func parseShiftTime(t string) (time.Time, error) {
	if parsed, err := time.Parse("15:04:05", t); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("15:04", t); err == nil {
		return parsed, nil
	}
	return time.Time{}, fmt.Errorf("formato inválido '%s': usa HH:MM o HH:MM:SS", t)
}

// validateShiftTimes verifica que end > start y que la diferencia sea >= 1 hora
func validateShiftTimes(startTime, endTime string) error {
	start, err := parseShiftTime(startTime)
	if err != nil {
		return fmt.Errorf("startTime: %w", err)
	}

	end, err := parseShiftTime(endTime)
	if err != nil {
		return fmt.Errorf("endTime: %w", err)
	}

	diff := end.Sub(start)

	if diff <= 0 {
		return errors.New("la hora de fin debe ser posterior a la hora de inicio")
	}
	if diff < time.Hour {
		return errors.New("el turno debe tener al menos 1 hora de duración")
	}

	return nil
}

func (u *shiftUseCase) CreateShift(ctx context.Context, shift *domain.Shift) error {
	if err := validation.Validate.Struct(shift); err != nil {
		return err
	}
	if err := validateShiftTimes(shift.StartTime, shift.EndTime); err != nil {
		return err
	}
	if shift.ID == "" {
		shift.ID = uuid.New().String()
	}
	return u.shiftRepo.Create(ctx, shift)
}

func (u *shiftUseCase) GetAllShifts(ctx context.Context) ([]domain.Shift, error) {
	return u.shiftRepo.GetAll(ctx)
}

func (u *shiftUseCase) GetShift(ctx context.Context, id string) (*domain.Shift, error) {
	return u.shiftRepo.GetByID(ctx, id)
}

func (u *shiftUseCase) UpdateShift(ctx context.Context, shift *domain.Shift) error {
	_, err := u.shiftRepo.GetByID(ctx, shift.ID)
	if err != nil {
		return errors.New("turno no encontrado")
	}
	if err := validateShiftTimes(shift.StartTime, shift.EndTime); err != nil {
		return err
	}
	return u.shiftRepo.Update(ctx, shift)
}

func (u *shiftUseCase) DeleteShift(ctx context.Context, id string) error {
	_, err := u.shiftRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("turno no encontrado")
	}
	return u.shiftRepo.Delete(ctx, id)
}
