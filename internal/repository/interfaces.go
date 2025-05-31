package repository

import (
	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

type MongoStorageInterface interface {
	// Paste
	SavePaste(model.Paste) error
	GetPasteByID(string) (*model.Paste, error)
	DeletePaste(string) error
	GetAllPastes() ([]model.Paste, error)

	// User
	SaveUser(model.User) error
	GetUserByID(string) (*model.User, error)
	DeleteUser(string) error
	GetAllUsers() ([]model.User, error)

	// ShortURL
	SaveShortURL(model.ShortURL) error
	GetShortURLByID(string) (*model.ShortURL, error)
	DeleteShortURL(string) error
	GetAllShortURLs() ([]model.ShortURL, error)

	// Stats
	SaveStats(model.Stats) error
	GetStatsByID(string) (*model.Stats, error)
	DeleteStats(string) error
	GetAllStats() ([]model.Stats, error)
}
