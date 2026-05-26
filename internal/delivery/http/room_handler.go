package http

import (
	"encoding/json"
	"net/http"

	"ancianato-backend/internal/domain"
	"github.com/go-chi/chi/v5"
)

type RoomHandler struct {
	roomUseCase domain.RoomUseCase
}

func NewRoomHandler(r chi.Router, useCase domain.RoomUseCase) {
	handler := &RoomHandler{
		roomUseCase: useCase,
	}

	r.Route("/room", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var room domain.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.roomUseCase.CreateRoom(r.Context(), &room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode((map[string]string{"message": "habitación creada", "roomId": room.ID}))
}

func (h *RoomHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomUseCase.GetAllRooms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	room, err := h.roomUseCase.GetRoom(r.Context(), id)
	if err != nil {
		http.Error(w, "Habitación no encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var room domain.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	room.ID = id
	err := h.roomUseCase.UpdateRoom(r.Context(), &room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode((map[string]string{"message": "habitación actualizada", "roomId": room.ID}))
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.roomUseCase.DeleteRoom(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Habitación eliminada"})
}
