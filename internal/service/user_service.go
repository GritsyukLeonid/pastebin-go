package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/logging"
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type userService struct {
	storage repository.StorageInterface
	logger  logging.Logger
}

func NewUserService(storage repository.StorageInterface, logger logging.Logger) UserService {
	return &userService{storage: storage, logger: logger}
}

func (s *userService) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	u.ID = time.Now().UnixNano()

	if err := s.storage.SaveUser(u); err != nil {
		return model.User{}, err
	}

	_ = s.logger.LogChange("user", fmt.Sprintf("%d", u.ID), "created")
	return u, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (model.User, error) {
	user, err := s.storage.GetUserByID(id)
	if err != nil {
		return model.User{}, err
	}
	return *user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	err := s.storage.DeleteUser(id)
	if err == nil {
		_ = s.logger.LogChange("user", id, "deleted")
	}
	return err
}

func (s *userService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.storage.GetAllUsers()
}
