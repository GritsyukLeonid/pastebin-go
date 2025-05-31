package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

// GetUsersHandler возвращает всех пользователей
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/users [get]
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := userService.ListUsers(r.Context())
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
// @Param id path int true "ID пользователя"
// @Success 200 {object} model.User
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Пользователь не найден"
// @Router /api/user/{id} [get]
func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	id := idStr // ID в Mongo и сервисе — string, а не int64

	user, err := userService.GetUserByID(r.Context(), id)
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
// @Param user body model.User true "Данные пользователя"
// @Success 201 {object} model.User
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/user [post]
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if u.Username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}
	created, err := userService.CreateUser(r.Context(), u)
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
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	id := idStr // ID как string, чтобы не преобразовывать

	if err := userService.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
