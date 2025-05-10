package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/GritsyukLeonid/pastebin-go/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/GritsyukLeonid/pastebin-go/internal/httpserver"
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

func main() {
	repository.LoadData()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ch := make(chan model.Storable)

	//go service.GenerateAndSendObjects(ctx, ch)
	go service.StoreFromChannel(ctx, ch)
	go service.LogChanges(ctx)

	// Запуск HTTP-сервера
	httpserver.RegisterHandlers()

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpserver.NewRouter(),
	}

	go func() {
		log.Println("HTTP server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")

	cancel()
	close(ch)

	ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
