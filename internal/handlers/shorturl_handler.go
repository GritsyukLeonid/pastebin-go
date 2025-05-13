package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parts := strings.Split(r.URL.Path, "/")

	if r.Method == http.MethodGet && len(parts) == 4 && parts[3] != "" {
		id := parts[3]
		url, err := repository.GetShortURLByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(url)
		return
	}

	if r.Method == http.MethodPut && len(parts) == 4 && parts[3] != "" {
		id := parts[3]
		var updated model.ShortURL
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		err := repository.UpdateShortURL(id, &updated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodDelete && len(parts) == 4 && parts[3] != "" {
		id := parts[3]
		err := repository.DeleteShortURL(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		urls := repository.GetAllShortURLs()
		json.NewEncoder(w).Encode(urls)
		return
	}

	http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
}
