# Guia de Contribuição — Leitor USBN

Obrigado por considerar contribuir com melhorias neste projeto! Este documento fornece diretrizes para manter a qualidade e consistência do código.

## Como começar

1. **Fork** o repositório em https://github.com/lucioPintanel/leitor-usbn
2. **Clone** seu fork localmente
3. Crie uma **branch de feature** a partir de `main`:
   ```bash
   git checkout -b feature/sua-feature
   ```
4. Faça suas alterações
5. **Teste** localmente (veja seção abaixo)
6. **Commit** seguindo as convenções (veja seção abaixo)
7. **Push** para seu fork e abra um **Pull Request**

## Padrão de Commit

Usamos **Conventional Commits** para facilitar a leitura do histórico.

### Formato
```
<tipo>(<escopo>): <assunto>

<corpo — opcional>

<footer — opcional>
```

### Tipos
- `feat`: Nova feature
- `fix`: Correção de bug
- `docs`: Mudança em documentação
- `refactor`: Refatoração (sem mudar comportamento)
- `test`: Adição ou modificação de testes
- `chore`: Alterações em build, deps, configs
- `perf`: Melhoria de performance

### Exemplos
```
feat(processor): add retry logic with exponential backoff

Implemented a new ProcessResult type to track ISBN processing status.
Added automatic retry with exponential backoff for API failures.

Closes #12
```

```
fix(reader): validate ISBN length before processing
```

```
docs: add architecture diagram and contribution guide
```

## Como testar localmente

### Rodar a aplicação CLI
```bash
go run ./src
```

### Rodar a UI web
```bash
go run ./src/web
# depois acesess http://localhost:8080
```

### Executar testes
```bash
go test ./...
```

### Verificar formatação e linting
```bash
go fmt ./...
go vet ./...
```

## Padrões de código

### Go Style Guide
- Use `gofmt` para formatar código
- Mantenha funções pequenas e focadas
- Adicione comentários para pacotes públicos e funções exportadas
- Use nomes descritivos para variáveis

### Estrutura de packages
- `api/` — cliente OpenLibrary e conversão de dados
- `database/` — operações SQLite (DAO/Repository)
- `reader/` — leitura de ISBNs (File, USB)
- `processor/` — orquestração e worker pool
- `config/` — configuração centralizada
- `models/` — tipos/estruturas de domínio
- `src/` — executáveis (main.go, web/main.go)

## Processos de PR

1. **Abrir PR**: descreva o problema/feature que está resolvendo
2. **Testes**: certifique-se de que `go test ./...` passa
3. **Documentação**: atualize `docs/` se a mudança afeta arquitetura ou uso
4. **Review**: aguarde feedback; responda comentários
5. **Merge**: merge via squash para manter histórico limpo

## Melhorias sugeridas (roadmap)

- [ ] Extrair interface para `database.Database` (facilita mocking)
- [ ] Adaptar `api.BookAPIClient` para aceitar `context.Context`
- [ ] Adicionar testes unitários com mocks
- [ ] Normalização e validação de ISBN (checksum, formato)
- [ ] Pipeline de migrations para schema (golang-migrate)
- [ ] CI/CD (GitHub Actions): go vet, go fmt, testes
- [ ] Linting (golangci-lint)
- [ ] Suporte a logging estruturado (logrus ou zerolog)

## Dúvidas?

- Abra uma [Issue](https://github.com/lucioPintanel/leitor-usbn/issues)
- Consulte [docs/ASSISTANT.md](docs/ASSISTANT.md) para saber mais sobre o papel do engenheiro/mentor

---
Obrigado por contribuir! 🎉
