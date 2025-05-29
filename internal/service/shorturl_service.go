package service

import (
	"context"
	"fmt"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type shortURLService struct{}

func NewShortURLService() ShortURLService {
	return &shortURLService{}
}

func (s *shortURLService) CreateShortURL(ctx context.Context, su model.ShortURL) (model.ShortURL, error) {
	if err := repository.StoreObject(&su); err != nil {
		return model.ShortURL{}, err
	}
	return su, nil
}

func (s *shortURLService) GetShortURL(ctx context.Context, id string) (model.ShortURL, error) {
	su, err := repository.GetShortURLByID(id)
	if err != nil {
		return model.ShortURL{}, fmt.Errorf("short url not found: %w", err)
	}
	return *su, nil
}

func (s *shortURLService) DeleteShortURL(ctx context.Context, id string) error {
	return repository.DeleteShortURL(id)
}

func (s *shortURLService) ListShortURLs(ctx context.Context) ([]*model.ShortURL, error) {
	return repository.GetAllShortURLs(), nil
}
