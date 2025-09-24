# Release Guide

Este guia explica como criar releases automatizadas no projeto YRN.

## 🏷️ Como Criar uma Release

### 1. Preparação
Certifique-se de que:
- Todos os testes estão passando
- O código está na branch `main`
- As mudanças estão documentadas

### 2. Criar e Pushar Tag
```bash
# Exemplo para versão 1.0.0
git tag v1.0.0
git push origin v1.0.0
```

### 3. Automação
O workflow **unificado** irá automaticamente:
- ✅ Verificar se imagem Docker já existe (do CI)
- ✅ Buildar Docker apenas se necessário
- ✅ Buildar binários para Linux, macOS e Windows
- ✅ Criar arquivos compactados (.tar.gz, .zip)
- ✅ Gerar changelog baseado nos commits
- ✅ Criar release no GitHub com assets

## 📋 Versionamento Semântico

Siga o padrão [Semantic Versioning](https://semver.org/):

### Formato: `v{MAJOR}.{MINOR}.{PATCH}`

- **MAJOR** (`v2.0.0`): Breaking changes
- **MINOR** (`v1.1.0`): Novas funcionalidades (compatível)
- **PATCH** (`v1.0.1`): Bug fixes

### Exemplos:
```bash
# Nova funcionalidade
git tag v1.1.0

# Bug fix
git tag v1.0.1

# Breaking change
git tag v2.0.0

# Pre-release (marcada como prerelease)
git tag v1.1.0-beta.1
git tag v1.1.0-rc.1
```

## 🚀 O que é Gerado

### Binários
- `yrn-agent-{os}-{arch}`
- `yrn-connector-{os}-{arch}`
- `yrn-api-{os}-{arch}`

### Plataformas Suportadas
- Linux AMD64
- macOS AMD64
- Windows AMD64

### Arquivos de Release
- `yrn-v{version}-linux-amd64.tar.gz`
- `yrn-v{version}-darwin-amd64.tar.gz`
- `yrn-v{version}-windows-amd64.zip`

### Imagens Docker
- `ghcr.io/yrn-go/yrn:{version}`
- `ghcr.io/yrn-go/yrn:latest`

## 🔧 Troubleshooting

### Erro: "Tag already exists"
```bash
# Remove tag local e remota
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Recrie a tag
git tag v1.0.0
git push origin v1.0.0
```

### Erro: "Workflow not triggered"
- Verifique se a tag segue o padrão `v*.*.*`
- Confirme que o workflow está na branch `main`
- Verifique permissões do repositório

### Re-executar Workflow
1. Vá para Actions no GitHub
2. Encontre o workflow "Release"
3. Clique em "Re-run jobs"

## 📝 Changelog Automático

O changelog é gerado automaticamente baseado nos commits desde a última tag:

### Formato dos Commits (Recomendado)
```
feat: adiciona nova funcionalidade X
fix: corrige bug no componente Y
docs: atualiza documentação Z
chore: atualiza dependências
```

### Exemplo de Changelog Gerado:
```markdown
## What's Changed
- feat: adiciona plugin HTTP com retry (a1b2c3d)
- fix: corrige validação de schema JSON (e4f5g6h)
- docs: atualiza README com novos exemplos (i7j8k9l)

## 🚀 Installation
### Docker
\`\`\`bash
docker pull ghcr.io/yrn-go/yrn:v1.0.0
\`\`\`
```

## 🎯 Próximos Passos

Após criar a release:

1. **Teste a release**: Baixe e teste os binários
2. **Anuncie**: Comunique a nova versão
3. **Documente**: Atualize documentação se necessário
4. **Monitor**: Acompanhe issues relacionadas à nova versão

## 🔍 Verificações Pós-Release

- [ ] Release aparece na página Releases
- [ ] Binários estão funcionando
- [ ] Imagem Docker foi publicada
- [ ] Deployment foi atualizado (se aplicável)
- [ ] Changelog está correto
- [ ] Links estão funcionando