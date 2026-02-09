# Leitor USBN - Sistema de Leitura e Consulta de Livros

Um sistema em Go para ler ISBNs de livros (via arquivo ou scanner USB), consultar dados em uma API pública e armazenar em um banco de dados local SQLite com boas práticas de modelagem relacional.

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Arquitetura](#arquitetura)
- [Requisitos](#requisitos)
- [Instalação](#instalação)
- [Configuração](#configuração)
- [Uso](#uso)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Componentes](#componentes)
- [Banco de Dados](#banco-de-dados)
- [Exemplos](#exemplos)
- [Próximas Melhorias](#próximas-melhorias)

## 🎯 Visão Geral

O projeto **Leitor USBN** automatiza o processo de:

1. **Leitura de ISBNs** - Via arquivo `.txt` ou scanner de código de barras USB
2. **Consulta de API** - Integração com OpenLibrary para obter dados dos livros
3. **Armazenamento** - Banco de dados SQLite com normalização adequada
4. **Processamento Paralelo** - Múltiplos workers para processar ISBNs simultaneamente

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────┐
│              LEITOR USBN                        │
├─────────────────────────────────────────────────┤
│                                                 │
│  Config (JSON)                                  │
│      ↓                                          │
│  ┌─────────────┐                               │
│  │ ISBNReader  │ ← FileISBNReader              │
│  │ (Interface) │ ← BarcodeReaderUSB            │
│  └──────┬──────┘                               │
│         ↓                                      │
│  ┌──────────────┐                              │
│  │  Processor   │ (orquestrador)               │
│  │  (Workers)   │                              │
│  └──────┬───────┘                              │
│         ↓                                      │
│  ┌──────────────────────┐                      │
│  │   BookAPIClient      │ (OpenLibrary API)    │
│  └──────┬───────────────┘                      │
│         ↓                                      │
│  ┌──────────────────────┐                      │
│  │   Database           │ (SQLite)             │
│  │   (Repository)       │                      │
│  └──────────────────────┘                      │
└─────────────────────────────────────────────────┘
```

## 📦 Requisitos

- **Go** 1.21 ou superior
- **SQLite3** (automaticamente incluído via pacote)
- **Internet** para consultar OpenLibrary API

## 🚀 Instalação

### 1. Clonar o repositório

```bash
cd c:\Projects\Go\Leitor_USBN\agent
```

### 2. Instalar dependências

```bash
go mod download
```

### 3. Verificar instalação

```bash
go run ./src/main.go -h
```

## ⚙️ Configuração

### Arquivo config.json

Localizável em `config/config.json`:

```json
{
  "database": {
    "type": "sqlite",
    "path": "./books.db"
  },
  "api": {
    "provider": "openlibrary",
    "baseUrl": "https://openlibrary.org/api/books",
    "timeout": 10
  },
  "reader": {
    "inputFile": "./config/isbn_list.txt",
    "type": "file",
    "verbose": true
  },
  "processor": {
    "maxWorkers": 1,
    "delayBetweenRequests": 500,
    "maxRetries": 3,
    "verbose": true
  }
}
```

### Parâmetros

#### Database
- `type`: Tipo de banco (apenas "sqlite" por enquanto)
- `path`: Caminho para o arquivo SQLite

#### API
- `provider`: Provedor da API (apenas "openlibrary")
- `baseUrl`: URL base da API
- `timeout`: Timeout em segundos para requisições

#### Reader
- `inputFile`: Caminho para arquivo de ISBNs
- `type`: Tipo de leitor ("file" ou "barcode")
- `verbose`: Ativa logs detalhados

#### Processor
- `maxWorkers`: Número de workers paralelos (recomendado: 1-4)
- `delayBetweenRequests`: Delay entre requisições em ms
- `maxRetries`: Número de tentativas por ISBN
- `verbose`: Ativa logs detalhados

## 📖 Uso

### Executar com configuração padrão

```bash
go run ./src/main.go
```

### Executar com configuração customizada

```bash
go run ./src/main.go -config ./config/seu_config.json
```

### Compilar para executável

```bash
go build -o leitor_usbn.exe ./src/main.go
.\leitor_usbn.exe
```

### Criar lista de ISBNs

Crie um arquivo `config/isbn_list.txt`:

```txt
# Comentários começam com #
0132350882
978-0201633610
9780201633610
0596007124
```

## 📁 Estrutura do Projeto

```
agent/
├── config/
│   ├── config.json           # Configurações da aplicação
│   ├── config.go             # Parser de configurações
│   └── isbn_list.txt         # Lista de ISBNs para processar
├── database/
│   ├── db.go                 # Inicialização do SQLite
│   └── repository.go         # Operações CRUD
├── api/
│   ├── client.go             # Cliente HTTP para OpenLibrary
│   └── types.go              # Estruturas de dados da API
├── reader/
│   ├── reader.go             # Interface ISBNReader
│   ├── file_reader.go        # Leitor de arquivo
│   └── barcode_reader.go     # Leitor de scanner USB
├── processor/
│   └── processor.go          # Orquestrador principal
├── models/
│   └── book.go               # Modelo de dados
├── src/
│   ├── main.go               # Programa principal
│   ├── test_isbn.go          # Teste de ISBN único
│   └── test_etapa4.go        # Teste de leitura de arquivo
├── go.mod                    # Dependências do projeto
└── README.md                 # Este arquivo
```

## 🔧 Componentes

### 1. **Reader** (`reader/`)

Define como os ISBNs são lidos:

- `ISBNReader` (interface)
  - `FileISBNReader` - Lê de arquivo .txt
  - `BarcodeReaderUSB` - Integração com scanner USB

**Uso:**
```go
config := reader.ReaderConfig{FilePath: "./isbn_list.txt"}
reader := reader.NewFileISBNReader(config)
err := reader.Start(ctx)
for isbn := range reader.Read() {
    // Processar ISBN
}
```

### 2. **API Client** (`api/`)

Integração com OpenLibrary:

- `BookAPIClient` - Cliente HTTP
- `OpenLibraryResponse` - Estrutura da resposta
- `BookData` - Dados normalizados

**Uso:**
```go
client := api.NewBookAPIClient("https://openlibrary.org/api/books", 10)
book, err := client.GetBookByISBN("0132350882")
data := api.ConvertToBookData(book)
```

### 3. **Database** (`database/`)

Gerenciamento do SQLite com boas práticas:

- `NewDatabase()` - Cria conexão
- `InitSchema()` - Cria tabelas
- `GetOrCreateAuthor()` - Evita duplicatas
- `GetOrCreatePublisher()` - Evita duplicatas
- `SaveBook()` - Insere/atualiza livro

**Uso:**
```go
db, _ := database.NewDatabase("./books.db")
db.InitSchema()
author, _ := db.GetOrCreateAuthor("Robert C. Martin")
book := &database.Book{
    ISBN: "0132350882",
    Title: "Clean Code",
    AuthorID: &author.ID,
}
db.SaveBook(book)
```

### 4. **Processor** (`processor/`)

Orquestra o fluxo completo:

- Coordena múltiplos workers
- Gerencia retry automático
- Coleta estatísticas
- Thread-safe

**Uso:**
```go
config := processor.ProcessorConfig{
    MaxWorkers: 2,
    DelayBetweenRequests: 500 * time.Millisecond,
}
proc := processor.NewProcessor(db, apiClient, reader, config)
proc.Process(ctx)
proc.PrintSummary()
```

### 5. **Config** (`config/`)

Gerenciador de configurações em JSON:

- Carrega arquivo JSON
- Valores padrão automáticos
- Validações

**Uso:**
```go
cfg, _ := config.LoadConfig("./config/config.json")
fmt.Println(cfg.Database.Path)
```

## 💾 Banco de Dados

### Schema

O banco utiliza 3 tabelas normalizadas:

#### Tabela: `authors`
```sql
CREATE TABLE authors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Tabela: `publishers`
```sql
CREATE TABLE publishers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Tabela: `books`
```sql
CREATE TABLE books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    isbn TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    author_id INTEGER,
    publisher_id INTEGER,
    publish_date TEXT,
    pages INTEGER,
    description TEXT,
    cover_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES authors(id),
    FOREIGN KEY (publisher_id) REFERENCES publishers(id)
);
```

### Índices

Criados automaticamente para otimizar buscas:
- `idx_books_isbn`
- `idx_books_author_id`
- `idx_books_publisher_id`
- `idx_authors_name`
- `idx_publishers_name`

### Boas Práticas Implementadas

✅ **Normalização** - Tabelas separadas para autores e editoras  
✅ **Foreign Keys** - Integridade referencial  
✅ **Unique Constraints** - Evita duplicatas  
✅ **Índices** - Performance otimizada  
✅ **Timestamps** - Auditoria (created_at, updated_at)  
✅ **NULL Safety** - Campos opcionais com ponteiros  

## 📚 Exemplos

### Exemplo 1: Processar um ISBN único

```bash
go run ./src/test_isbn.go
```

Testa a integração completa com um ISBN fixo (Clean Code).

### Exemplo 2: Processar arquivo de ISBNs

```bash
go run ./src/test_etapa4.go
```

Processa todos os ISBNs do arquivo configurado.

### Exemplo 3: Executar aplicação completa

```bash
go run ./src/main.go
```

Executa a aplicação com todas as configurações.

### Exemplo 4: Múltiplos workers

Edite `config/config.json`:

```json
{
  "processor": {
    "maxWorkers": 3
  }
}
```

```bash
go run ./src/main.go -config ./config/config.json
```

## 🔄 Fluxo de Execução

```
1. Carregar configurações (config.json)
2. Inicializar banco de dados (SQLite)
3. Criar schema se não existir
4. Inicializar cliente API (OpenLibrary)
5. Criar leitor de ISBNs (arquivo ou scanner)
6. Iniciar processamento
7. Para cada ISBN:
   a. Consultar API
   b. Converter dados
   c. Obter/criar autor
   d. Obter/criar editora
   e. Salvar no banco
   f. Coletar resultado
8. Exibir resumo
9. Finalizar
```

## 🚀 Próximas Melhorias

### Curto Prazo
- [ ] Implementar scanner USB real com biblioteca keybd_event
- [ ] Adicionar validação de ISBN (algoritmo Luhn)
- [ ] Exportar dados para CSV/Excel
- [ ] Interface CLI melhorada (colors, progress bar)

### Médio Prazo
- [ ] Dashboard web (Go + HTML/CSS/JS)
- [ ] API REST para consultar dados
- [ ] Suporte para múltiplas APIs (Google Books, etc)
- [ ] Cache de requisições

### Longo Prazo
- [ ] Autenticação de usuários
- [ ] Histórico de importações
- [ ] Notificações (email, Slack)
- [ ] Integração com sistemas de biblioteca
- [ ] Containerização (Docker)
- [ ] CI/CD pipeline

## 🐛 Troubleshooting

### Erro: "não é possível abrir arquivo de configuração"
**Solução:** Verifique o caminho do arquivo de configuração

```bash
go run ./src/main.go -config ./config/config.json
```

### Erro: "ISBN não encontrado na API"
**Solução:** O ISBN pode estar inválido ou não existir no OpenLibrary. Verifique em https://openlibrary.org/

### API timeout
**Solução:** Aumente o timeout em `config.json`:

```json
{
  "api": {
    "timeout": 20
  }
}
```

### Erro de acesso ao banco de dados
**Solução:** Verifique permissões de leitura/escrita no diretório:

```bash
# Windows
icacls . /grant %USERNAME%:F

# Linux/Mac
chmod 755 ./
```

## 📞 Suporte

Para dúvidas ou problemas:

1. Verifique logs em modo verbose (`verbose: true`)
2. Consulte a documentação do OpenLibrary: https://openlibrary.org/developers/api
3. Revise a estrutura do projeto em [Estrutura do Projeto](#estrutura-do-projeto)

## 📄 Licença

Este projeto é livre para uso e modificação.

---

**Versão:** 1.0  
**Data:** Fevereiro 2026  
**Linguagem:** Go 1.21+
