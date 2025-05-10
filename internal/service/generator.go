package service

import (
	"context"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

func GenerateAndSendObjects(ctx context.Context, ch chan<- model.Storable) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ch <- model.NewPaste("example", time.Minute)
			ch <- model.NewUser("Leonid")
			ch <- model.NewStats("xyz123")
			ch <- model.NewShortURL("https://example.com")
		}
	}
}
