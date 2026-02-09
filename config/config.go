package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Config contém todas as configurações da aplicação
type Config struct {
	Database DatabaseConfig `json:"database"`
	API      APIConfig      `json:"api"`
	Reader   ReaderConfig   `json:"reader"`
	Processor ProcessorConfig `json:"processor"`
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

// APIConfig configurações da API
type APIConfig struct {
	Provider string `json:"provider"`
	BaseURL  string `json:"baseUrl"`
	Timeout  int    `json:"timeout"`
}

// ReaderConfig configurações do leitor
type ReaderConfig struct {
	InputFile string `json:"inputFile"`
	Type      string `json:"type"`
	Verbose   bool   `json:"verbose"`
}

// ProcessorConfig configurações do processador
type ProcessorConfig struct {
	MaxWorkers           int `json:"maxWorkers"`
	DelayBetweenRequests int `json:"delayBetweenRequests"`
	MaxRetries           int `json:"maxRetries"`
	Verbose              bool `json:"verbose"`
}

// LoadConfig carrega configurações de um arquivo JSON
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de configuração: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}

	// Validações
	if config.Database.Type == "" {
		config.Database.Type = "sqlite"
	}
	if config.Database.Path == "" {
		config.Database.Path = "./books.db"
	}
	if config.API.BaseURL == "" {
		config.API.BaseURL = "https://openlibrary.org/api/books"
	}
	if config.API.Timeout == 0 {
		config.API.Timeout = 10
	}
	if config.Reader.Type == "" {
		config.Reader.Type = "file"
	}
	if config.Reader.InputFile == "" {
		config.Reader.InputFile = "./isbn_list.txt"
	}
	if config.Processor.MaxWorkers == 0 {
		config.Processor.MaxWorkers = 1
	}
	if config.Processor.DelayBetweenRequests == 0 {
		config.Processor.DelayBetweenRequests = 500
	}
	if config.Processor.MaxRetries == 0 {
		config.Processor.MaxRetries = 3
	}

	return &config, nil
}
