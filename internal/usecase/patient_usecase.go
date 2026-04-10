package usecase

import (
	"context"
	"errors"
	"time"

	"ancianato-backend/internal/infrastructure/validation"

	"ancianato-backend/internal/domain"
	"github.com/google/uuid"
)

type patientUseCase struct {
	patientRepo domain.PatientRepository
	roomRepo    domain.RoomRepository
	deviceRepo  domain.DeviceRepository
}

// NewPatientUseCase inyecta el repositorio en la capa de negocio
func NewPatientUseCase(repo domain.PatientRepository, roomRepo domain.RoomRepository, deviceRepo domain.DeviceRepository) domain.PatientUseCase {
	return &patientUseCase{patientRepo: repo, roomRepo: roomRepo, deviceRepo: deviceRepo}
}

func (u *patientUseCase) CreatePatient(ctx context.Context, patient *domain.Patient) error {
	// 1. Validaciones estructurales y de etiquetas
	if err := validation.Validate.Struct(patient); err != nil {
		return err
	}

	// 1.5 Validaciones de negocio: La fecha de nacimiento no puede estar en el futuro
	if patient.DateOfBirth.After(time.Now()) {
		return errors.New("la fecha de nacimiento no puede ser en el futuro")
	}

	// 1.6 Validar que la habitación asignada existe en la BD
	if _, err := u.roomRepo.GetByID(ctx, patient.RoomID); err != nil {
		return errors.New("la habitación con ID '" + patient.RoomID + "' no existe")
	}

	// 1.7 Validar que los wearables existen
	for _, w := range patient.WearableDevices {
		if _, err := u.deviceRepo.GetByID(ctx, w.ID); err != nil {
			return errors.New("el wearable con ID '" + w.ID + "' no existe")
		}
	}

	// 1.8 Validar que los wearables no estén ya asignados a otro paciente
	if err := u.validateWearablesUniqueness(ctx, "", patient.WearableDevices); err != nil {
		return err
	}

	// 2. Asignar ID si viene vacío (UUID)
	if patient.ID == "" {
		patient.ID = uuid.New().String()
	}

	if patient.EmergencyContact != nil && patient.EmergencyContact.ID == "" {
		patient.EmergencyContact.ID = uuid.New().String()
	}

	for i := range patient.Allergies {
		if patient.Allergies[i].ID == "" {
			patient.Allergies[i].ID = uuid.New().String()
		}
	}
	for i := range patient.Diseases {
		if patient.Diseases[i].ID == "" {
			patient.Diseases[i].ID = uuid.New().String()
		}
	}

	// 3. Guardar usando Repositorio
	err := u.patientRepo.Create(ctx, patient)
	if err != nil {
		return err
	}

	// TODO: Posteriormente aquí agregaremos la llamada a Kafka para emitir el evento "PatientCreated"
	return nil
}

func (u *patientUseCase) GetPatient(ctx context.Context, id string) (*domain.Patient, error) {
	return u.patientRepo.GetByID(ctx, id)
}

func (u *patientUseCase) GetAllPatients(ctx context.Context) ([]domain.Patient, error) {
	return u.patientRepo.GetAll(ctx)
}

func (u *patientUseCase) UpdatePatient(ctx context.Context, patient *domain.Patient) error {
	// Validar que la habitación asignada existe en la BD
	if _, err := u.roomRepo.GetByID(ctx, patient.RoomID); err != nil {
		return errors.New("la habitación con ID '" + patient.RoomID + "' no existe")
	}

	// First check if the patient exists
	_, err := u.patientRepo.GetByID(ctx, patient.ID)
	if err != nil {
		return errors.New("paciente no encontrado")
	}

	// Validar wearables
	for _, w := range patient.WearableDevices {
		if _, err := u.deviceRepo.GetByID(ctx, w.ID); err != nil {
			return errors.New("el wearable con ID '" + w.ID + "' no existe")
		}
	}

	// Validar que los wearables no estén asignados a otro paciente
	if err := u.validateWearablesUniqueness(ctx, patient.ID, patient.WearableDevices); err != nil {
		return err
	}

	if patient.EmergencyContact != nil && patient.EmergencyContact.ID == "" {
		patient.EmergencyContact.ID = uuid.New().String()
	}

	for i := range patient.Allergies {
		if patient.Allergies[i].ID == "" {
			patient.Allergies[i].ID = uuid.New().String()
		}
	}
	for i := range patient.Diseases {
		if patient.Diseases[i].ID == "" {
			patient.Diseases[i].ID = uuid.New().String()
		}
	}

	return u.patientRepo.Update(ctx, patient)
}

func (u *patientUseCase) DeletePatient(ctx context.Context, id string) error {
	_, err := u.patientRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("paciente no encontrado")
	}

	return u.patientRepo.Delete(ctx, id)
}

func (u *patientUseCase) validateWearablesUniqueness(ctx context.Context, patientID string, newWearables []domain.Wearable) error {
	if len(newWearables) == 0 {
		return nil
	}
	allPatients, err := u.patientRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, p := range allPatients {
		if p.ID == patientID {
			continue // Ignorar al paciente actual si es una actualización
		}
		for _, w := range p.WearableDevices {
			for _, newW := range newWearables {
				if w.ID == newW.ID {
					return errors.New("el wearable con ID '" + w.ID + "' ya se encuentra asignado al paciente " + p.FirstName + " " + p.LastName)
				}
			}
		}
	}
	return nil
}
