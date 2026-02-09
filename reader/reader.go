package reader

import (
	"context"
)

// ISBNReader define a interface para diferentes formas de leitura de ISBN
type ISBNReader interface {
	// Start inicia o leitor
	Start(ctx context.Context) error

	// Stop para o leitor
	Stop() error

	// Read retorna um canal com os ISBNs lidos
	Read() <-chan string

	// GetType retorna o tipo do leitor
	GetType() string

	// IsRunning indica se o leitor está ativo
	IsRunning() bool
}

// ReaderConfig contém configurações para os leitores
type ReaderConfig struct {
	// Para FileISBNReader
	FilePath string

	// Para BarcodeReaderUSB
	DevicePath string
	Timeout    int // em segundos

	// Geral
	Verbose bool
}
