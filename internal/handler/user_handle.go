package handler

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Failed to decode the user"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.usecase.RegisterUser(context.Background(), &u); err != nil {
		http.Error(w, "Failed to register the user"+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetAllUser(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch user"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id not found", http.StatusBadRequest)
		return
	}
	user, err := h.usecase.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get the user by the req id"+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	u := &models.User{}
	json.NewDecoder(r.Body).Decode(u)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id not found", http.StatusBadRequest)
		return
	}

	u.ID = id
	h.usecase.UpdateUser(r.Context(), u)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Failed to extract id from ResponseWritere url", http.StatusBadRequest)
		return
	}

	err = h.usecase.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete ResponseWritere user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User Deleted Successfully")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Failed to decode the user", http.StatusBadRequest)
		return
	}

	token, err := h.usecase.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, "Failed to get the token "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
