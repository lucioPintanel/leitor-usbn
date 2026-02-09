package reader

import (
	"context"
	"fmt"
	"log"
	"time"
)

// BarcodeReaderUSB lê ISBNs de um scanner USB (que funciona como teclado)
type BarcodeReaderUSB struct {
	isbnChan  chan string
	stopChan  chan struct{}
	isRunning bool
	verbose   bool
	timeout   time.Duration
	currentISBN string
}

// NewBarcodeReaderUSB cria uma nova instância do leitor USB
func NewBarcodeReaderUSB(config ReaderConfig) *BarcodeReaderUSB {
	timeout := time.Duration(config.Timeout) * time.Second
	if config.Timeout == 0 {
		timeout = 10 * time.Second // default
	}

	return &BarcodeReaderUSB{
		isbnChan: make(chan string, 100),
		stopChan: make(chan struct{}),
		verbose:  config.Verbose,
		timeout:  timeout,
	}
}

// Start inicia a leitura do scanner USB
// Nota: Isso é uma implementação de exemplo
// Em produção, você precisará usar uma biblioteca para capturar input de dispositivo HID
func (b *BarcodeReaderUSB) Start(ctx context.Context) error {
	if b.isRunning {
		return fmt.Errorf("leitor USB já está ativo")
	}

	b.isRunning = true

	if b.verbose {
		log.Println("=== Leitor USB de Código de Barras Iniciado ===")
		log.Println("Aguardando leitura de barcodes do scanner...")
		log.Println("Scanner funciona como teclado - configure para enviar Enter/Return ao final do código")
	}

	go func() {
		defer func() {
			b.isRunning = false
			close(b.isbnChan)
		}()

		// Para usar em produção, você precisará:
		// 1. Usar github.com/go-vgo/robotgo para capturar input global
		// 2. Ou usar github.com/bendahl/uinput para capturar eventos de dispositivo
		// 3. Ou usar github.com/tmc/keyboardinput para input de teclado

		// Implementação básica usando polling simulado
		timer := time.NewTimer(b.timeout)

		select {
		case <-b.stopChan:
			if b.verbose {
				log.Println("Leitor USB interrompido pelo usuário")
			}
			return
		case <-ctx.Done():
			if b.verbose {
				log.Println("Contexto cancelado")
			}
			return
		case <-timer.C:
			if b.verbose {
				log.Println("Timeout - nenhum código lido no intervalo")
			}
			return
		}
	}()

	return nil
}

// Stop para a leitura do scanner USB
func (b *BarcodeReaderUSB) Stop() error {
	if !b.isRunning {
		return fmt.Errorf("leitor USB não está ativo")
	}

	close(b.stopChan)
	return nil
}

// Read retorna o canal de ISBNs
func (b *BarcodeReaderUSB) Read() <-chan string {
	return b.isbnChan
}

// GetType retorna o tipo do leitor
func (b *BarcodeReaderUSB) GetType() string {
	return "BarcodeReaderUSB"
}

// IsRunning indica se o leitor está ativo
func (b *BarcodeReaderUSB) IsRunning() bool {
	return b.isRunning
}

// SimulateBarcodeScan simula uma leitura de barcode (para testes)
func (b *BarcodeReaderUSB) SimulateBarcodeScan(isbn string) error {
	if !b.isRunning {
		return fmt.Errorf("leitor não está ativo")
	}

	select {
	case b.isbnChan <- isbn:
		if b.verbose {
			log.Printf("Barcode simulado: %s", isbn)
		}
		return nil
	case <-b.stopChan:
		return fmt.Errorf("leitor foi parado")
	}
}

// NotaImplementacao comenta sobre as bibliotecas necessárias
// Para implementação real com captura de input do scanner:
// 1. github.com/go-vgo/robotgo - simula/captura input de teclado global
// 2. github.com/bendahl/uinput - acesso a eventos de dispositivo Linux
// 3. github.com/ravernkoh/keyboardinput - input de teclado simples
// 4. github.com/micmonay/keybd_event - simula eventos de teclado
