package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

// Моки

type mockStorage struct {
	pastes map[string]model.Paste
}

func (m *mockStorage) SavePaste(p model.Paste) error {
	m.pastes[p.ID] = p
	return nil
}

func (m *mockStorage) GetPasteByID(id string) (*model.Paste, error) {
	p, ok := m.pastes[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &p, nil
}

func (m *mockStorage) DeletePaste(id string) error {
	delete(m.pastes, id)
	return nil
}

func (m *mockStorage) GetAllPastes() ([]model.Paste, error) {
	result := make([]model.Paste, 0, len(m.pastes))
	for _, p := range m.pastes {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockStorage) SaveUser(model.User) error                       { return nil }
func (m *mockStorage) GetUserByID(string) (*model.User, error)         { return nil, nil }
func (m *mockStorage) DeleteUser(string) error                         { return nil }
func (m *mockStorage) GetAllUsers() ([]model.User, error)              { return nil, nil }
func (m *mockStorage) SaveShortURL(model.ShortURL) error               { return nil }
func (m *mockStorage) GetShortURLByID(string) (*model.ShortURL, error) { return nil, nil }
func (m *mockStorage) DeleteShortURL(string) error                     { return nil }
func (m *mockStorage) GetAllShortURLs() ([]model.ShortURL, error)      { return nil, nil }
func (m *mockStorage) SaveStats(model.Stats) error                     { return nil }
func (m *mockStorage) GetStatsByID(string) (*model.Stats, error)       { return nil, nil }
func (m *mockStorage) DeleteStats(string) error                        { return nil }
func (m *mockStorage) GetAllStats() ([]model.Stats, error)             { return nil, nil }

type mockLogger struct{}

func (l *mockLogger) LogChange(entity, id, action string) error {
	return nil
}

// Тесты

func setupTestPasteService() PasteService {
	storage := &mockStorage{pastes: make(map[string]model.Paste)}
	logger := &mockLogger{}
	return NewPasteService(storage, logger)
}

func TestCreatePaste(t *testing.T) {
	service := setupTestPasteService()
	ctx := context.Background()

	paste := model.Paste{Content: "Hello", CreatedAt: time.Now()}
	created, err := service.CreatePaste(ctx, paste)

	assert.NoError(t, err)
	assert.Equal(t, paste.Content, created.Content)
	assert.NotEmpty(t, created.ID)
}

func TestGetPasteByID(t *testing.T) {
	service := setupTestPasteService()
	ctx := context.Background()

	paste := model.Paste{Content: "Fetch me", CreatedAt: time.Now()}
	created, _ := service.CreatePaste(ctx, paste)

	fetched, err := service.GetPasteByID(ctx, created.ID)

	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, paste.Content, fetched.Content)
}

func TestDeletePaste(t *testing.T) {
	service := setupTestPasteService()
	ctx := context.Background()

	paste := model.Paste{Content: "To delete", CreatedAt: time.Now()}
	created, _ := service.CreatePaste(ctx, paste)

	err := service.DeletePaste(ctx, created.ID)
	assert.NoError(t, err)

	_, err = service.GetPasteByID(ctx, created.ID)
	assert.Error(t, err)
}

func TestGetNonexistentPaste(t *testing.T) {
	service := setupTestPasteService()
	ctx := context.Background()

	_, err := service.GetPasteByID(ctx, "nonexistent-id")
	assert.Error(t, err)
}
