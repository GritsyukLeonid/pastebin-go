package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/GritsyukLeonid/pastebin-go/internal/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/GritsyukLeonid/pastebin-go/internal/handlers"
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
	"github.com/gorilla/mux"
)

// @title Pastebin API
// @version 1.0
// @description API for managing pastes, users, stats, and short URLs.
// @host localhost:8080
// @BasePath /api

func main() {
	repository.LoadData()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ch := make(chan model.Storable)

	go service.StoreFromChannel(ctx, ch)
	go service.LogChanges(ctx)

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()

	// Paste endpoints
	api.HandleFunc("/paste", handlers.CreatePasteHandler).Methods(http.MethodPost)
	api.HandleFunc("/paste/{id}", handlers.UpdatePasteHandler).Methods(http.MethodPut)
	api.HandleFunc("/paste/{id}", handlers.DeletePasteHandler).Methods(http.MethodDelete)

	// User endpoints
	api.HandleFunc("/user", handlers.GetUsersHandler).Methods(http.MethodGet)
	api.HandleFunc("/user/{id}", handlers.GetUserByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/user", handlers.CreateUserHandler).Methods(http.MethodPost)
	api.HandleFunc("/user/{id}", handlers.UpdateUserHandler).Methods(http.MethodPut)
	api.HandleFunc("/user/{id}", handlers.DeleteUserHandler).Methods(http.MethodDelete)

	// Stats endpoints
	api.HandleFunc("/stats", handlers.GetAllStatsHandler).Methods(http.MethodGet)
	api.HandleFunc("/stat/{id}", handlers.GetStatsByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/stats", handlers.CreateStatsHandler).Methods(http.MethodPost)
	api.HandleFunc("/stat/{id}", handlers.UpdateStatsHandler).Methods(http.MethodPut)
	api.HandleFunc("/stat/{id}", handlers.DeleteStatsHandler).Methods(http.MethodDelete)

	// Short URL endpoints
	api.HandleFunc("/shorturls", handlers.GetAllShortURLsHandler).Methods(http.MethodGet)
	api.HandleFunc("/shorturl/{id}", handlers.GetShortURLByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/shorturl/{id}", handlers.UpdateShortURLHandler).Methods(http.MethodPut)
	api.HandleFunc("/shorturl/{id}", handlers.DeleteShortURLHandler).Methods(http.MethodDelete)

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("HTTP server started on :8080")
		log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	cancel()
	close(ch)

	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}
