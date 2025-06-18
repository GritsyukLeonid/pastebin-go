package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

type CreateUserRequest struct {
	Username string `json:"username"`
}

// GetUsersHandler возвращает всех пользователей
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/users [get]
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Ошибка при получении пользователей", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByIDHandler возвращает пользователя по ID
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по заданному ID
// @Tags users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} model.User
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Пользователь не найден"
// @Router /api/user/{id} [get]
func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUserHandler создает нового пользователя
// @Summary Создать нового пользователя
// @Description Создает нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body handlers.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} model.User
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/user [post]
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if u.Username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}
	created, err := h.service.CreateUser(r.Context(), model.User{
		Username: u.Username,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// DeleteUserHandler удаляет пользователя по ID
// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Tags users
// @Param id path int true "ID пользователя"
// @Success 204 {string} string "Пользователь удален"
// @Failure 404 {string} string "Пользователь не найден"
// @Router /api/user/{id} [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
