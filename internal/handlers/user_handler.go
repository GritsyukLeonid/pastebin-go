package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

// GetUsersHandler возвращает всех пользователей
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/users [get]
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := repository.GetAllUsers()
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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	user, err := repository.GetUserByID(id)
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
	if err := repository.AddUser(&u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// UpdateUserHandler обновляет пользователя по ID
// @Summary Обновить пользователя
// @Description Обновляет пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body model.User true "Обновленные данные пользователя"
// @Success 200 {object} model.User
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Пользователь не найден"
// @Router /api/user/{id} [put]
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateUser(id, &u); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
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
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
