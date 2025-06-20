package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/logging"
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type pasteService struct {
	storage      repository.StorageInterface
	logger       logging.Logger
	statsService StatsService
	shortService ShortURLService
}

func NewPasteService(storage repository.StorageInterface, logger logging.Logger, stats StatsService, short ShortURLService) PasteService {
	return &pasteService{
		storage:      storage,
		logger:       logger,
		statsService: stats,
		shortService: short,
	}
}

func (s *pasteService) CreatePaste(ctx context.Context, p model.Paste) (model.Paste, error) {
	now := time.Now()
	p.ID = fmt.Sprintf("%d", now.UnixNano())
	p.CreatedAt = now

	hash := sha1.New()
	hash.Write([]byte(p.Content + now.String()))
	p.Hash = fmt.Sprintf("%x", hash.Sum(nil))[:10]

	if err := s.storage.SavePaste(p); err != nil {
		return model.Paste{}, err
	}

	if len(p.Hash) >= 6 {
		_, _ = s.shortService.CreateShortURL(ctx, model.ShortURL{
			ID:       p.Hash[:6],
			Original: p.Hash,
		})
	}

	if p.ExpiresAt.Before(now) {
		return model.Paste{}, errors.New("expiration must be in the future")
	}

	_ = s.logger.LogChange("paste", p.ID, "created")
	return p, nil
}

func (s *pasteService) GetPasteByID(ctx context.Context, id string) (model.Paste, error) {
	paste, err := s.storage.GetPasteByID(id)
	if err != nil {
		return model.Paste{}, err
	}
	return *paste, nil
}

func (s *pasteService) GetPasteByHash(ctx context.Context, hash string) (model.Paste, error) {
	paste, err := s.storage.GetPasteByHash(hash)
	if err != nil {
		return model.Paste{}, err
	}

	_ = s.statsService.IncrementViews(ctx, paste.ID)
	return *paste, nil
}

func (s *pasteService) DeletePaste(ctx context.Context, id string) error {
	err := s.storage.DeletePaste(id)
	if err == nil {
		_ = s.logger.LogChange("paste", id, "deleted")
	}
	return err
}

func (s *pasteService) ListPastes(ctx context.Context) ([]model.Paste, error) {
	return s.storage.GetAllPastes()
}
