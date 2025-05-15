package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GritsyukLeonid/pastebin-go/internal/handlers"
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

	go service.StoreFromChannel(ctx, ch)
	go service.LogChanges(ctx)

	http.HandleFunc("/api/pastes", handlers.PasteHandler)
	http.HandleFunc("/api/paste/", handlers.PasteHandler)
	http.HandleFunc("/api/users", handlers.UserHandler)
	http.HandleFunc("/api/user/", handlers.UserHandler)
	http.HandleFunc("/api/stats", handlers.StatsHandler)
	http.HandleFunc("/api/stat/", handlers.StatsHandler)
	http.HandleFunc("/api/shorturls", handlers.ShortURLHandler)
	http.HandleFunc("/api/shorturl/", handlers.ShortURLHandler)

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("HTTP server started on :8080")
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
