package service

import (
	"context"
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
			repository.StoreObject(obj)
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
			repository.CheckAndLogChanges()
		}
	}
}
