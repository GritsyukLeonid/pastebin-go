package repository

import (
	"encoding/json"
	"fmt"
	"log"
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

func StoreObject(obj model.Storable) {
	switch v := obj.(type) {
	case *model.Paste:
		pasteMutex.Lock()
		defer pasteMutex.Unlock()
		Pastes = append(Pastes, v)
		saveJSON("pastes.json", Pastes)
	case *model.User:
		userMutex.Lock()
		defer userMutex.Unlock()
		Users = append(Users, v)
		saveJSON("users.json", Users)
	case *model.Stats:
		statsMutex.Lock()
		defer statsMutex.Unlock()
		StatsSet = append(StatsSet, v)
		saveJSON("stats.json", StatsSet)
	case *model.ShortURL:
		urlMutex.Lock()
		defer urlMutex.Unlock()
		URLs = append(URLs, v)
		saveJSON("urls.json", URLs)
	default:
		log.Println("Unknown type:", v.GetTypeName())
	}
}

func CheckAndLogChanges() {
	logNew("Paste", &pasteMutex, func() int { return len(Pastes) }, func(i int) string {
		return fmt.Sprintf("Paste: %+v", Pastes[i])
	})
	logNew("User", &userMutex, func() int { return len(Users) }, func(i int) string {
		return fmt.Sprintf("User: %+v", Users[i])
	})
	logNew("Stats", &statsMutex, func() int { return len(StatsSet) }, func(i int) string {
		return fmt.Sprintf("Stats: %+v", StatsSet[i])
	})
	logNew("ShortURL", &urlMutex, func() int { return len(URLs) }, func(i int) string {
		return fmt.Sprintf("ShortURL: %+v", URLs[i])
	})
}

func logNew(key string, mu *sync.Mutex, countFn func() int, formatFn func(i int) string) {
	mu.Lock()
	defer mu.Unlock()
	current := countFn()
	prev := prevCounts[key]
	if current > prev {
		for i := prev; i < current; i++ {
			log.Println(formatFn(i))
		}
		prevCounts[key] = current
	}
}

func saveJSON(filename string, data any) {
	file, err := os.Create(filename)
	if err != nil {
		log.Println("Error creating file:", filename, err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		log.Println("Error encoding data to", filename, err)
	}
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
		log.Println("Error opening file:", filename, err)
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(target); err != nil {
		log.Println("Error decoding", filename, err)
	}
}
