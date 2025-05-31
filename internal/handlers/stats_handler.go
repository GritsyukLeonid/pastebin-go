package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

// GetAllStatsHandler возвращает все записи статистики
// @Summary Получить все статистики
// @Description Возвращает список всех статистик
// @Tags stats
// @Produce json
// @Success 200 {array} model.Stats
// @Router /api/stats [get]
func GetAllStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := statsService.ListStats(r.Context())
	if err != nil {
		http.Error(w, "Ошибка при получении статистики", http.StatusInternalServerError)
		return
	}
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
	stat, err := statsService.GetStatsByID(r.Context(), id)
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
// @Success 201 {object} model.Stats
// @Failure 400 {string} string "Некорректный ввод"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/stats [post]
func CreateStatsHandler(w http.ResponseWriter, r *http.Request) {
	var s model.Stats
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := statsService.CreateStats(r.Context(), s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
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
	if err := statsService.DeleteStats(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
