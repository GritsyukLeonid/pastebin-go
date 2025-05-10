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

// @Summary Get all pastes or create a new paste
// @Tags Paste
// @Accept json
// @Produce json
// @Success 200 {array} model.Paste
// @Success 201 {object} model.Paste
// @Router /api/paste [get]
// @Router /api/paste [post]
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

// @Summary Get, update or delete a paste by ID
// @Tags Paste
// @Accept json
// @Produce json
// @Param id path int true "Paste ID"
// @Success 200 {object} model.Paste
// @Success 204 {string} string "No Content"
// @Router /api/paste/{id} [get]
// @Router /api/paste/{id} [put]
// @Router /api/paste/{id} [delete]
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

// @Summary Get all users or create a new user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {array} model.User
// @Success 201 {object} model.User
// @Router /api/user [get]
// @Router /api/user [post]
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

// @Summary Get, update or delete a user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Success 204 {string} string "No Content"
// @Router /api/user/{id} [get]
// @Router /api/user/{id} [put]
// @Router /api/user/{id} [delete]
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

// @Summary Get all stats or create a new stat
// @Tags Stats
// @Accept json
// @Produce json
// @Success 200 {array} model.Stats
// @Success 201 {object} model.Stats
// @Router /api/stats [get]
// @Router /api/stats [post]
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

// @Summary Get, update or delete a stat by ID
// @Tags Stats
// @Accept json
// @Produce json
// @Param id path int true "Stat ID"
// @Success 200 {object} model.Stats
// @Success 204 {string} string "No Content"
// @Router /api/stats/{id} [get]
// @Router /api/stats/{id} [put]
// @Router /api/stats/{id} [delete]
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

// @Summary Get all short URLs or create a new short URL
// @Tags ShortURL
// @Accept json
// @Produce json
// @Success 200 {array} model.ShortURL
// @Success 201 {object} model.ShortURL
// @Router /api/url [get]
// @Router /api/url [post]
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

// @Summary Get, update or delete a short URL by ID
// @Tags ShortURL
// @Accept json
// @Produce json
// @Param id path int true "Short URL ID"
// @Success 200 {object} model.ShortURL
// @Success 204 {string} string "No Content"
// @Router /api/url/{id} [get]
// @Router /api/url/{id} [put]
// @Router /api/url/{id} [delete]
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
