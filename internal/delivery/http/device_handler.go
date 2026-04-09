package http

import (
	"encoding/json"
	"net/http"

	"ancianato-backend/internal/domain"
	"github.com/go-chi/chi/v5"
)

type DeviceHandler struct {
	deviceUseCase domain.DeviceUseCase
}

func NewDeviceHandler(r chi.Router, useCase domain.DeviceUseCase) {
	handler := &DeviceHandler{
		deviceUseCase: useCase,
	}

	r.Route("/device", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)
		r.Get("/{device_id}", handler.GetByID)
		r.Put("/{device_id}", handler.Update)
		r.Delete("/{device_id}", handler.Delete)
	})
}

func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var device domain.Wearable
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.deviceUseCase.CreateDevice(r.Context(), &device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"wearableId": device.ID,
		"Message":    "Dispositivo registrado exitosamente",
	})
}

func (h *DeviceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceUseCase.GetAllDevices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *DeviceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "device_id")
	device, err := h.deviceUseCase.GetDevice(r.Context(), id)
	if err != nil {
		http.Error(w, "Dispositivo no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (h *DeviceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "device_id")
	var device domain.Wearable
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	device.ID = id
	err := h.deviceUseCase.UpdateDevice(r.Context(), &device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"wearableId": device.ID,
		"Message":    "Información del dispositivo actualizada",
	})
}

func (h *DeviceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "device_id")
	err := h.deviceUseCase.DeleteDevice(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Message": "Dispositivo eliminado exitosamente",
	})
}
