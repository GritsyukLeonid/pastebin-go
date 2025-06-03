package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

type ShortURLHandler struct {
	service service.ShortURLService
}

func NewShortURLHandler(s service.ShortURLService) *ShortURLHandler {
	return &ShortURLHandler{service: s}
}

// GetShortURLByIDHandler получает короткий URL по ID
// @Summary Получить короткий URL
// @Description Возвращает короткий URL по ID
// @Tags shorturls
// @Produce json
// @Param id path string true "ID ShortURL"
// @Success 200 {object} model.ShortURL
// @Failure 404 {string} string "ShortURL не найден"
// @Router /api/shorturl/{id} [get]
func (h *ShortURLHandler) GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/shorturl/")
	url, err := h.service.GetShortURLByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

// DeleteShortURLHandler удаляет короткий URL по ID
// @Summary Удалить короткий URL
// @Description Удаляет короткий URL по ID
// @Tags shorturls
// @Param id path string true "ID ShortURL"
// @Success 200 {string} string "ShortURL удалён"
// @Failure 404 {string} string "ShortURL не найден"
// @Router /api/shorturl/{id} [delete]
func (h *ShortURLHandler) DeleteShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/shorturl/")
	if err := h.service.DeleteShortURL(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetAllShortURLsHandler возвращает все короткие URL
// @Summary Получить все короткие URL
// @Description Возвращает список всех коротких URL
// @Tags shorturls
// @Produce json
// @Success 200 {array} model.ShortURL
// @Router /api/shorturls [get]
func (h *ShortURLHandler) GetAllShortURLsHandler(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.ListShortURLs(r.Context())
	if err != nil {
		http.Error(w, "Ошибка при получении URL", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}
