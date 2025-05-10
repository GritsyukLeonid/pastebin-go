package service

import (
	"context"
	"log"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

func StoreFromChannel(ctx context.Context, ch <-chan model.Storable) {
	for {
		select {
		case <-ctx.Done():
			return
		case obj, ok := <-ch:
			if !ok {
				return
			}
			if err := repository.StoreObject(obj); err != nil {
				log.Println("Error storing object:", err)
			}
		}
	}
}

func LogChanges(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			newObjects := repository.DetectNewObjects()
			for typ, entries := range newObjects {
				for _, entry := range entries {
					log.Printf("%s: %s\n", typ, entry)
				}
			}
		}
	}
}
