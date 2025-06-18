package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

type ShortURLHandler struct {
	service      service.ShortURLService
	pasteService service.PasteService
}

func NewShortURLHandler(s service.ShortURLService, ps service.PasteService) *ShortURLHandler {
	return &ShortURLHandler{
		service:      s,
		pasteService: ps,
	}
}

// CreateShortURLHandler создаёт короткий URL, используя укороченный hash
// @Summary Создать короткий URL
// @Description Использует hash пасты как основу для короткой ссылки
// @Tags shorturls
// @Param hash path string true "Hash пасты"
// @Success 201 {object} model.ShortURL
// @Failure 400 {string} string "Хэш слишком короткий"
// @Failure 500 {string} string "Ошибка сервиса"
// @Router /api/shorturl/{hash} [post]
func (h *ShortURLHandler) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash, ok := vars["hash"]
	if !ok {
		http.Error(w, "missing hash", http.StatusBadRequest)
		return
	}

	if len(hash) < 6 {
		http.Error(w, "Хэш слишком короткий для сокращения", http.StatusBadRequest)
		return
	}

	short := model.NewShortURL(hash, hash[:6])

	created, err := h.service.CreateShortURL(r.Context(), *short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// ResolveShortURLHandler возвращает пасту по короткому коду
// @Summary Получить пасту по короткой ссылке
// @Description Возвращает пасту, связанную с коротким URL
// @Tags shorturls
// @Param code path string true "Короткий код"
// @Success 200 {object} model.Paste
// @Failure 404 {string} string "ShortURL или паста не найдена"
// @Router /s/{code} [get]
func (h *ShortURLHandler) ResolveShortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	if !ok {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	short, err := h.service.GetShortURLByID(r.Context(), code)
	if err != nil {
		http.Error(w, "ShortURL не найден", http.StatusNotFound)
		return
	}

	paste, err := h.pasteService.GetPasteByHash(r.Context(), short.Original)
	if err != nil {
		http.Error(w, "Паста не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paste)
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
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

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
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

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
