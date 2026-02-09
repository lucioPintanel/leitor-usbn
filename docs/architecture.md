# Diagrama de Arquitetura (Mermaid)

```mermaid
flowchart TD
  subgraph CLI[CLI Runner]
    CLI[src/main.go]
  end

  subgraph Readers[Entrada / Readers]
    FileReader[FileISBNReader\n`reader/file_reader.go`]
    USBReader[BarcodeReaderUSB\n`reader/barcode_reader.go`]
  end

  subgraph Core[Core / Orquestração]
    Processor[Processor\n`processor/processor.go`]
    Workers((Worker Pool))
  end

  subgraph Adapters[Adapters Externos]
    APIClient[OpenLibrary API Client\n`api/client.go`]
    Database[SQLite Adapter\n`database/db.go`]
  end

  subgraph UI[UI / Admin]
    WebServer[Web UI + API\n`src/web/main.go`]
    Templates[`src/web/templates/books.html`]
  end

  subgraph Config[Config]
    ConfigFile[`config/config.json`]
  end

  %% Fluxo
  CLI -->|carrega cfg| ConfigFile
  CLI --> Processor
  FileReader --> Processor
  USBReader --> Processor
  Processor --> Workers
  Workers --> APIClient
  Workers --> Database
  Processor --> Database
  APIClient -->|consulta externa| ExternalAPI[(openlibrary.org)]
  WebServer --> Database
  WebServer --> Templates
  CLI --> WebServer

  %% Artefatos / Models
  Models[models/book.go] --- Processor
  Models --- Database

  classDef adapters fill:#f9f,stroke:#333,stroke-width:1px
  class Adapters adapters
```

Legenda rápida:
- `Processor`: orquestra leitura → consulta → persistência; implementa worker pool.
- `Readers`: portas de entrada (file / USB) — implementam a interface `ISBNReader`.
- `APIClient` e `Database`: adapters externos (hexagonal pattern).
- `WebServer`: expõe uma API interna e UI que consome os dados persistidos.

Arquivo: `docs/architecture.md`
