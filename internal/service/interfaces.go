package service

import (
	"context"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

type PasteService interface {
	CreatePaste(ctx context.Context, p model.Paste) (model.Paste, error)
	GetPasteByID(ctx context.Context, id string) (model.Paste, error)
	DeletePaste(ctx context.Context, id string) error
	ListPastes(ctx context.Context) ([]model.Paste, error)
}

type UserService interface {
	CreateUser(ctx context.Context, u model.User) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]model.User, error)
}

type ShortURLService interface {
	CreateShortURL(ctx context.Context, u model.ShortURL) (model.ShortURL, error)
	GetShortURLByID(ctx context.Context, id string) (model.ShortURL, error)
	DeleteShortURL(ctx context.Context, id string) error
	ListShortURLs(ctx context.Context) ([]model.ShortURL, error)
}

type StatsService interface {
	CreateStats(ctx context.Context, s model.Stats) (model.Stats, error)
	GetStatsByID(ctx context.Context, id string) (model.Stats, error)
	DeleteStats(ctx context.Context, id string) error
	ListStats(ctx context.Context) ([]model.Stats, error)
}
