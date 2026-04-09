package repository

import (
	"context"
	"time"

	"ancianato-backend/internal/domain"
	"ancianato-backend/internal/infrastructure/cache"
	"gorm.io/gorm"
)

type patientRepository struct {
	db *gorm.DB
}

// NewPatientRepository retorna la implementación concreta inyectando la conexión de base de datos
func NewPatientRepository(db *gorm.DB) domain.PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(patient).Error; err != nil {
			return err
		}

		// Insert wearables relationships
		for _, w := range patient.WearableDevices {
			pw := domain.PatientWearable{
				PatientID:    patient.ID,
				WearableID:   w.ID,
				AssignedDate: time.Now(),
			}
			if err := tx.Create(&pw).Error; err != nil {
				return err
			}
		}

		cache.AppCache.Delete("patient:ALL")
		return nil
	})
}

func (r *patientRepository) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	cacheKey := "patient:" + id
	if cachedData, found := cache.AppCache.Get(cacheKey); found {
		return cachedData.(*domain.Patient), nil
	}

	var patient domain.Patient
	// Preload nos permite traer las relaciones completas
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("PatientConditions.MedicalCondition").
		Preload("PatientContacts.EmergencyContact").
		Preload("PatientWearables.Wearable").
		First(&patient, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	mapPatientTransientFields(&patient)

	cache.AppCache.Set(cacheKey, &patient, 5*time.Minute)
	return &patient, nil
}

func (r *patientRepository) GetAll(ctx context.Context) ([]domain.Patient, error) {
	cacheKey := "patient:ALL"
	if cachedData, found := cache.AppCache.Get(cacheKey); found {
		return cachedData.([]domain.Patient), nil
	}

	var patients []domain.Patient
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("PatientConditions.MedicalCondition").
		Preload("PatientContacts.EmergencyContact").
		Preload("PatientWearables.Wearable").
		Find(&patients).Error
	if err != nil {
		return nil, err
	}

	for i := range patients {
		mapPatientTransientFields(&patients[i])
	}

	cache.AppCache.Set(cacheKey, patients, 5*time.Minute)
	return patients, nil
}

// mapPatientTransientFields puebla los arrays correspondientes (Allergies, Diseases, etc) 
// para que el JSON de respuesta devuelva la estructura acordada sin anidaciones complejas de la relación.
func mapPatientTransientFields(p *domain.Patient) {
	p.Allergies = make([]domain.MedicalCondition, 0)
	p.Diseases = make([]domain.MedicalCondition, 0)
	p.WearableDevices = make([]domain.Wearable, 0)

	for _, pc := range p.PatientConditions {
		// Asignamos el diagnóstico individual de esta relación
		cond := pc.MedicalCondition
		cond.Diagnostic = pc.Diagnostic
		if cond.AllergenType != "" {
			p.Allergies = append(p.Allergies, cond)
		} else {
			p.Diseases = append(p.Diseases, cond)
		}
	}

	if len(p.PatientContacts) > 0 {
		contact := p.PatientContacts[0].EmergencyContact
		contact.Relationship = p.PatientContacts[0].Relationship
		p.EmergencyContact = &contact
	}

	for _, pw := range p.PatientWearables {
		w := pw.Wearable
		// Evitamos recursividad
		p.WearableDevices = append(p.WearableDevices, w)
	}
}

func (r *patientRepository) Update(ctx context.Context, patient *domain.Patient) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update patient base fields
		if err := tx.Save(patient).Error; err != nil {
			return err
		}

		// Update Wearables (borrar antiguos e insertar nuevos)
		if err := tx.Where("patient_id = ?", patient.ID).Delete(&domain.PatientWearable{}).Error; err != nil {
			return err
		}

		for _, w := range patient.WearableDevices {
			pw := domain.PatientWearable{
				PatientID:    patient.ID,
				WearableID:   w.ID,
				AssignedDate: time.Now(),
			}
			if err := tx.Create(&pw).Error; err != nil {
				return err
			}
		}

		cache.AppCache.Delete("patient:" + patient.ID)
		cache.AppCache.Delete("patient:ALL")
		return nil
	})
}

func (r *patientRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Patient{}).Error
	if err == nil {
		cache.AppCache.Delete("patient:" + id)
		cache.AppCache.Delete("patient:ALL")
	}
	return err
}
