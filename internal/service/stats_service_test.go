package service

import (
	"context"
	"testing"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestStatsService_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	service := NewStatsService()

	stats := model.Stats{ID: "s1", Views: 42}
	created, err := service.CreateStats(ctx, stats)
	assert.NoError(t, err)
	assert.Equal(t, stats.ID, created.ID)
	assert.Equal(t, stats.Views, created.Views)

	got, err := service.GetStats(ctx)
	assert.NoError(t, err)

	var found *model.Stats
	for _, s := range got {
		if s.ID == created.ID {
			found = &s
			break
		}
	}
	assert.NotNil(t, found)
	assert.Equal(t, created.Views, found.Views)
}

func TestStatsService_Delete(t *testing.T) {
	ctx := context.Background()
	service := NewStatsService()

	stats := model.Stats{ID: "s2", Views: 100}
	created, err := service.CreateStats(ctx, stats)
	assert.NoError(t, err)

	err = service.DeleteStats(ctx, created.ID)
	assert.NoError(t, err)

	got, err := service.GetStats(ctx)
	assert.NoError(t, err)

	var found bool
	for _, s := range got {
		if s.ID == created.ID {
			found = true
			break
		}
	}
	assert.False(t, found)
}
