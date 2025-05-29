package service

import (
	"context"
	"testing"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPasteService_CreateAndGetPaste(t *testing.T) {
	type testCase struct {
		name    string
		content string
		expect  string
	}

	tests := []testCase{
		{
			name:    "Valid paste creation",
			content: "Hello, World!",
			expect:  "Hello, World!",
		},
		{
			name:    "Empty content",
			content: "",
			expect:  "",
		},
	}

	service := NewPasteService()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			paste := model.Paste{
				Content:   tc.content,
				CreatedAt: time.Now(),
			}

			created, err := service.CreatePaste(ctx, paste)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, created.Content)

			fetched, err := service.GetPasteByID(ctx, created.ID)
			assert.NoError(t, err)
			assert.Equal(t, created.ID, fetched.ID)
			assert.Equal(t, tc.expect, fetched.Content)
		})
	}
}

func TestPasteService_DeletePaste(t *testing.T) {
	ctx := context.Background()
	service := NewPasteService()

	paste := model.Paste{Content: "To be deleted", CreatedAt: time.Now()}
	created, err := service.CreatePaste(ctx, paste)
	assert.NoError(t, err)

	err = service.DeletePaste(ctx, created.ID)
	assert.NoError(t, err)

	_, err = service.GetPasteByID(ctx, created.ID)
	assert.Error(t, err)
}

func TestPasteService_GetNonexistentPaste(t *testing.T) {
	ctx := context.Background()
	service := NewPasteService()

	_, err := service.GetPasteByID(ctx, "nonexistent-id")
	assert.Error(t, err)
}
