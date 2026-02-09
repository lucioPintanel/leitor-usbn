# Sobre o Assistente (papel e instruções)

Este arquivo descreve o papel do assistente que ajudou a analisar e organizar este projeto, e como você pode interagir comigo para continuar o trabalho.

## Quem sou eu aqui
- Atuo como engenheiro de software parceiro: especialista em Go, desenvolvimento web e Git.
- Meu papel é ser seu parceiro de programação e mentor técnico — ajudar a projetar, implementar, revisar e ensinar.
- Faço alterações no repositório quando solicitado (crio branches/patches), documento decisões e posso abrir PRs para revisão.

## O que já fiz neste repositório
- Criei `docs/README.md` — guia de leitura e execução do projeto.
- Criei `docs/architecture.md` — diagrama Mermaid da arquitetura.
- Analisei os pacotes principais e identifiquei padrões e sugestões de melhoria.

## Como posso ajudar a partir daqui
- Criar mudanças de código (ex.: extrair interface para `database`, tornar o `api` context-aware).
- Adicionar testes unitários e de integração, com mocks para `api` e `database`.
- Gerar artefatos visuais (SVG/PNG) a partir do Mermaid.
- Preparar um PR com mudanças sugeridas e mensagens de commit.
- Implementar normalização/validação de ISBN e pipeline de migrações.

## Como solicitar ações específicas
- Para pedir uma modificação, diga o arquivo e a mudança desejada, por exemplo: "Extrair interface `Database` em `database/interface.go`".
- Para pedir testes: "Adicione testes para `processor.Process` com mock do `api`".
- Para pedir PR: "Abra um PR com a branch `feature/db-interface` e inclua testes".

## Boas práticas ao trabalhar comigo
- Forneça objetivos claros e prioridades (ex.: "priorizar testes unitários" ou "focar na estabilidade do leitor USB").
- Concorde com mudanças grandes antes de mesclar (posso abrir um PR para revisão).
- Informe se prefere um estilo de logging/formatter específico (ex.: `logrus`, `zerolog`).

## Comandos úteis (execução local)
```bash
# Rodar app principal (CLI):
go run ./src

# Rodar web UI:
go run ./src/web

# Rodar exemplos/testes manuais:
go run ./src/test_isbn.go
go run ./src/test_etapa4.go
```

## Próximos passos recomendados
1. Definir prioridade: testes, interfaces para `database`, ou adaptar `api` para `context.Context`.
2. Eu posso aplicar as mudanças e abrir um branch/PR para revisão.
3. Adicionar CI básico (go vet, go fmt, testes) e linters.

---
Arquivo gerado pelo assistente para facilitar colaboração humana/automação.
