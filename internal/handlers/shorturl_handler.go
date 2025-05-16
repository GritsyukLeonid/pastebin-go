package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

// GetShortURLByIDHandler получает короткий URL по ID
// @Summary Получить короткий URL
// @Description Возвращает короткий URL по ID
// @Tags shorturls
// @Produce json
// @Param id path string true "ID ShortURL"
// @Success 200 {object} model.ShortURL
// @Failure 404 {string} string "ShortURL не найден"
// @Router /api/shorturl/{id} [get]
func GetShortURLByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/shorturl/")
	url, err := repository.GetShortURLByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

// UpdateShortURLHandler обновляет короткий URL по ID
// @Summary Обновить короткий URL
// @Description Обновляет короткий URL по ID
// @Tags shorturls
// @Accept json
// @Produce json
// @Param id path string true "ID ShortURL"
// @Param shorturl body model.ShortURL true "Данные ShortURL"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Некорректный ввод"
// @Failure 404 {string} string "ShortURL не найден"
// @Router /api/shorturl/{id} [put]
func UpdateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/shorturl/")
	var updated model.ShortURL
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := repository.UpdateShortURL(id, &updated); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteShortURLHandler удаляет короткий URL по ID
// @Summary Удалить короткий URL
// @Description Удаляет короткий URL по ID
// @Tags shorturls
// @Param id path string true "ID ShortURL"
// @Success 200 {string} string "ShortURL удалён"
// @Failure 404 {string} string "ShortURL не найден"
// @Router /api/shorturl/{id} [delete]
func DeleteShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/shorturl/")
	if err := repository.DeleteShortURL(id); err != nil {
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
func GetAllShortURLsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	urls := repository.GetAllShortURLs()
	json.NewEncoder(w).Encode(urls)
}
