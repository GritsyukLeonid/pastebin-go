package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

// Моки

type mockShortURLStorage struct {
	shorts map[string]model.ShortURL
}

func (m *mockShortURLStorage) SaveShortURL(u model.ShortURL) error {
	m.shorts[u.ID] = u
	return nil
}

func (m *mockShortURLStorage) GetShortURLByID(id string) (*model.ShortURL, error) {
	u, ok := m.shorts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &u, nil
}

func (m *mockShortURLStorage) DeleteShortURL(id string) error {
	delete(m.shorts, id)
	return nil
}

func (m *mockShortURLStorage) GetAllShortURLs() ([]model.ShortURL, error) {
	out := make([]model.ShortURL, 0, len(m.shorts))
	for _, v := range m.shorts {
		out = append(out, v)
	}
	return out, nil
}

func (m *mockShortURLStorage) SavePaste(model.Paste) error               { return nil }
func (m *mockShortURLStorage) GetPasteByID(string) (*model.Paste, error) { return nil, nil }
func (m *mockShortURLStorage) DeletePaste(string) error                  { return nil }
func (m *mockShortURLStorage) GetAllPastes() ([]model.Paste, error)      { return nil, nil }
func (m *mockShortURLStorage) SaveUser(model.User) error                 { return nil }
func (m *mockShortURLStorage) GetUserByID(string) (*model.User, error)   { return nil, nil }
func (m *mockShortURLStorage) DeleteUser(string) error                   { return nil }
func (m *mockShortURLStorage) GetAllUsers() ([]model.User, error)        { return nil, nil }
func (m *mockShortURLStorage) SaveStats(model.Stats) error               { return nil }
func (m *mockShortURLStorage) GetStatsByID(string) (*model.Stats, error) { return nil, nil }
func (m *mockShortURLStorage) DeleteStats(string) error                  { return nil }
func (m *mockShortURLStorage) GetAllStats() ([]model.Stats, error)       { return nil, nil }

type shortMockLogger struct{}

func (l *shortMockLogger) LogChange(entity, id, action string) error {
	return nil
}

func setupShortService() ShortURLService {
	storage := &mockShortURLStorage{shorts: make(map[string]model.ShortURL)}
	logger := &shortMockLogger{}
	return NewShortURLService(storage, logger)
}

// Тесты

func TestCreateShortURL(t *testing.T) {
	service := setupShortService()
	ctx := context.Background()

	input := model.ShortURL{ID: "s1", Original: "https://example.com"}
	created, err := service.CreateShortURL(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, input.Original, created.Original)
	assert.Equal(t, input.ID, created.ID)
}

func TestGetShortURLByID(t *testing.T) {
	service := setupShortService()
	ctx := context.Background()

	input := model.ShortURL{ID: "s2", Original: "https://go.dev"}
	_, _ = service.CreateShortURL(ctx, input)

	got, err := service.GetShortURLByID(ctx, "s2")

	assert.NoError(t, err)
	assert.Equal(t, "s2", got.ID)
	assert.Equal(t, "https://go.dev", got.Original)
}

func TestDeleteShortURL(t *testing.T) {
	service := setupShortService()
	ctx := context.Background()

	input := model.ShortURL{ID: "delme", Original: "https://delete.com"}
	_, _ = service.CreateShortURL(ctx, input)

	err := service.DeleteShortURL(ctx, "delme")
	assert.NoError(t, err)

	_, err = service.GetShortURLByID(ctx, "delme")
	assert.Error(t, err)
}
