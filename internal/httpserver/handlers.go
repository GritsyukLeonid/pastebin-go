package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

var (
	pasteMu = &sync.Mutex{}
	userMu  = &sync.Mutex{}
	statsMu = &sync.Mutex{}
	urlMu   = &sync.Mutex{}
)

func RegisterHandlers() {
	http.HandleFunc("/api/paste", handlePastes)
	http.HandleFunc("/api/paste/", handlePasteByID)
	http.HandleFunc("/api/user", handleUsers)
	http.HandleFunc("/api/user/", handleUserByID)
	http.HandleFunc("/api/stats", handleStats)
	http.HandleFunc("/api/stats/", handleStatByID)
	http.HandleFunc("/api/url", handleURLs)
	http.HandleFunc("/api/url/", handleURLByID)
}

func handlePastes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		pasteMu.Lock()
		defer pasteMu.Unlock()
		json.NewEncoder(w).Encode(repository.Pastes)
	case http.MethodPost:
		var p model.Paste
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		pasteMu.Lock()
		defer pasteMu.Unlock()
		repository.Pastes = append(repository.Pastes, &p)
		repository.SaveAll()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePasteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/paste/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	pasteMu.Lock()
	defer pasteMu.Unlock()

	if id < 0 || id >= len(repository.Pastes) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(repository.Pastes[id])
	case http.MethodPut:
		var p model.Paste
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		repository.Pastes[id] = &p
		repository.SaveAll()
		json.NewEncoder(w).Encode(p)
	case http.MethodDelete:
		repository.Pastes = append(repository.Pastes[:id], repository.Pastes[id+1:]...)
		repository.SaveAll()
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		userMu.Lock()
		defer userMu.Unlock()
		json.NewEncoder(w).Encode(repository.Users)
	case http.MethodPost:
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		userMu.Lock()
		defer userMu.Unlock()
		repository.Users = append(repository.Users, &u)
		repository.SaveAll()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	userMu.Lock()
	defer userMu.Unlock()

	if id < 0 || id >= len(repository.Users) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(repository.Users[id])
	case http.MethodPut:
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		repository.Users[id] = &u
		repository.SaveAll()
		json.NewEncoder(w).Encode(u)
	case http.MethodDelete:
		repository.Users = append(repository.Users[:id], repository.Users[id+1:]...)
		repository.SaveAll()
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		statsMu.Lock()
		defer statsMu.Unlock()
		json.NewEncoder(w).Encode(repository.StatsSet)
	case http.MethodPost:
		var s model.Stats
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		statsMu.Lock()
		defer statsMu.Unlock()
		repository.StatsSet = append(repository.StatsSet, &s)
		repository.SaveAll()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(s)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleStatByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/stats/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	statsMu.Lock()
	defer statsMu.Unlock()

	if id < 0 || id >= len(repository.StatsSet) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(repository.StatsSet[id])
	case http.MethodPut:
		var s model.Stats
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		repository.StatsSet[id] = &s
		repository.SaveAll()
		json.NewEncoder(w).Encode(s)
	case http.MethodDelete:
		repository.StatsSet = append(repository.StatsSet[:id], repository.StatsSet[id+1:]...)
		repository.SaveAll()
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		urlMu.Lock()
		defer urlMu.Unlock()
		json.NewEncoder(w).Encode(repository.URLs)
	case http.MethodPost:
		var u model.ShortURL
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		urlMu.Lock()
		defer urlMu.Unlock()
		repository.URLs = append(repository.URLs, &u)
		repository.SaveAll()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleURLByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/url/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	urlMu.Lock()
	defer urlMu.Unlock()

	if id < 0 || id >= len(repository.URLs) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(repository.URLs[id])
	case http.MethodPut:
		var u model.ShortURL
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		repository.URLs[id] = &u
		repository.SaveAll()
		json.NewEncoder(w).Encode(u)
	case http.MethodDelete:
		repository.URLs = append(repository.URLs[:id], repository.URLs[id+1:]...)
		repository.SaveAll()
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
