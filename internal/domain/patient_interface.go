package domain

import "context"

// PatientRepository define los métodos que la capa de datos debe implementar (Inversión de Dependencias)
type PatientRepository interface {
	Create(ctx context.Context, patient *Patient) error
	GetByID(ctx context.Context, id string) (*Patient, error)
	GetAll(ctx context.Context) ([]Patient, error)
	Update(ctx context.Context, patient *Patient) error
	Delete(ctx context.Context, id string) error
}

// PatientUseCase define las reglas de negocio para los pacientes
type PatientUseCase interface {
	CreatePatient(ctx context.Context, patient *Patient) error
	GetPatient(ctx context.Context, id string) (*Patient, error)
	GetAllPatients(ctx context.Context) ([]Patient, error)
	UpdatePatient(ctx context.Context, patient *Patient) error
	DeletePatient(ctx context.Context, id string) error
}
