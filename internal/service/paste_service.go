package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type pasteService struct{}

func NewPasteService() PasteService {
	repository.LoadData()
	return &pasteService{}
}

func (s *pasteService) CreatePaste(ctx context.Context, p model.Paste) (model.Paste, error) {
	p.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	p.CreatedAt = time.Now()

	if err := repository.StoreObject(&p); err != nil {
		return model.Paste{}, err
	}
	return p, nil
}

func (s *pasteService) GetPasteByID(ctx context.Context, id string) (model.Paste, error) {
	p, err := repository.GetPasteByID(id)
	if err != nil {
		return model.Paste{}, err
	}
	return *p, nil
}

func (s *pasteService) DeletePaste(ctx context.Context, id string) error {
	return repository.DeletePaste(id)
}

func (s *pasteService) ListPastes(ctx context.Context) ([]model.Paste, error) {
	plist := repository.GetAllPastes()
	out := make([]model.Paste, len(plist))
	for i, v := range plist {
		out[i] = *v
	}
	return out, nil
}
