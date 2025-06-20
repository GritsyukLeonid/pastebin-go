package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

type StatsHandler struct {
	statsService service.StatsService
	pasteService service.PasteService
}

func NewStatsHandler(statsSvc service.StatsService, pasteSvc service.PasteService) *StatsHandler {
	return &StatsHandler{
		statsService: statsSvc,
		pasteService: pasteSvc,
	}
}

type CreateStatsRequest struct{}

// @Summary Получить всю статистику
// @Description Возвращает список всех записей статистики просмотров
// @Tags stats
// @Produce json
// @Success 200 {array} model.Stats
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/stats [get]
func (h *StatsHandler) GetAllStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsService.ListStats(r.Context())
	if err != nil {
		http.Error(w, "Ошибка при получении статистики", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// @Summary Получить статистику по ID
// @Description Возвращает статистику просмотров для указанного ID пасты
// @Tags stats
// @Produce json
// @Param id path string true "ID статистики (равен ID пасты)"
// @Success 200 {object} model.Stats
// @Failure 400 {string} string "ID отсутствует"
// @Failure 404 {string} string "Статистика не найдена"
// @Router /api/stat/{id} [get]
func (h *StatsHandler) GetStatsByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	stat, err := h.statsService.GetStatsByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stat)
}

// @Summary Создать новую запись статистики
// @Description Создаёт пустую запись статистики. Используется редко, т.к. обычно статистика создаётся автоматически.
// @Tags stats
// @Accept json
// @Produce json
// @Param stats body handlers.CreateStatsRequest true "Пустой объект запроса"
// @Success 201 {object} model.Stats
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/stats [post]
func (h *StatsHandler) CreateStatsHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.statsService.CreateStats(r.Context(), model.Stats{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// @Summary Удалить статистику
// @Description Удаляет статистику просмотров по ID пасты
// @Tags stats
// @Param id path string true "ID статистики (равен ID пасты)"
// @Success 204 {string} string "Статистика удалена"
// @Failure 400 {string} string "ID отсутствует"
// @Failure 404 {string} string "Статистика не найдена"
// @Router /api/stat/{id} [delete]
func (h *StatsHandler) DeleteStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := h.statsService.DeleteStats(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Получить популярные пасты
// @Description Возвращает список самых просматриваемых паст (по убыванию просмотров)
// @Tags stats
// @Produce json
// @Param limit query int false "Максимальное количество записей (по умолчанию 5)"
// @Success 200 {array} model.Paste
// @Failure 500 {string} string "Ошибка при получении статистики"
// @Failure 404 {string} string "Популярные пасты не найдены"
// @Router /api/paste/popular [get]
func (h *StatsHandler) GetPopularPastesHandler(w http.ResponseWriter, r *http.Request) {
	limit := 5
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	stats, err := h.statsService.ListTopStats(r.Context(), limit)
	if err != nil {
		http.Error(w, "Ошибка при получении популярных паст", http.StatusInternalServerError)
		return
	}

	var popularPastes []model.Paste
	for _, stat := range stats {
		paste, err := h.pasteService.GetPasteByID(r.Context(), stat.ID)
		if err != nil {
			log.Printf("❌ Не найдена паста с ID %s: %v", stat.ID, err)
			continue
		}
		paste.Views = stat.Views
		popularPastes = append(popularPastes, paste)
	}

	if len(popularPastes) == 0 {
		http.Error(w, "No popular pastes found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(popularPastes)
}
