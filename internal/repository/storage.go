package repository

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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

func StoreFromChannel(ctx context.Context, ch <-chan model.Storable) {
	for {
		select {
		case <-ctx.Done():
			return
		case obj, ok := <-ch:
			if !ok {
				return
			}
			switch v := obj.(type) {
			case *model.Paste:
				pasteMutex.Lock()
				Pastes = append(Pastes, v)
				pasteMutex.Unlock()
			case *model.User:
				userMutex.Lock()
				Users = append(Users, v)
				userMutex.Unlock()
			case *model.Stats:
				statsMutex.Lock()
				StatsSet = append(StatsSet, v)
				statsMutex.Unlock()
			case *model.ShortURL:
				urlMutex.Lock()
				URLs = append(URLs, v)
				urlMutex.Unlock()
			default:
				log.Println("Unknown type:", v.GetTypeName())
			}
		}
	}
}

func LogChanges(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logNew("Paste", &pasteMutex, func() int {
				return len(Pastes)
			}, func(i int) string {
				return fmt.Sprintf("Paste: %+v", Pastes[i])
			})
			logNew("User", &userMutex, func() int {
				return len(Users)
			}, func(i int) string {
				return fmt.Sprintf("User: %+v", Users[i])
			})
			logNew("Stats", &statsMutex, func() int {
				return len(StatsSet)
			}, func(i int) string {
				return fmt.Sprintf("Stats: %+v", StatsSet[i])
			})
			logNew("ShortURL", &urlMutex, func() int {
				return len(URLs)
			}, func(i int) string {
				return fmt.Sprintf("ShortURL: %+v", URLs[i])
			})
		}
	}
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
