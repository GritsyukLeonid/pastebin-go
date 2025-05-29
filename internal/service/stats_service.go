package service

import (
	"context"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type statsService struct{}

func NewStatsService() StatsService {
	repository.LoadData()
	return &statsService{}
}

func (s *statsService) CreateStats(ctx context.Context, stats model.Stats) (model.Stats, error) {
	if err := repository.StoreObject(&stats); err != nil {
		return model.Stats{}, err
	}
	return stats, nil
}

func (s *statsService) RecordView(ctx context.Context, id string) error {
	st, err := repository.GetStatsByID(id)
	if err != nil {
		return err
	}

	st.Views++
	return repository.UpdateStats(id, st)
}

func (s *statsService) GetStats(ctx context.Context) ([]model.Stats, error) {
	list := repository.GetAllStats()
	stats := make([]model.Stats, len(list))
	for i, s := range list {
		stats[i] = *s
	}
	return stats, nil
}

func (s *statsService) DeleteStats(ctx context.Context, id string) error {
	return repository.DeleteStats(id)
}
