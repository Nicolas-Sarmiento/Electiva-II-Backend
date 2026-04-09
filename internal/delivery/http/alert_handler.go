package http

import (
	"encoding/json"
	"net/http"

	"ancianato-backend/internal/domain"
	"github.com/go-chi/chi/v5"
)

type AlertHandler struct {
	alertUseCase domain.AlertUseCase
}

func NewAlertHandler(r chi.Router, useCase domain.AlertUseCase) {
	handler := &AlertHandler{
		alertUseCase: useCase,
	}

	r.Route("/alert", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)
		r.Get("/{alert_id}", handler.GetByID)
		r.Put("/{alert_id}", handler.Update)
		r.Delete("/{alert_id}", handler.Delete)
	})
}

func (h *AlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	var alert domain.Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.alertUseCase.CreateAlert(r.Context(), &alert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"alert_id":  alert.ID,
		"createdAt": alert.CreatedAt,
		"Message":   "Alerta creada exitosamente",
	})
}

func (h *AlertHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.alertUseCase.GetAllAlerts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func (h *AlertHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "alert_id")
	alert, err := h.alertUseCase.GetAlert(r.Context(), id)
	if err != nil {
		http.Error(w, "Alerta no encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

func (h *AlertHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "alert_id")
	var alert domain.Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	alert.ID = id
	err := h.alertUseCase.UpdateAlert(r.Context(), &alert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"alertId": alert.ID,
		"Message": "Alerta actualizada",
	})
}

func (h *AlertHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "alert_id")
	err := h.alertUseCase.DeleteAlert(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Message": "Alerta eliminada exitosamente",
	})
}
