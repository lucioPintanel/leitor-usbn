# Sobre .gitignore

O arquivo `.gitignore` instruir o Git a não rastrear certos arquivos e pastas. Isso evita que artefatos indesejados (como binários compilados, bancos de dados locais e arquivos de configuração sensíveis) sejam adicionados ao repositório.

## Estrutura do nosso .gitignore

### Binários e Artefatos Go
```
*.exe
*.exe~
*.dll
*.so
*.so.*
*.dylib
*.test
*.out
vendor/
go.work
```
- Mantém o repositório limpo de executáveis compilados
- Ignora diretório `vendor/` (dependências locais)
- Ignora saída de coverage (`*.out`)

### IDEs e Editores
```
.vscode/
.idea/
*.swp
*.swo
*~
.DS_Store
*.iml
```
- Evita arquivos de configuração pessoal de IDEs (VS Code, IntelliJ, etc.)
- Ignora arquivos temporários de editores

### Banco de Dados Locais
```
books.db
books_etapa4.db
*.db-journal
```
- **Importante**: Banco de dados **não** é versionado
- Cada desenvolvedor gera seu próprio DB local ao executar a app
- Se precisar compartilhar schema, use migrations (veja `docs/ASSISTANT.md`)

### Configuração Local
```
config/config.local.json
```
- Para valores de configuração específicos da máquina
- Exemplo: porta diferente, caminho local para ISBN list, etc.
- Criada localmente, nunca adicionada ao git

### Build e Logs
```
dist/
build/
*.log
tmp/
temp/
```
- Diretórios de saída de build
- Logs da aplicação
- Arquivos temporários

## Como usar

### Verificar se um arquivo será ignorado
```bash
git check-ignore -v <arquivo>
```

### Forçar adicionar um arquivo ignorado (use com cuidado!)
```bash
git add -f <arquivo>
```

### Remover arquivo já versionado, mas agora ignorado
```bash
git rm --cached <arquivo>
git commit -m "chore: remove tracked file that should be ignored"
```

## Bom prática

1. **Nunca versione**: senhas, tokens, chaves API, arquivos DB locais
2. **Sempre versione**: código-fonte, testes, documentação, schema SQL
3. **Quando em dúvida**: adicione ao `.gitignore` e providencie um exemplo/template

---
Para mais informações sobre `.gitignore`, veja: https://git-scm.com/docs/gitignore
