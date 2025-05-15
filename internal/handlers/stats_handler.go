package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] != "stats" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(parts) == 4 {
			id := parts[3]
			stat, err := repository.GetStatsByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(stat)
			return
		}
		stats := repository.GetAllStats()
		json.NewEncoder(w).Encode(stats)

	case http.MethodPost:
		var s model.Stats
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := repository.StoreObject(&s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		if len(parts) < 4 {
			http.Error(w, "missing ID in URL", http.StatusBadRequest)
			return
		}
		id := parts[3]
		var updated model.Stats
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := repository.UpdateStats(id, &updated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		if len(parts) < 4 {
			http.Error(w, "missing ID in URL", http.StatusBadRequest)
			return
		}
		id := parts[3]
		err := repository.DeleteStats(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}
