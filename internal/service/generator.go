package service

import (
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

func GenerateAndStoreObjects() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		p := model.NewPaste("example", time.Minute)
		u := model.NewUser("Leonid")
		s := model.NewStats("xyz123")
		url := model.NewShortURL("https://example.com")

		repository.Store(p)
		repository.Store(u)
		repository.Store(s)
		repository.Store(url)
	}
}
