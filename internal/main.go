package main

import (
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

func main() {
	ch := make(chan model.Storable)

	go service.GenerateAndSendObjects(ch)
	go repository.StoreFromChannel(ch)
	go repository.LogChanges()

	select {} // блокировка main, чтобы не завершился
}
