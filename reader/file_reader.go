package reader

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

// FileISBNReader lê ISBNs de um arquivo de texto
type FileISBNReader struct {
	filePath  string
	isbnChan  chan string
	stopChan  chan struct{}
	isRunning bool
	verbose   bool
}

// NewFileISBNReader cria uma nova instância do leitor de arquivo
func NewFileISBNReader(config ReaderConfig) *FileISBNReader {
	return &FileISBNReader{
		filePath: config.FilePath,
		isbnChan: make(chan string, 100), // buffer para evitar bloqueios
		stopChan: make(chan struct{}),
		verbose:  config.Verbose,
	}
}

// Start inicia a leitura do arquivo
func (f *FileISBNReader) Start(ctx context.Context) error {
	if f.isRunning {
		return fmt.Errorf("leitor de arquivo já está ativo")
	}

	f.isRunning = true

	go func() {
		defer func() {
			f.isRunning = false
			close(f.isbnChan)
		}()

		file, err := os.Open(f.filePath)
		if err != nil {
			log.Printf("Erro ao abrir arquivo: %v", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 0

		for scanner.Scan() {
			select {
			case <-f.stopChan:
				if f.verbose {
					log.Println("Leitura de arquivo interrompida pelo usuário")
				}
				return
			case <-ctx.Done():
				if f.verbose {
					log.Println("Contexto cancelado")
				}
				return
			default:
			}

			lineNumber++
			line := strings.TrimSpace(scanner.Text())

			// Ignorar linhas vazias e comentários
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// Validação básica de ISBN
			if len(line) < 10 {
				if f.verbose {
					log.Printf("Linha %d: ISBN inválido (muito curto): %s", lineNumber, line)
				}
				continue
			}

			if f.verbose {
				log.Printf("Linha %d: ISBN lido: %s", lineNumber, line)
			}

			// Enviar ISBN pelo canal
			select {
			case f.isbnChan <- line:
			case <-f.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Erro ao ler arquivo: %v", err)
			return
		}

		if f.verbose {
			log.Printf("Total de %d linhas lidas do arquivo", lineNumber)
		}
	}()

	return nil
}

// Stop para a leitura do arquivo
func (f *FileISBNReader) Stop() error {
	if !f.isRunning {
		return fmt.Errorf("leitor de arquivo não está ativo")
	}

	close(f.stopChan)
	return nil
}

// Read retorna o canal de ISBNs
func (f *FileISBNReader) Read() <-chan string {
	return f.isbnChan
}

// GetType retorna o tipo do leitor
func (f *FileISBNReader) GetType() string {
	return "FileISBNReader"
}

// IsRunning indica se o leitor está ativo
func (f *FileISBNReader) IsRunning() bool {
	return f.isRunning
}
