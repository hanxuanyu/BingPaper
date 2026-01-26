package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	err := Init("")
	if err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	cfg := GetConfig()
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}

	if cfg.DB.Type != "sqlite" {
		t.Errorf("Expected DB type sqlite, got %s", cfg.DB.Type)
	}
}
