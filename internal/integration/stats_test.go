package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllStats(t *testing.T) {
	skipIfNotIntegration(t)

	resp, err := http.Get("http://localhost:8080/api/stats")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var stats []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&stats)
	resp.Body.Close()
}
