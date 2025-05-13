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

func GetAllPastes() []*model.Paste {
	pasteMutex.Lock()
	defer pasteMutex.Unlock()
	return Pastes
}

func GetPasteByID(id string) (*model.Paste, error) {
	pasteMutex.Lock()
	defer pasteMutex.Unlock()
	for _, p := range Pastes {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("paste not found")
}

func UpdatePaste(id string, newPaste *model.Paste) error {
	pasteMutex.Lock()
	defer pasteMutex.Unlock()
	for i, p := range Pastes {
		if p.ID == id {
			Pastes[i] = newPaste
			return saveJSON("pastes.json", Pastes)
		}
	}
	return fmt.Errorf("paste not found")
}

func DeletePaste(id string) error {
	pasteMutex.Lock()
	defer pasteMutex.Unlock()
	for i, p := range Pastes {
		if p.ID == id {
			Pastes = append(Pastes[:i], Pastes[i+1:]...)
			return saveJSON("pastes.json", Pastes)
		}
	}
	return fmt.Errorf("paste not found")
}

func GetAllUsers() []*model.User {
	userMutex.Lock()
	defer userMutex.Unlock()
	return Users
}
func GetUserByID(id int64) (*model.User, error) {
	userMutex.Lock()
	defer userMutex.Unlock()
	for _, u := range Users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
func UpdateUser(id int64, newUser *model.User) error {
	userMutex.Lock()
	defer userMutex.Unlock()
	for i, u := range Users {
		if u.ID == id {
			Users[i] = newUser
			return saveJSON("users.json", Users)
		}
	}
	return fmt.Errorf("user not found")
}
func DeleteUser(id int64) error {
	userMutex.Lock()
	defer userMutex.Unlock()
	for i, u := range Users {
		if u.ID == id {
			Users = append(Users[:i], Users[i+1:]...)
			return saveJSON("users.json", Users)
		}
	}
	return fmt.Errorf("user not found")
}
func AddUser(u *model.User) error {
	userMutex.Lock()
	defer userMutex.Unlock()
	Users = append(Users, u)
	return saveJSON("users.json", Users)
}

func GetAllStats() []*model.Stats {
	statsMutex.Lock()
	defer statsMutex.Unlock()
	return StatsSet
}
func GetStatsByID(id string) (*model.Stats, error) {
	statsMutex.Lock()
	defer statsMutex.Unlock()
	for _, s := range StatsSet {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, fmt.Errorf("stats not found")
}
func UpdateStats(id string, newStats *model.Stats) error {
	statsMutex.Lock()
	defer statsMutex.Unlock()
	for i, s := range StatsSet {
		if s.ID == id {
			StatsSet[i] = newStats
			return saveJSON("stats.json", StatsSet)
		}
	}
	return fmt.Errorf("stats not found")
}
func DeleteStats(id string) error {
	statsMutex.Lock()
	defer statsMutex.Unlock()
	for i, s := range StatsSet {
		if s.ID == id {
			StatsSet = append(StatsSet[:i], StatsSet[i+1:]...)
			return saveJSON("stats.json", StatsSet)
		}
	}
	return fmt.Errorf("stats not found")
}

func GetAllShortURLs() []*model.ShortURL {
	urlMutex.Lock()
	defer urlMutex.Unlock()
	return URLs
}
func GetShortURLByID(id string) (*model.ShortURL, error) {
	urlMutex.Lock()
	defer urlMutex.Unlock()
	for _, u := range URLs {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fmt.Errorf("shorturl not found")
}
func UpdateShortURL(id string, newURL *model.ShortURL) error {
	urlMutex.Lock()
	defer urlMutex.Unlock()
	for i, u := range URLs {
		if u.ID == id {
			URLs[i] = newURL
			return saveJSON("urls.json", URLs)
		}
	}
	return fmt.Errorf("shorturl not found")
}
func DeleteShortURL(id string) error {
	urlMutex.Lock()
	defer urlMutex.Unlock()
	for i, u := range URLs {
		if u.ID == id {
			URLs = append(URLs[:i], URLs[i+1:]...)
			return saveJSON("urls.json", URLs)
		}
	}
	return fmt.Errorf("shorturl not found")
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
