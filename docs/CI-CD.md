# CI/CD — Integração Contínua e Entrega Contínua

Este documento explica como o **CI/CD** funciona neste projeto usando **GitHub Actions**.

## O que é CI/CD?

### CI (Continuous Integration) — Integração Contínua
Automatiza testes e validações a cada push/PR:
- ✅ Executa testes (`go test`)
- ✅ Verifica formatação (`go fmt`)
- ✅ Roda linting (`go vet`, `golangci-lint`)
- ✅ Compila o código (`go build`)
- ✅ Valida cobertura de testes (coverage)

**Benefício**: Detecta problemas rapidamente, antes de mesclar código.

### CD (Continuous Deployment) — Entrega Contínua
(Não implementado ainda neste projeto)

Após CI passar, automaticamente:
- Deploy em servidor de staging
- Deploy em produção
- Gera releases automáticas

## GitHub Actions

**GitHub Actions** é o serviço nativo do GitHub para CI/CD.

### Como funciona

1. Você cria um arquivo `.yaml` em `.github/workflows/`
2. Define **gatilhos** (quando executar):
   - `push` — ao fazer commit
   - `pull_request` — ao abrir/atualizar PR
   - `schedule` — agendado (ex.: diariamente)
3. Define **jobs** (tarefas):
   - Qual SO usar (ubuntu, windows, macos)
   - Quais passos executar
4. GitHub faz a execução em VMs
5. Resultado aparece na aba **Actions** do repo

## Workflow do projeto (`.github/workflows/go.yml`)

### Jobs e etapas

#### **Job 1: test** — Testes, Build, Coverage
Executa em **múltiplas combinações**:
- Sistemas: Ubuntu, Windows, macOS
- Versões Go: 1.21, 1.22

**Passos:**
1. ✅ Checkout do código
2. ✅ Instala Go
3. ✅ Cache de dependências (acelera workflow)
4. ✅ Download de módulos (`go mod download`)
5. ✅ Verifica formatação (`gofmt`)
6. ✅ Linting (`go vet`)
7. ✅ **Testes com race detector** (`go test -race`)
8. ✅ Gera relatório de coverage (`.out`)
9. ✅ **Upload para Codecov** (agregador de coverage)
10. ✅ Compila CLI (`go build -o leitor-usbn ./src`)
11. ✅ Compila Web UI (`go build -o leitor-usbn-web ./src/web`)

#### **Job 2: lint** — Linting avançado
Usa **golangci-lint** (ferramenta profissional):
- Detecta anti-patterns Go
- Verifica código morto
- Identifica complexidade muito alta

### Estratégia de Matrix

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    go-version: ['1.21', '1.22']
```

Isso gera **6 combinações** (3 SOs × 2 versões Go), garantindo compatibilidade.

## Como visualizar resultados

### Na interface GitHub

1. Abra o repositório
2. Clique em **Actions** (barra superior)
3. Escolha um workflow na lista
4. Veja status: ✅ (passou) ou ❌ (falhou)
5. Clique para ver logs detalhados

### No seu PR

- Badge de status aparece automaticamente
- Se falhar, receberá uma anotação

### Localmente (antes de fazer push)

Rode os mesmos comandos:
```bash
go fmt ./...
go vet ./...
go test -v -race ./...
go build -v -o leitor-usbn ./src
go build -v -o leitor-usbn-web ./src/web
```

## Configuração de badges

Você pode adicionar badges de status no `README.md` principal:

```markdown
![Go CI](https://github.com/lucioPintanel/leitor-usbn/actions/workflows/go.yml/badge.svg)
```

## Variáveis e segredos

Para adicionar **segredos** (tokens, senhas) sem expô-los no código:

1. Repo → **Settings** → **Secrets and variables** → **Actions**
2. Clique em **New repository secret**
3. Adicione nome (ex.: `CODECOV_TOKEN`) e valor
4. Use no workflow: `${{ secrets.CODECOV_TOKEN }}`

## Próximos passos sugeridos

- [ ] Habilitar **branch protection rules**:
  - Exigir CI passar antes de mesclar PR
  - Exigir 1+ review
  - Desabilitar push direto em `main`
- [ ] Integrar com **Codecov** para rastrear cobertura ao longo do tempo
- [ ] Adicionar **pre-commit hooks** localment (rodá quick checks antes de push)
- [ ] Configurar notificações (Slack, email) para falhas

## Troubleshooting

### Workflow não dispara
- Verifique o **trigger** (`on: push`, `on: pull_request`)
- Confira se branch está nas `branches:` configuradas

### Falha em formatação
```bash
# Corrigir localmente
go fmt ./...
git add .
git commit -m "style: fix code formatting"
git push
```

### Race detector falhando
- Significa acesso simultâneo a mesma variável sem sincronização
- Use `sync.Mutex`, `sync.RWMutex` ou `chan` (seu projeto já usa!)

### Coverage muito baixa
- Foque em converter `test_*.go` em testes unitários prédios (`*_test.go`)
- Use mocks para `api` e `database`

## Referências

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing](https://golang.org/pkg/testing/)
- [golangci-lint](https://golangci-lint.run/)
- [Codecov for Go](https://docs.codecov.com/docs/golang)

---
Mantendo CI/CD robusto = código confiável! ✅
