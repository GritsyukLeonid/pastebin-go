package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCreateAndGet(t *testing.T) {
	skipIfNotIntegration(t)

	user := map[string]interface{}{
		"username": "integration_user",
	}
	body, _ := json.Marshal(user)

	resp, err := http.Post("http://localhost:8080/api/user", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created struct {
		ID       int64    `json:"id"`
		Username string   `json:"username"`
		Posts    []string `json:"posts"`
	}
	err = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, "integration_user", created.Username)
	assert.Greater(t, created.ID, int64(0))

	resp, err = http.Get("http://localhost:8080/api/user/" + fmt.Sprintf("%d", created.ID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched struct {
		ID       int64    `json:"id"`
		Username string   `json:"username"`
		Posts    []string `json:"posts"`
	}
	err = json.NewDecoder(resp.Body).Decode(&fetched)
	resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Username, fetched.Username)
	assert.Empty(t, fetched.Posts)
}
