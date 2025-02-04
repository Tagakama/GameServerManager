package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	yamlContent := `
proxy:
  ip: "127.0.0.1"
  port: "8080"
  code: "admin123"
local_storage:
  directory: "/tmp/"
  filename_exe: "server.exe"
`
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	SetConfigPath(tmpFile.Name())

	cfg := LoadConfig()

	if cfg.Proxy.Ip != "127.0.0.1" {
		t.Errorf("Expected proxy IP 127.0.0.1, got %s", cfg.Proxy.Ip)
	}
	if cfg.Proxy.Port != "8080" {
		t.Errorf("Expected proxy port 8080, got %s", cfg.Proxy.Port)
	}
	if cfg.LocalStorage.Directory != "/tmp/" {
		t.Errorf("Expected directory /tmp/, got %s", cfg.LocalStorage.Directory)
	}
}
