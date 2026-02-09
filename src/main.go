package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leitor-usbn/api"
	"leitor-usbn/config"
	"leitor-usbn/database"
	"leitor-usbn/processor"
	"leitor-usbn/reader"
)

func main() {
	// Flags
	configPath := flag.String("config", "./config/config.json", "Caminho para arquivo de configuração")
	flag.Parse()

	fmt.Println("=== LEITOR USBN - Sistema de Leitura e Consulta de Livros ===\n")

	// Carregar configurações
	fmt.Println("[1] Carregando configurações...")
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}
	fmt.Printf("✓ Configurações carregadas de: %s\n\n", *configPath)

	// Inicializar banco de dados
	fmt.Println("[2] Inicializando banco de dados...")
	db, err := database.NewDatabase(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	err = db.InitSchema()
	if err != nil {
		log.Fatalf("Erro ao criar schema: %v", err)
	}
	fmt.Printf("✓ Banco de dados inicializado: %s\n\n", cfg.Database.Path)

	// Inicializar cliente API
	fmt.Println("[3] Inicializando cliente API...")
	apiClient := api.NewBookAPIClient(cfg.API.BaseURL, cfg.API.Timeout)
	fmt.Printf("✓ Cliente API criado: %s\n\n", cfg.API.Provider)

	// Criar leitor de ISBNs
	fmt.Println("[4] Configurando leitor de ISBNs...")
	readerConfig := reader.ReaderConfig{
		FilePath: cfg.Reader.InputFile,
		Verbose:  cfg.Reader.Verbose,
	}

	var isbnReader reader.ISBNReader
	switch cfg.Reader.Type {
	case "file":
		isbnReader = reader.NewFileISBNReader(readerConfig)
	case "barcode":
		isbnReader = reader.NewBarcodeReaderUSB(readerConfig)
	default:
		log.Fatalf("Tipo de leitor desconhecido: %s", cfg.Reader.Type)
	}

	fmt.Printf("✓ Leitor de ISBNs configurado: %s\n\n", isbnReader.GetType())

	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Capturar sinais para graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\n\nSinal recebido: %v\n", sig)
		cancel()
	}()

	// Iniciar leitor
	fmt.Println("[5] Iniciando leitura de ISBNs...")
	err = isbnReader.Start(ctx)
	if err != nil {
		log.Fatalf("Erro ao iniciar leitor: %v", err)
	}
	fmt.Println("✓ Leitor iniciado\n")

	// Criar processador
	fmt.Println("[6] Configurando processador...")
	processorConfig := processor.ProcessorConfig{
		MaxWorkers:           cfg.Processor.MaxWorkers,
		DelayBetweenRequests: time.Duration(cfg.Processor.DelayBetweenRequests) * time.Millisecond,
		MaxRetries:           cfg.Processor.MaxRetries,
		Verbose:              cfg.Processor.Verbose,
	}

	proc := processor.NewProcessor(db, apiClient, isbnReader, processorConfig)
	fmt.Printf("✓ Processador criado com %d worker(s)\n\n", processorConfig.MaxWorkers)

	// Processar ISBNs
	fmt.Println("[7] Processando ISBNs...")
	fmt.Println("==================================================")

	startTime := time.Now()
	err = proc.Process(ctx)
	if err != nil {
		log.Fatalf("Erro ao processar: %v", err)
	}

	elapsed := time.Since(startTime)

	// Imprimir resumo
	proc.PrintSummary()

	// Total de livros no banco
	totalBooks, err := db.CountBooks()
	if err == nil {
		fmt.Printf("\nTotal de livros no banco de dados: %d\n", totalBooks)
	}

	fmt.Printf("Tempo total: %v\n", elapsed)
	fmt.Println("\n✓ Aplicação finalizada com sucesso!")
}

