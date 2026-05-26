package http

import (
	"encoding/json"
	"net/http"

	"ancianato-backend/internal/domain"
	"github.com/go-chi/chi/v5"
)

type AlertTypeHandler struct {
	alertTypeUseCase domain.AlertTypeUseCase
}

func NewAlertTypeHandler(r chi.Router, useCase domain.AlertTypeUseCase) {
	handler := &AlertTypeHandler{
		alertTypeUseCase: useCase,
	}

	r.Route("/alert-type", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

func (h *AlertTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var alertType domain.AlertType
	if err := json.NewDecoder(r.Body).Decode(&alertType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.alertTypeUseCase.CreateAlertType(r.Context(), &alertType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode((map[string]string{"message": "tipo de alerta creada", "alertTypeId": alertType.ID}))
}

func (h *AlertTypeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	types, err := h.alertTypeUseCase.GetAllAlertTypes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func (h *AlertTypeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	alertType, err := h.alertTypeUseCase.GetAlertType(r.Context(), id)
	if err != nil {
		http.Error(w, "Tipo de alerta no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alertType)
}

func (h *AlertTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var alertType domain.AlertType
	if err := json.NewDecoder(r.Body).Decode(&alertType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	alertType.ID = id
	err := h.alertTypeUseCase.UpdateAlertType(r.Context(), &alertType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode((map[string]string{"message": "tipo de alerta actualizada", "alertTypeId": alertType.ID}))
}

func (h *AlertTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.alertTypeUseCase.DeleteAlertType(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Tipo de alerta eliminado"})
}
