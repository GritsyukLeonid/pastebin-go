package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

// Моки

type mockStatsStorage struct {
	stats map[string]model.Stats
}

func (m *mockStatsStorage) SaveStats(s model.Stats) error {
	m.stats[s.ID] = s
	return nil
}

func (m *mockStatsStorage) GetStatsByID(id string) (*model.Stats, error) {
	s, ok := m.stats[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &s, nil
}

func (m *mockStatsStorage) DeleteStats(id string) error {
	delete(m.stats, id)
	return nil
}

func (m *mockStatsStorage) GetAllStats() ([]model.Stats, error) {
	var out []model.Stats
	for _, s := range m.stats {
		out = append(out, s)
	}
	return out, nil
}

func (m *mockStatsStorage) SavePaste(model.Paste) error                     { return nil }
func (m *mockStatsStorage) GetPasteByID(string) (*model.Paste, error)       { return nil, nil }
func (m *mockStatsStorage) DeletePaste(string) error                        { return nil }
func (m *mockStatsStorage) GetAllPastes() ([]model.Paste, error)            { return nil, nil }
func (m *mockStatsStorage) SaveUser(model.User) error                       { return nil }
func (m *mockStatsStorage) GetUserByID(string) (*model.User, error)         { return nil, nil }
func (m *mockStatsStorage) DeleteUser(string) error                         { return nil }
func (m *mockStatsStorage) GetAllUsers() ([]model.User, error)              { return nil, nil }
func (m *mockStatsStorage) SaveShortURL(model.ShortURL) error               { return nil }
func (m *mockStatsStorage) GetShortURLByID(string) (*model.ShortURL, error) { return nil, nil }
func (m *mockStatsStorage) DeleteShortURL(string) error                     { return nil }
func (m *mockStatsStorage) GetAllShortURLs() ([]model.ShortURL, error)      { return nil, nil }

type statsMockLogger struct{}

func (l *statsMockLogger) LogChange(entity, id, action string) error {
	return nil
}

func setupStatsService() StatsService {
	storage := &mockStatsStorage{stats: make(map[string]model.Stats)}
	logger := &statsMockLogger{}
	return NewStatsService(storage, logger)
}

// Тесты

func TestCreateStats(t *testing.T) {
	service := setupStatsService()
	ctx := context.Background()

	s := model.Stats{ID: "s1", Views: 42}
	created, err := service.CreateStats(ctx, s)

	assert.NoError(t, err)
	assert.Equal(t, s.ID, created.ID)
	assert.Equal(t, s.Views, created.Views)
}

func TestGetStats(t *testing.T) {
	service := setupStatsService()
	ctx := context.Background()

	s := model.Stats{ID: "s2", Views: 77}
	_, _ = service.CreateStats(ctx, s)

	got, err := service.ListStats(ctx)
	assert.NoError(t, err)

	var found *model.Stats
	for _, stat := range got {
		if stat.ID == s.ID {
			found = &stat
			break
		}
	}
	assert.NotNil(t, found)
	assert.Equal(t, s.Views, found.Views)
}

func TestDeleteStats(t *testing.T) {
	service := setupStatsService()
	ctx := context.Background()

	s := model.Stats{ID: "s3", Views: 100}
	_, _ = service.CreateStats(ctx, s)

	err := service.DeleteStats(ctx, s.ID)
	assert.NoError(t, err)

	got, err := service.ListStats(ctx)
	assert.NoError(t, err)

	var found bool
	for _, stat := range got {
		if stat.ID == s.ID {
			found = true
			break
		}
	}
	assert.False(t, found)
}
