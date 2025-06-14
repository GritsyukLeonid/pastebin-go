package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

// Моки

type mockUserStorage struct {
	users map[int64]model.User
}

func (m *mockUserStorage) SaveUser(u model.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserStorage) GetUserByID(id string) (*model.User, error) {
	for _, u := range m.users {
		if fmt.Sprintf("%d", u.ID) == id {
			return &u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserStorage) DeleteUser(id string) error {
	for uid, u := range m.users {
		if fmt.Sprintf("%d", u.ID) == id {
			delete(m.users, uid)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockUserStorage) GetAllUsers() ([]model.User, error) {
	var out []model.User
	for _, u := range m.users {
		out = append(out, u)
	}
	return out, nil
}

func (m *mockUserStorage) GetPasteByHash(hash string) (*model.Paste, error) {
	return nil, nil
}

func (m *mockUserStorage) IncrementStatsViews(id string) error {
	return nil
}

func (m *mockUserStorage) SavePaste(model.Paste) error                     { return nil }
func (m *mockUserStorage) GetPasteByID(string) (*model.Paste, error)       { return nil, nil }
func (m *mockUserStorage) DeletePaste(string) error                        { return nil }
func (m *mockUserStorage) GetAllPastes() ([]model.Paste, error)            { return nil, nil }
func (m *mockUserStorage) SaveShortURL(model.ShortURL) error               { return nil }
func (m *mockUserStorage) GetShortURLByID(string) (*model.ShortURL, error) { return nil, nil }
func (m *mockUserStorage) DeleteShortURL(string) error                     { return nil }
func (m *mockUserStorage) GetAllShortURLs() ([]model.ShortURL, error)      { return nil, nil }
func (m *mockUserStorage) SaveStats(model.Stats) error                     { return nil }
func (m *mockUserStorage) GetStatsByID(string) (*model.Stats, error)       { return nil, nil }
func (m *mockUserStorage) DeleteStats(string) error                        { return nil }
func (m *mockUserStorage) GetAllStats() ([]model.Stats, error)             { return nil, nil }

type userMockLogger struct{}

func (l *userMockLogger) LogChange(entity, id, action string) error {
	return nil
}

func setupUserService() UserService {
	storage := &mockUserStorage{users: make(map[int64]model.User)}
	logger := &userMockLogger{}
	return NewUserService(storage, logger)
}

// Тесты

func TestCreateUser(t *testing.T) {
	service := setupUserService()
	ctx := context.Background()

	user := model.User{Username: "alice"}
	created, err := service.CreateUser(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, created.Username)
	assert.NotZero(t, created.ID)
}

func TestGetUserByID(t *testing.T) {
	service := setupUserService()
	ctx := context.Background()

	user := model.User{Username: "bob"}
	created, _ := service.CreateUser(ctx, user)

	got, err := service.GetUserByID(ctx, fmt.Sprintf("%d", created.ID))

	assert.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, created.Username, got.Username)
}

func TestDeleteUser(t *testing.T) {
	service := setupUserService()
	ctx := context.Background()

	user := model.User{Username: "charlie"}
	created, _ := service.CreateUser(ctx, user)

	err := service.DeleteUser(ctx, fmt.Sprintf("%d", created.ID))
	assert.NoError(t, err)

	_, err = service.GetUserByID(ctx, fmt.Sprintf("%d", created.ID))
	assert.Error(t, err)
}
