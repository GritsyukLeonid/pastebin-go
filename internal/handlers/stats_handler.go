package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

type StatsHandler struct {
	service service.StatsService
}

func NewStatsHandler(s service.StatsService) *StatsHandler {
	return &StatsHandler{service: s}
}

type CreateStatsRequest struct{}

// GetAllStatsHandler возвращает все записи статистики
// @Summary Получить все статистики
// @Description Возвращает список всех статистик
// @Tags stats
// @Produce json
// @Success 200 {array} model.Stats
// @Router /api/stats [get]
func (h *StatsHandler) GetAllStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.ListStats(r.Context())
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
func (h *StatsHandler) GetStatsByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	stat, err := h.service.GetStatsByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stat)
}

// @Summary Создать статистику
// @Description Создает новую запись статистики (ID и views генерируются на сервере)
// @Tags stats
// @Accept json
// @Produce json
// @Param stats body handlers.CreateStatsRequest true "Пустой объект запроса"
// @Success 201 {object} model.Stats
// @Failure 400 {string} string "Некорректный ввод"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/stats [post]
func (h *StatsHandler) CreateStatsHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateStats(r.Context(), model.Stats{})
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
func (h *StatsHandler) DeleteStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteStats(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetPopularPastesHandler возвращает N самых популярных записей по количеству просмотров
// @Summary Получить популярные пасты
// @Description Возвращает топ паст по количеству просмотров
// @Tags stats
// @Produce json
// @Param limit query int false "Максимальное количество записей (по умолчанию 5)"
// @Success 200 {array} model.Stats
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/paste/popular [get]
func (h *StatsHandler) GetPopularPastesHandler(w http.ResponseWriter, r *http.Request) {
	limit := 5 // значение по умолчанию

	// читаем limit из query-параметра
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	stats, err := h.service.ListTopStats(r.Context(), limit)
	if err != nil {
		http.Error(w, "Ошибка при получении популярных паст", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
