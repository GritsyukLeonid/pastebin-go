package repository

import (
	"fmt"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

var (
	Pastes   []model.Paste
	Users    []model.User
	StatsSet []model.Stats
	URLs     []model.ShortURL
)

func Store(obj model.Storable) {
	switch v := obj.(type) {
	case *model.Paste:
		Pastes = append(Pastes, *v)
	case *model.User:
		Users = append(Users, *v)
	case *model.Stats:
		StatsSet = append(StatsSet, *v)
	case *model.ShortURL:
		URLs = append(URLs, *v)
	default:
		fmt.Println("Unknown type:", v.GetTypeName())
	}
}
