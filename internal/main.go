// @title Pastebin API
// @version 1.0
// @description API для управления пастами, пользователями, статистикой и короткими URL
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/GritsyukLeonid/pastebin-go/internal/docs"
	"github.com/GritsyukLeonid/pastebin-go/internal/handlers"
	"github.com/GritsyukLeonid/pastebin-go/internal/logging"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("ошибка инициализации миграционного драйвера: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("ошибка создания миграции: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("ошибка применения миграции: %v", err)
	}

	log.Println("Миграции успешно применены")
}

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/pastebin?sslmode=disable"
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("не удалось подключиться к PostgreSQL: %v", err)
	}
	defer db.Close()

	runMigrations(db)

	postgresStorage := repository.NewPostgresStorage(db)

	go func() {
		for {
			if err := postgresStorage.DeleteExpiredPastes(); err != nil {
				log.Printf("ошибка при удалении просроченных записей: %v", err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	redisLogger := logging.NewRedisLogger(redisAddr, 10*time.Minute)

	statsService := service.NewStatsService(postgresStorage, redisLogger)

	handlers.InitHandlers(
		service.NewPasteService(postgresStorage, redisLogger, statsService),
		service.NewUserService(postgresStorage, redisLogger),
		statsService,
		service.NewShortURLService(postgresStorage, redisLogger),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/paste", handlers.Paste.CreatePasteHandler).Methods(http.MethodPost)
	api.HandleFunc("/paste/{id}", handlers.Paste.DeletePasteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/paste/{id}", handlers.Paste.GetPasteByIDHandler).Methods(http.MethodGet)

	api.HandleFunc("/user", handlers.User.GetUsersHandler).Methods(http.MethodGet)
	api.HandleFunc("/user/{id}", handlers.User.GetUserByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/user", handlers.User.CreateUserHandler).Methods(http.MethodPost)
	api.HandleFunc("/user/{id}", handlers.User.DeleteUserHandler).Methods(http.MethodDelete)

	api.HandleFunc("/stats", handlers.Stats.GetAllStatsHandler).Methods(http.MethodGet)
	api.HandleFunc("/stat/{id}", handlers.Stats.GetStatsByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/stats", handlers.Stats.CreateStatsHandler).Methods(http.MethodPost)
	api.HandleFunc("/stat/{id}", handlers.Stats.DeleteStatsHandler).Methods(http.MethodDelete)

	api.HandleFunc("/shorturls", handlers.ShortURL.GetAllShortURLsHandler).Methods(http.MethodGet)
	api.HandleFunc("/shorturl/{id}", handlers.ShortURL.GetShortURLByIDHandler).Methods(http.MethodGet)
	api.HandleFunc("/shorturl/{id}", handlers.ShortURL.DeleteShortURLHandler).Methods(http.MethodDelete)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("🚀 HTTP server started on :8080")
		log.Println("📚 Swagger UI: http://localhost:8080/swagger/index.html")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Shutting down server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}
}
