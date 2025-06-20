package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

type PasteHandler struct {
	service      service.PasteService
	statsService service.StatsService
}

type CreatePasteRequest struct {
	Content   string    `json:"content"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type PasteCreateResponse struct {
	ID       string `json:"id"`
	Hash     string `json:"hash"`
	ShortURL string `json:"short_url"`
}

func NewPasteHandler(pasteSvc service.PasteService, statsSvc service.StatsService) *PasteHandler {
	return &PasteHandler{
		service:      pasteSvc,
		statsService: statsSvc,
	}
}

// @Summary Создать новую пасту
// @Description Создает новую пасту с указанным содержимым и временем истечения. Возвращает ID, hash и короткий URL.
// @Tags pastes
// @Accept json
// @Produce json
// @Param paste body handlers.CreatePasteRequest true "Данные пасты"
// @Success 201 {object} PasteCreateResponse
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/paste [post]
func (h *PasteHandler) CreatePasteHandler(w http.ResponseWriter, r *http.Request) {

	var req CreatePasteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "content required", http.StatusBadRequest)
		return
	}

	if req.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Expiration date must be in the future", http.StatusBadRequest)
		return
	}

	paste := model.Paste{
		Content:   req.Content,
		ExpiresAt: req.ExpiresAt,
	}

	created, err := h.service.CreatePaste(r.Context(), paste)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PasteCreateResponse{
		ID:       created.ID,
		Hash:     created.Hash,
		ShortURL: fmt.Sprintf("http://localhost:8080/s/%s", created.Hash[:6]),
	})

}

// @Summary Удалить пасту по ID
// @Description Удаляет существующую пасту по её уникальному ID
// @Tags pastes
// @Param id path string true "ID пасты"
// @Success 204 {string} string "Паста удалена"
// @Failure 400 {string} string "Некорректный ID"
// @Failure 404 {string} string "Паста не найдена"
// @Router /api/paste/{id} [delete]
func (h *PasteHandler) DeletePasteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePaste(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Получить пасту по ID
// @Description Возвращает полную информацию о пасте по её ID. Также увеличивает счётчик просмотров.
// @Tags pastes
// @Produce json
// @Param id path string true "ID пасты"
// @Success 200 {object} model.Paste
// @Failure 400 {string} string "Некорректный ID"
// @Failure 404 {string} string "Паста не найдена"
// @Router /api/paste/{id} [get]
func (h *PasteHandler) GetPasteByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	paste, err := h.service.GetPasteByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	_ = h.statsService.IncrementViews(r.Context(), paste.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paste)
}

// @Summary Получить пасту по hash
// @Description Возвращает пасту по уникальному hash. Также увеличивает счётчик просмотров.
// @Tags pastes
// @Produce json
// @Param hash path string true "Hash пасты"
// @Success 200 {object} model.Paste
// @Failure 400 {string} string "Некорректный hash"
// @Failure 404 {string} string "Паста не найдена"
// @Router /api/paste/hash/{hash} [get]
func (h *PasteHandler) GetPasteByHashHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash, ok := vars["hash"]
	if !ok {
		http.Error(w, "missing hash", http.StatusBadRequest)
		return
	}

	paste, err := h.service.GetPasteByHash(r.Context(), hash)
	if err != nil {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	_ = h.statsService.IncrementViews(r.Context(), paste.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paste)
}
