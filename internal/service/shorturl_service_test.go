package service

import (
	"context"
	"testing"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestShortURLService_CreateAndGet(t *testing.T) {
	type testCase struct {
		name     string
		id       string
		original string
	}

	tests := []testCase{
		{"Simple short url", "short1", "https://example.com"},
		{"Empty original", "short2", ""},
	}

	service := NewShortURLService()
	ctx := context.Background()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := model.ShortURL{ID: tc.id, Original: tc.original}
			created, err := service.CreateShortURL(ctx, input)
			assert.NoError(t, err)
			assert.Equal(t, tc.original, created.Original)

			got, err := service.GetShortURL(ctx, created.ID)
			assert.NoError(t, err)
			assert.Equal(t, created.ID, got.ID)
			assert.Equal(t, tc.original, got.Original)
		})
	}
}

func TestShortURLService_Delete(t *testing.T) {
	ctx := context.Background()
	service := NewShortURLService()

	su := model.ShortURL{ID: "d1", Original: "https://delete.me"}
	created, err := service.CreateShortURL(ctx, su)
	assert.NoError(t, err)

	err = service.DeleteShortURL(ctx, created.ID)
	assert.NoError(t, err)

	_, err = service.GetShortURL(ctx, created.ID)
	assert.Error(t, err)
}
