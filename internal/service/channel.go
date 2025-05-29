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
		case obj := <-ch:
			if err := repository.StoreObject(obj); err != nil {
				log.Printf("failed to store object: %v", err)
			}
		}
	}
}

func LogChanges(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			newObjects := repository.DetectNewObjects()
			for k, v := range newObjects {
				for _, obj := range v {
					log.Printf("[New %s] %s", k, obj)
				}
			}
		}
	}
}
