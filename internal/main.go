package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	go service.GenerateAndSendObjects(ctx, ch)
	go service.StoreFromChannel(ctx, ch)
	go service.LogChanges(ctx)

	<-stop
	cancel()
	close(ch)
}
