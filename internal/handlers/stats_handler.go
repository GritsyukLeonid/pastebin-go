package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

// GetAllStatsHandler возвращает все записи статистики
// @Summary Получить все статистики
// @Description Возвращает список всех статистик
// @Tags stats
// @Produce json
// @Success 200 {array} model.Stats
// @Router /api/stats [get]
func GetAllStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := repository.GetAllStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetStatsByIDHandler возвращает статистику по ID
// @Summary Получить статистику по ID
// @Description Возвращает статистику по заданному ID
// @Tags stats
// @Produce json
// @Param id path string true "ID статистики"
// @Success 200 {object} model.Stats
// @Failure 404 {string} string "Статистика не найдена"
// @Router /api/stat/{id} [get]
func GetStatsByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/stat/")
	stat, err := repository.GetStatsByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stat)
}

// CreateStatsHandler создает новую статистику
// @Summary Создать статистику
// @Description Создает новую запись статистики
// @Tags stats
// @Accept json
// @Produce json
// @Param stats body model.Stats true "Статистика"
// @Success 201 {string} string "Статистика создана"
// @Failure 400 {string} string "Некорректный ввод"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/stats [post]
func CreateStatsHandler(w http.ResponseWriter, r *http.Request) {
	var s model.Stats
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := repository.StoreObject(&s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UpdateStatsHandler обновляет статистику по ID
// @Summary Обновить статистику
// @Description Обновляет запись статистики по ID
// @Tags stats
// @Accept json
// @Param id path string true "ID статистики"
// @Param stats body model.Stats true "Обновленная статистика"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Некорректный ввод"
// @Failure 404 {string} string "Статистика не найдена"
// @Router /api/stat/{id} [put]
func UpdateStatsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/stat/")
	var updated model.Stats
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := repository.UpdateStats(id, &updated); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteStatsHandler удаляет статистику по ID
// @Summary Удалить статистику
// @Description Удаляет запись статистики по ID
// @Tags stats
// @Param id path string true "ID статистики"
// @Success 204 {string} string "Статистика удалена"
// @Failure 404 {string} string "Статистика не найдена"
// @Router /api/stat/{id} [delete]
func DeleteStatsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/stat/")
	if err := repository.DeleteStats(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
