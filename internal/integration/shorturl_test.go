package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNonexistentShortURL(t *testing.T) {
	skipIfNotIntegration(t)

	resp, err := http.Get("http://localhost:8080/api/shorturl/fake123")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}
