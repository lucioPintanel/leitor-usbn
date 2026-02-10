package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig testa carregamento de configuração válida
func TestLoadConfig(t *testing.T) {
	// Usar arquivo de config existente
	configPath := "./config.json"

	// Verificar se arquivo existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("config.json não encontrado em %s", configPath)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig returned nil")
	}

	// Verificar defaults foram aplicados
	if cfg.Database.Type == "" {
		t.Error("Database.Type é obrigatório")
	}
	if cfg.API.BaseURL == "" {
		t.Error("API.BaseURL é obrigatório")
	}
	if cfg.Processor.MaxWorkers == 0 {
		t.Error("Processor.MaxWorkers deve ser > 0")
	}
}

// TestLoadConfigNotFound testa erro ao carregar arquivo inexistente
func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("LoadConfig should error for nonexistent file")
	}
}

// TestConfigDefaults testa aplicação de valores default
func TestConfigDefaults(t *testing.T) {
	// Criar arquivo temporário com config vazia/mínima
	tmpFile := filepath.Join(t.TempDir(), "config.json")
	
	// Escrever um JSON mínimo
	err := os.WriteFile(tmpFile, []byte(`{}`), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	// Verificar que defaults foram aplicados
	if cfg.Database.Type != "sqlite" {
		t.Errorf("Database.Type default = %s, want sqlite", cfg.Database.Type)
	}
	if cfg.Database.Path != "./books.db" {
		t.Errorf("Database.Path default = %s, want ./books.db", cfg.Database.Path)
	}
	if cfg.API.Timeout != 10 {
		t.Errorf("API.Timeout default = %d, want 10", cfg.API.Timeout)
	}
	if cfg.Processor.MaxWorkers != 1 {
		t.Errorf("Processor.MaxWorkers default = %d, want 1", cfg.Processor.MaxWorkers)
	}
	if cfg.Processor.DelayBetweenRequests != 500 {
		t.Errorf("Processor.DelayBetweenRequests default = %d, want 500", cfg.Processor.DelayBetweenRequests)
	}
	if cfg.Processor.MaxRetries != 3 {
		t.Errorf("Processor.MaxRetries default = %d, want 3", cfg.Processor.MaxRetries)
	}
}

// TestConfigPresetsValues testa que valores fornecidos não são sobrescritom
func TestConfigPreserveValues(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")
	
	configJSON := `{
		"database": {
			"type": "postgres",
			"path": "/custom/path/db.sql"
		},
		"api": {
			"timeout": 30
		},
		"processor": {
			"maxWorkers": 10,
			"maxRetries": 5
		}
	}`
	
	err := os.WriteFile(tmpFile, []byte(configJSON), 0644)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.Database.Type != "postgres" {
		t.Errorf("Database.Type = %s, want postgres", cfg.Database.Type)
	}
	if cfg.Database.Path != "/custom/path/db.sql" {
		t.Errorf("Database.Path = %s, want /custom/path/db.sql", cfg.Database.Path)
	}
	if cfg.API.Timeout != 30 {
		t.Errorf("API.Timeout = %d, want 30", cfg.API.Timeout)
	}
	if cfg.Processor.MaxWorkers != 10 {
		t.Errorf("Processor.MaxWorkers = %d, want 10", cfg.Processor.MaxWorkers)
	}
	if cfg.Processor.MaxRetries != 5 {
		t.Errorf("Processor.MaxRetries = %d, want 5", cfg.Processor.MaxRetries)
	}
}
