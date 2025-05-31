package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/logging"
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type pasteService struct {
	storage repository.MongoStorageInterface
	logger  logging.Logger
}

func NewPasteService(storage repository.MongoStorageInterface, logger logging.Logger) PasteService {
	return &pasteService{storage: storage, logger: logger}
}

func (s *pasteService) CreatePaste(ctx context.Context, p model.Paste) (model.Paste, error) {
	p.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	p.CreatedAt = time.Now()

	if err := s.storage.SavePaste(p); err != nil {
		return model.Paste{}, err
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
