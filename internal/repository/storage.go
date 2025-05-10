package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
)

var (
	pasteMutex sync.Mutex
	userMutex  sync.Mutex
	statsMutex sync.Mutex
	urlMutex   sync.Mutex

	Pastes   []*model.Paste
	Users    []*model.User
	StatsSet []*model.Stats
	URLs     []*model.ShortURL

	prevCounts = map[string]int{
		"Paste":    0,
		"User":     0,
		"Stats":    0,
		"ShortURL": 0,
	}
)

func StoreObject(obj model.Storable) error {
	switch v := obj.(type) {
	case *model.Paste:
		pasteMutex.Lock()
		defer pasteMutex.Unlock()
		Pastes = append(Pastes, v)
		return saveJSON("pastes.json", Pastes)
	case *model.User:
		userMutex.Lock()
		defer userMutex.Unlock()
		Users = append(Users, v)
		return saveJSON("users.json", Users)
	case *model.Stats:
		statsMutex.Lock()
		defer statsMutex.Unlock()
		StatsSet = append(StatsSet, v)
		return saveJSON("stats.json", StatsSet)
	case *model.ShortURL:
		urlMutex.Lock()
		defer urlMutex.Unlock()
		URLs = append(URLs, v)
		return saveJSON("urls.json", URLs)
	default:
		return fmt.Errorf("unknown type: %s", v.GetTypeName())
	}
}

func DetectNewObjects() map[string][]string {
	result := make(map[string][]string)

	pasteMutex.Lock()
	if len(Pastes) > prevCounts["Paste"] {
		for i := prevCounts["Paste"]; i < len(Pastes); i++ {
			result["Paste"] = append(result["Paste"], fmt.Sprintf("%+v", Pastes[i]))
		}
		prevCounts["Paste"] = len(Pastes)
	}
	pasteMutex.Unlock()

	userMutex.Lock()
	if len(Users) > prevCounts["User"] {
		for i := prevCounts["User"]; i < len(Users); i++ {
			result["User"] = append(result["User"], fmt.Sprintf("%+v", Users[i]))
		}
		prevCounts["User"] = len(Users)
	}
	userMutex.Unlock()

	statsMutex.Lock()
	if len(StatsSet) > prevCounts["Stats"] {
		for i := prevCounts["Stats"]; i < len(StatsSet); i++ {
			result["Stats"] = append(result["Stats"], fmt.Sprintf("%+v", StatsSet[i]))
		}
		prevCounts["Stats"] = len(StatsSet)
	}
	statsMutex.Unlock()

	urlMutex.Lock()
	if len(URLs) > prevCounts["ShortURL"] {
		for i := prevCounts["ShortURL"]; i < len(URLs); i++ {
			result["ShortURL"] = append(result["ShortURL"], fmt.Sprintf("%+v", URLs[i]))
		}
		prevCounts["ShortURL"] = len(URLs)
	}
	urlMutex.Unlock()

	return result
}

func saveJSON(filename string, data any) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func LoadData() {
	loadJSON("pastes.json", &Pastes)
	loadJSON("users.json", &Users)
	loadJSON("stats.json", &StatsSet)
	loadJSON("urls.json", &URLs)

	prevCounts["Paste"] = len(Pastes)
	prevCounts["User"] = len(Users)
	prevCounts["Stats"] = len(StatsSet)
	prevCounts["ShortURL"] = len(URLs)
}

func loadJSON(filename string, target any) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}
	defer file.Close()

	_ = json.NewDecoder(file).Decode(target)
}
