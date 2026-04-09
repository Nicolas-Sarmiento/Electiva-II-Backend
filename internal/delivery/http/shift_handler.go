package http

import (
	"encoding/json"
	"net/http"

	"ancianato-backend/internal/domain"
	"github.com/go-chi/chi/v5"
)

type ShiftHandler struct {
	shiftUseCase domain.ShiftUseCase
}

func NewShiftHandler(r chi.Router, useCase domain.ShiftUseCase) {
	handler := &ShiftHandler{
		shiftUseCase: useCase,
	}

	r.Route("/shifts", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

func (h *ShiftHandler) Create(w http.ResponseWriter, r *http.Request) {
	var shift domain.Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.shiftUseCase.CreateShift(r.Context(), &shift)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shift)
}

func (h *ShiftHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	shifts, err := h.shiftUseCase.GetAllShifts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shifts)
}

func (h *ShiftHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	shift, err := h.shiftUseCase.GetShift(r.Context(), id)
	if err != nil {
		http.Error(w, "Turno no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shift)
}

func (h *ShiftHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var shift domain.Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shift.ID = id
	err := h.shiftUseCase.UpdateShift(r.Context(), &shift)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shift)
}

func (h *ShiftHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.shiftUseCase.DeleteShift(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Turno eliminado"})
}
