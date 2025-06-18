package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

type PasteHandler struct {
	service service.PasteService
}

type CreatePasteRequest struct {
	Content   string    `json:"content"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewPasteHandler(s service.PasteService) *PasteHandler {
	return &PasteHandler{service: s}
}

// CreatePasteHandler создаёт новую пасту
// @Summary Создать новую запись
// @Description Создает новую запись paste
// @Tags pastes
// @Accept json
// @Produce json
// @Param paste body handlers.CreatePasteRequest true "Paste объект"
// @Success 201 {object} model.Paste
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 500 {string} string "Ошибка на сервере"
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

	type PasteCreateResponse struct {
		Hash     string `json:"hash"`
		ShortURL string `json:"short_url"`
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PasteCreateResponse{
		Hash:     created.Hash,
		ShortURL: created.Hash[:6],
	})

}

// DeletePasteHandler удаляет пасту по ID
// @Summary Удалить запись
// @Description Удаляет paste по ID
// @Tags pastes
// @Param id path string true "ID пасты"
// @Success 204 "Паста удалена"
// @Failure 404 {string} string "Paste не найден"
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

// GetPasteByIDHandler возвращает пасту по ID
// @Summary Получить запись по ID
// @Description Возвращает paste по ID
// @Tags pastes
// @Produce json
// @Param id path string true "ID пасты"
// @Success 200 {object} model.Paste
// @Failure 404 {string} string "Paste не найден"
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paste)
}

// GetPasteByHashHandler возвращает пасту по хэшу
// @Summary Получить запись по хэшу
// @Description Возвращает paste по hash
// @Tags pastes
// @Produce json
// @Param hash path string true "Hash пасты"
// @Success 200 {object} model.Paste
// @Failure 404 {string} string "Paste не найден"
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

	w.Header().Set("Content-Type", "application/json")
	type PasteContentResponse struct {
		Content string `json:"content"`
	}

	json.NewEncoder(w).Encode(PasteContentResponse{
		Content: paste.Content,
	})
}
