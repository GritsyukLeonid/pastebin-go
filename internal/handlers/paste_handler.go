package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

// CreatePasteHandler создаёт новую пасту
// @Summary Создать новую запись
// @Description Создает новую запись paste
// @Tags pastes
// @Accept json
// @Produce json
// @Param paste body model.Paste true "Paste объект"
// @Success 201 {object} model.Paste
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 500 {string} string "Ошибка на сервере"
// @Router /api/paste [post]
func CreatePasteHandler(w http.ResponseWriter, r *http.Request) {
	var p model.Paste
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.ExpiresAt.IsZero() {
		p.ExpiresAt = p.CreatedAt.Add(7 * 24 * time.Hour)
	}

	if err := repository.StoreObject(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// UpdatePasteHandler обновляет существующую пасту по ID
// @Summary Обновить запись
// @Description Обновляет paste по ID
// @Tags pastes
// @Accept json
// @Produce json
// @Param id path string true "ID пасты"
// @Param paste body model.Paste true "Новые данные пасты"
// @Success 200 {object} model.Paste
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 404 {string} string "Paste не найден"
// @Router /api/paste/{id} [put]
func UpdatePasteHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	var p model.Paste
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := repository.UpdatePaste(id, &p); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

// DeletePasteHandler удаляет пасту по ID
// @Summary Удалить запись
// @Description Удаляет paste по ID
// @Tags pastes
// @Param id path string true "ID пасты"
// @Success 204 {string} string "Паста удалена"
// @Failure 404 {string} string "Paste не найден"
// @Router /api/paste/{id} [delete]
func DeletePasteHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	if err := repository.DeletePaste(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
