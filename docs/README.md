# Guia de Leitura e Entendimento do Projeto — Leitor USBN

Este documento orienta um desenvolvedor (ou você mesmo) sobre como ler, entender e testar rapidamente o projeto.

## Objetivo
- Fornecer um roteiro claro para estudar a arquitetura, componentes e fluxo de dados do repositório.
- Incluir comandos básicos para executar a aplicação CLI e a UI web, além de pontos recomendados para revisão e melhorias.

## Estrutura do repositório (resumida)
- `src/main.go` — runner CLI principal (executa leitura, consulta e persistência).
- `src/web/main.go` — pequeno servidor HTTP que serve UI e API interna.
- `config/` — `config.json` e `config.go` (carregamento e defaults).
- `reader/` — interface `ISBNReader` e implementações (`file_reader.go`, `barcode_reader.go`).
- `processor/` — orquestrador que cria workers e processa ISBNs.
- `api/` — cliente e adaptadores para OpenLibrary (`client.go`, `types.go`).
- `database/` — adapter SQLite com DAO/Repository (`db.go`, `repository.go`, `views.go`).
- `models/` — modelos de domínio (ex.: `book.go`).

## Passos rápidos para explorar o código
1. Abra `config/config.json` para ver valores default (DB, API, leitor).
2. Leia `reader/reader.go` para entender a interface `ISBNReader`.
3. Abra `processor/processor.go` — aqui está o fluxo: recebe ISBNs → consulta API → grava no DB.
4. Verifique `api/client.go` e `api/types.go` para ver como os dados externos são mapeados e convertidos.
5. Veja `database/db.go` e `database/repository.go` para entender schema e operações de persistência.
6. Se quiser a UI, leia `src/web/main.go` e `src/web/templates/books.html`.

## Como executar (exemplos)

Executar a aplicação CLI (processar `config/isbn_list.txt`):
```bash
cd c:\Projects\Go\Leitor_USBN\agent
go run ./src
```

Executar a UI web (porta 8080 por padrão):
```bash
cd c:\Projects\Go\Leitor_USBN\agent
go run ./src/web
# depois acessar http://localhost:8080
```

Executar os pequenos testes / demos (existem dois arquivos em `src/`):
```bash
go run ./src/test_isbn.go
go run ./src/test_etapa4.go
```

Observação: o projeto usa SQLite (`github.com/mattn/go-sqlite3`). Os arquivos DB gerados padrão são `books.db` e `books_etapa4.db`.

## Checklist para entender a arquitetura
- [ ] Identificar as portas (interfaces) e adapters: `ISBNReader`, `api` (OpenLibrary), `database` (SQLite).
- [ ] Mapear o fluxo de dados (leitura → processamento → persistência → UI).
- [ ] Verificar áreas concorrentes: `processor` (workers, canais, mutex) e uso de `context.Context`.
- [ ] Revisar tratamento de erro e logging (consistência e enriquecimento de erros).
- [ ] Checar pontos de configuração (timeouts, delays, número de workers).

## Padrões observados (para revisar)
- Ports & Adapters (hexagonal): `reader` é a porta; `file`/`barcode` são adapters.
- Worker pool / producer-consumer em `processor`.
- Repository/DAO no pacote `database`.

## Sugestões de melhorias e próximos passos
- Tornar o cliente da API e os métodos de DB dependentes de `context.Context` para cancelamento consistente.
- Extrair uma interface para a camada `database` para facilitar testes (mocks) do `processor`.
- Adicionar testes unitários (mocks para `api` e `database`) e testes de integração mínimos.
- Normalizar ISBNs (remoção de hífens, validação de checksum) na camada `reader`.
- Introduzir migrações para o schema (ex.: `golang-migrate`) em vez de DDL inline.

## Pontos de referência (arquivos-chave)
- `config/config.go`, `config/config.json`
- `reader/reader.go`, `reader/file_reader.go`, `reader/barcode_reader.go`
- `processor/processor.go`
- `api/client.go`, `api/types.go`
- `database/db.go`, `database/repository.go`, `database/views.go`
- `src/main.go`, `src/web/main.go`

## Como contribuir / validar mudanças localmente
1. Faça mudanças em uma branch de feature.
2. Rodar `go vet` e `go fmt`.
3. Executar os testes (quando adicionados) e validar manualmente com `go run ./src`.

---
Arquivo criado automaticamente para orientar leitura do projeto.
