package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/logging"
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type shortURLService struct {
	storage repository.StorageInterface
	logger  logging.Logger
}

func NewShortURLService(storage repository.StorageInterface, logger logging.Logger) ShortURLService {
	return &shortURLService{storage: storage, logger: logger}
}

func (s *shortURLService) CreateShortURL(ctx context.Context, u model.ShortURL) (model.ShortURL, error) {
	u.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	if err := s.storage.SaveShortURL(u); err != nil {
		return model.ShortURL{}, err
	}

	_ = s.logger.LogChange("shorturl", u.ID, "created")
	return u, nil
}

func (s *shortURLService) GetShortURLByID(ctx context.Context, id string) (model.ShortURL, error) {
	url, err := s.storage.GetShortURLByID(id)
	if err != nil {
		return model.ShortURL{}, err
	}
	return *url, nil
}

func (s *shortURLService) DeleteShortURL(ctx context.Context, id string) error {
	err := s.storage.DeleteShortURL(id)
	if err == nil {
		_ = s.logger.LogChange("shorturl", id, "deleted")
	}
	return err
}

func (s *shortURLService) ListShortURLs(ctx context.Context) ([]model.ShortURL, error) {
	return s.storage.GetAllShortURLs()
}
