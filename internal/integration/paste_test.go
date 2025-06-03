package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasteLifecycle(t *testing.T) {
	skipIfNotIntegration(t)

	paste := map[string]interface{}{
		"user_id":            "test-user",
		"content":            "integration test content",
		"expiration_minutes": 10,
	}
	body, _ := json.Marshal(paste)

	resp, err := http.Post("http://localhost:8080/api/paste", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	assert.NoError(t, err)

	pasteID, ok := created["id"].(string)
	assert.True(t, ok, "ID должен быть строкой")

	resp, err = http.Get("http://localhost:8080/api/paste/" + pasteID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&fetched)
	resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, paste["content"], fetched["content"])

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/paste/"+pasteID, nil)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()
}
