package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	saveFunc      func(model.Paste) error
	getByIDFunc   func(string) (*model.Paste, error)
	getAllFunc    func() ([]model.Paste, error)
	deleteFunc    func(string) error
	getByHashFunc func(string) (*model.Paste, error)
}

func (m *mockStorage) SavePaste(p model.Paste) error                { return m.saveFunc(p) }
func (m *mockStorage) GetPasteByID(id string) (*model.Paste, error) { return m.getByIDFunc(id) }
func (m *mockStorage) GetAllPastes() ([]model.Paste, error)         { return m.getAllFunc() }
func (m *mockStorage) DeletePaste(id string) error                  { return m.deleteFunc(id) }
func (m *mockStorage) GetPasteByHash(hash string) (*model.Paste, error) {
	return m.getByHashFunc(hash)
}

func (m *mockStorage) SaveStats(model.Stats) error                     { return nil }
func (m *mockStorage) GetStatsByID(string) (*model.Stats, error)       { return nil, nil }
func (m *mockStorage) DeleteStats(string) error                        { return nil }
func (m *mockStorage) GetAllStats() ([]model.Stats, error)             { return nil, nil }
func (m *mockStorage) SaveUser(model.User) error                       { return nil }
func (m *mockStorage) GetUserByID(string) (*model.User, error)         { return nil, nil }
func (m *mockStorage) DeleteUser(string) error                         { return nil }
func (m *mockStorage) GetAllUsers() ([]model.User, error)              { return nil, nil }
func (m *mockStorage) SaveShortURL(model.ShortURL) error               { return nil }
func (m *mockStorage) GetShortURLByID(string) (*model.ShortURL, error) { return nil, nil }
func (m *mockStorage) DeleteShortURL(string) error                     { return nil }
func (m *mockStorage) GetAllShortURLs() ([]model.ShortURL, error)      { return nil, nil }

type mockLogger struct{}

func (m *mockLogger) LogChange(entity, id, action string) error { return nil }

type mockStatsService struct{}

func (m *mockStatsService) CreateStats(ctx context.Context, s model.Stats) (model.Stats, error) {
	return s, nil
}
func (m *mockStatsService) GetStatsByID(ctx context.Context, id string) (model.Stats, error) {
	return model.Stats{ID: id}, nil
}
func (m *mockStatsService) DeleteStats(ctx context.Context, id string) error {
	return nil
}
func (m *mockStatsService) ListStats(ctx context.Context) ([]model.Stats, error) {
	return nil, nil
}
func (m *mockStatsService) IncrementViews(ctx context.Context, id string) error {
	return nil
}

func TestCreatePaste(t *testing.T) {
	mockStorage := &mockStorage{
		saveFunc: func(p model.Paste) error {
			if p.Content == "fail" {
				return errors.New("save error")
			}
			return nil
		},
	}
	mockLogger := &mockLogger{}
	mockStats := &mockStatsService{}

	svc := NewPasteService(mockStorage, mockLogger, mockStats)

	ctx := context.Background()
	paste := model.Paste{Content: "test content", ExpiresAt: time.Now().Add(1 * time.Hour)}
	created, err := svc.CreatePaste(ctx, paste)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.CreatedAt)
	assert.NotEmpty(t, created.Hash)
}

func TestGetPasteByHash(t *testing.T) {
	expected := &model.Paste{ID: "123", Hash: "abc", Content: "test"}

	mockStorage := &mockStorage{
		getByHashFunc: func(h string) (*model.Paste, error) {
			if h == "abc" {
				return expected, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockLogger := &mockLogger{}
	mockStats := &mockStatsService{}

	svc := NewPasteService(mockStorage, mockLogger, mockStats)

	ctx := context.Background()
	res, err := svc.GetPasteByHash(ctx, "abc")
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, res.ID)
}
