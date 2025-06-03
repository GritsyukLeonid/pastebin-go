package integration

import (
	"os"
	"testing"
)

func skipIfNotIntegration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION") != "1" {
		t.Skip("Пропущен интеграционный тест")
	}
}
