package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if r.Method == http.MethodGet && parts[len(parts)-1] != "" && parts[len(parts)-2] == "user" {
		idStr := parts[len(parts)-1]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		u, err := repository.GetUserByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
		return
	}

	if r.Method == http.MethodGet {
		us := repository.GetAllUsers()
		json.NewEncoder(w).Encode(us)
		return
	}

	if r.Method == http.MethodPost {
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if u.Username == "" {
			http.Error(w, "username required", http.StatusBadRequest)
			return
		}
		if err := repository.AddUser(&u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
		return
	}

	if r.Method == http.MethodPut {
		idStr := parts[len(parts)-1]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if err := repository.UpdateUser(id, &u); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
		return
	}

	if r.Method == http.MethodDelete {
		idStr := parts[len(parts)-1]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := repository.DeleteUser(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
