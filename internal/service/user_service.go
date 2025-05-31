package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
)

type userService struct{}

func NewUserService() UserService {
	repository.LoadData()
	return &userService{}
}

func (s *userService) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	u.ID = int64(len(repository.GetAllUsers()) + 1)

	if err := repository.AddUser(&u); err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (model.User, error) {
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return model.User{}, fmt.Errorf("invalid user ID: %v", err)
	}
	u, err := repository.GetUserByID(uid)
	if err != nil {
		return model.User{}, err
	}
	return *u, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}
	return repository.DeleteUser(uid)
}

func (s *userService) ListUsers(ctx context.Context) ([]model.User, error) {
	raw := repository.GetAllUsers()
	users := make([]model.User, len(raw))
	for i, u := range raw {
		users[i] = *u
	}
	return users, nil
}
