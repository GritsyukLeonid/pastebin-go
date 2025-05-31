package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/GritsyukLeonid/pastebin-go/internal/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/GritsyukLeonid/pastebin-go/internal/handlers"
	"github.com/GritsyukLeonid/pastebin-go/internal/logging"
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
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	mongoStorage, err := repository.NewMongoStorage(mongoURI, "pastebin")
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}
	redisLogger := logging.NewRedisLogger(redisAddr, 10*time.Minute)

	// Прокидываем зависимости в хендлеры
	handlers.SetServices(
		service.NewPasteService(mongoStorage, redisLogger),
		service.NewUserService(mongoStorage, redisLogger),
		service.NewStatsService(mongoStorage, redisLogger),
		service.NewShortURLService(mongoStorage, redisLogger),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// Paste endpoints
	api.HandleFunc("/paste", handlers.CreatePasteHandler).Methods(http.MethodPost)
	api.HandleFunc("/paste/{id}", handlers.DeletePasteHandler).Methods(http.MethodDelete)

	// User endpoints
	api.HandleFunc("/user", handlers.GetUsersHandler).Methods(http.MethodGet)
	api.HandleFunc("/user/{id}", handlers.GetUserByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/user", handlers.CreateUserHandler).Methods(http.MethodPost)
	api.HandleFunc("/user/{id}", handlers.DeleteUserHandler).Methods(http.MethodDelete)

	// Stats endpoints
	api.HandleFunc("/stats", handlers.GetAllStatsHandler).Methods(http.MethodGet)
	api.HandleFunc("/stat/{id}", handlers.GetStatsByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/stats", handlers.CreateStatsHandler).Methods(http.MethodPost)
	api.HandleFunc("/stat/{id}", handlers.DeleteStatsHandler).Methods(http.MethodDelete)

	// Short URL endpoints
	api.HandleFunc("/shorturls", handlers.GetAllShortURLsHandler).Methods(http.MethodGet)
	api.HandleFunc("/shorturl/{id}", handlers.GetShortURLByIDHandler).Methods(http.MethodGet)
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
	log.Println("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
