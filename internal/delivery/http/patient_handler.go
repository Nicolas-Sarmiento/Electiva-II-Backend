package http

import (
	"encoding/json"
	"net/http"

	"ancianato-backend/internal/domain"

	"github.com/go-chi/chi/v5"
)

type PatientHandler struct {
	patientUseCase domain.PatientUseCase
}

// NewPatientHandler inyecta el Caso de Uso y conecta las rutas con los métodos
func NewPatientHandler(r chi.Router, useCase domain.PatientUseCase) {
	handler := &PatientHandler{
		patientUseCase: useCase,
	}

	// Agrupación de rutas
	r.Route("/patients", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

// Create parsea el payload http a struct y lo delega a la capa UseCase
func (h *PatientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var patient domain.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.patientUseCase.CreatePatient(r.Context(), &patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"message": "Paciente creado exitosamente",
			"patientId": patient.ID,
		})

}

func (h *PatientHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	patients, err := h.patientUseCase.GetAllPatients(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

func (h *PatientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	patient, err := h.patientUseCase.GetPatient(r.Context(), id)
	if err != nil {
		http.Error(w, "Paciente no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var patient domain.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	patient.ID = id
	err := h.patientUseCase.UpdatePatient(r.Context(), &patient)
	if err != nil {
		if err.Error() == "paciente no encontrado" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"patientId": patient.ID,
		"message":   "Información del paciente actualizada",
	})
}

func (h *PatientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.patientUseCase.DeletePatient(r.Context(), id)
	if err != nil {
		if err.Error() == "paciente no encontrado" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Paciente eliminado exitosamente",
	})
}
