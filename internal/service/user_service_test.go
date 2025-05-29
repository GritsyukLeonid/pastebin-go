package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateAndGetUser(t *testing.T) {
	type testCase struct {
		name     string
		username string
	}

	tests := []testCase{
		{"Basic user", "alice"},
		{"Empty username", ""},
	}

	service := NewUserService()
	ctx := context.Background()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user := model.User{Username: tc.username}
			created, err := service.CreateUser(ctx, user)
			assert.NoError(t, err)
			assert.Equal(t, tc.username, created.Username)

			got, err := service.GetUserByID(ctx, fmt.Sprintf("%d", created.ID))
			assert.NoError(t, err)
			assert.Equal(t, created.ID, got.ID)
			assert.Equal(t, tc.username, got.Username)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	service := NewUserService()

	user := model.User{Username: "bob"}
	created, err := service.CreateUser(ctx, user)
	assert.NoError(t, err)

	err = service.DeleteUser(ctx, fmt.Sprintf("%d", created.ID))
	assert.NoError(t, err)

	_, err = service.GetUserByID(ctx, fmt.Sprintf("%d", created.ID))
	assert.Error(t, err)
}
