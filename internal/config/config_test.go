package config

import "testing"

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()
	if len(cfg.DSN) == 0 {
		t.Errorf("Error DSN config %v", cfg.DSN)
	}

	if len(cfg.ServerAddress) == 0 {
		t.Errorf("Error ServerAddress config %v", cfg.ServerAddress)
	}

	if len(cfg.FileServerAddress) == 0 {
		t.Errorf("Error ServerAddress config %v", cfg.FileServerAddress)
	}

}
