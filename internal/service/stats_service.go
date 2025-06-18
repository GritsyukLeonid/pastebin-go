package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/logging"
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type statsService struct {
	storage repository.StorageInterface
	logger  logging.Logger
}

func NewStatsService(storage repository.StorageInterface, logger logging.Logger) StatsService {
	return &statsService{storage: storage, logger: logger}
}

func (s *statsService) CreateStats(ctx context.Context, stat model.Stats) (model.Stats, error) {
	stat.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	if err := s.storage.SaveStats(stat); err != nil {
		return model.Stats{}, err
	}

	_ = s.logger.LogChange("stats", stat.ID, "created")
	return stat, nil
}

func (s *statsService) GetStatsByID(ctx context.Context, id string) (model.Stats, error) {
	stat, err := s.storage.GetStatsByID(id)
	if err != nil {
		return model.Stats{}, err
	}
	return *stat, nil
}

func (s *statsService) DeleteStats(ctx context.Context, id string) error {
	err := s.storage.DeleteStats(id)
	if err == nil {
		_ = s.logger.LogChange("stats", id, "deleted")
	}
	return err
}

func (s *statsService) ListStats(ctx context.Context) ([]model.Stats, error) {
	return s.storage.GetAllStats()
}

func (s *statsService) IncrementViews(ctx context.Context, id string) error {
	return s.storage.IncrementStatsViews(id)
}

func (s *statsService) ListTopStats(ctx context.Context, limit int) ([]model.Stats, error) {
	allStats, err := s.storage.GetAllStats()
	if err != nil {
		return nil, err
	}

	// Сортируем по Views (по убыванию)
	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].Views > allStats[j].Views
	})

	if len(allStats) > limit {
		return allStats[:limit], nil
	}
	return allStats, nil
}
