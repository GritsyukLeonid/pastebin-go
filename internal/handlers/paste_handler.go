package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

func PasteHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if r.Method == http.MethodPost {
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
		return
	}
	if r.Method == http.MethodPost {
		var p model.Paste
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if p.Content == "" {
			http.Error(w, "content required", http.StatusBadRequest)
			return
		}
		repository.StoreObject(&p)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
		return
	}
	if r.Method == http.MethodPut {
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
		return
	}
	if r.Method == http.MethodDelete {
		id := parts[len(parts)-1]
		if err := repository.DeletePaste(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
