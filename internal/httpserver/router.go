package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/pastes", handlePastes).Methods("POST")
	r.HandleFunc("/pastes/{id}", handlePasteByID).Methods("GET")
	r.HandleFunc("/pastes/{id}", handlePasteByID).Methods("PUT")
	r.HandleFunc("/pastes/{id}", handlePasteByID).Methods("DELETE")

	// ShortURL handlers
	r.HandleFunc("/shorturls", handleURLs).Methods("POST")
	r.HandleFunc("/shorturls/{id}", handleURLByID).Methods("GET")
	r.HandleFunc("/shorturls/{id}", handleURLByID).Methods("PUT")
	r.HandleFunc("/shorturls/{id}", handleURLByID).Methods("DELETE")

	// Stats handlers
	r.HandleFunc("/stats", handleStats).Methods("POST")
	r.HandleFunc("/stats/{id}", handleStatByID).Methods("GET")
	r.HandleFunc("/stats/{id}", handleStatByID).Methods("PUT")
	r.HandleFunc("/stats/{id}", handleStatByID).Methods("DELETE")

	// User handlers
	r.HandleFunc("/users", handleUsers).Methods("POST")
	r.HandleFunc("/users/{id}", handleUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", handleUserByID).Methods("PUT")
	r.HandleFunc("/users/{id}", handleUserByID).Methods("DELETE")

	return r
}
