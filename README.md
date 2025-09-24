# YRN - Yet another Routine Network

[![Go Tests & Docker Build](https://github.com/yrn-go/yrn/actions/workflows/test-and-build.yaml/badge.svg)](https://github.com/yrn-go/yrn/actions/workflows/test-and-build.yaml)
[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://hub.docker.com)

YRN é uma plataforma de orquestração de serviços distribuída construída em Go, projetada para gerenciar fluxos de trabalho baseados em plugins com descoberta dinâmica de serviços e validação de esquemas JSON.

## 📋 Índice

- [Visão Geral](#-visão-geral)
- [Arquitetura](#-arquitetura)
- [Componentes Principais](#-componentes-principais)
- [Plugins Disponíveis](#-plugins-disponíveis)
- [Instalação e Uso](#-instalação-e-uso)
- [Configuração](#-configuração)
- [Desenvolvimento](#-desenvolvimento)
- [API Reference](#-api-reference)

## 🎯 Visão Geral

YRN oferece uma arquitetura de microserviços com os seguintes recursos principais:

- **Descoberta de Serviços**: Integração com HashiCorp Consul para registro e descoberta automática
- **Sistema de Plugins**: Arquitetura extensível com plugins para HTTP, Google Drive e autenticação OAuth
- **Validação de Esquemas**: Validação automática de entrada/saída usando JSON Schema
- **Flow Manager**: Sistema de orquestração de workflows com execução sequencial e paralela
- **Multi-tenancy**: Suporte nativo para múltiplos inquilinos
- **Observabilidade**: Métricas, logging e health checks integrados

## 🏗 Arquitetura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Agent       │    │   Connector     │    │      API        │
│   (Discovery)   │◄──►│  (Validation)   │◄──►│   (Interface)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────▼───────────────────────┘
                           ┌─────────────┐
                           │   Consul    │
                           │ (Registry)  │
                           └─────────────┘
```

### Componentes da Arquitetura

**Camada de Serviços (`cmd/`)**
- `agent` - Serviço de descoberta e agregação de esquemas
- `connector` - Serviço de validação baseado em JSON Schema
- `api` - API HTTP básica com endpoints de saúde

**Bibliotecas Compartilhadas (`pkg/`)**
- `ybase` - Framework core com registro no Consul e health checks
- `ylog` - Utilitários de logging estruturado
- `yctx` - Gerenciamento de contexto
- `plugin*` - Implementações de plugins (HTTP, Google Drive, etc.)

**Lógica de Negócio (`internal/`)**
- `producer` - Produção de mensagens/eventos
- `database` - Adaptadores para MongoDB e PostgreSQL
- `adapter` - Adaptadores para serviços externos
- `externalservice` - Integrações com serviços externos

**Módulos de Domínio (`module/`)**
- `flowmanager` - Gerenciamento de workflows e execução de plugins
- `team` - Gerenciamento de equipes
- `tenant` - Multi-tenancy
- `project` - Gerenciamento de projetos

## 🔧 Componentes Principais

### Flow Manager

O Flow Manager é o coração do sistema de orquestração:

```go
type Flow struct {
    Id               string       `json:"id"`
    Name             string       `json:"name"`
    Description      string       `json:"description"`
    Tenant           string       `json:"tenant"`
    FirstPluginToRun string       `json:"first_plugin_to_run"`
    Plugins          []FlowPlugin `json:"plugins"`
    Version          int          `json:"version"`
}
```

**Características:**
- Execução sequencial e condicional de plugins
- Compartilhamento de dados entre plugins
- Tratamento de erros configurável
- Versionamento de fluxos
- Suporte multi-tenant

### Sistema de Plugins

Todos os plugins implementam a interface `PluginExecutor`:

```go
type PluginExecutor interface {
    Do(ctx *yctx.Context, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (output any, err error)
}
```

**Validação de Esquemas**: Cada plugin possui um esquema JSON que define:
- Parâmetros de entrada obrigatórios e opcionais
- Tipos de dados aceitos
- Validações customizadas
- Documentação integrada

## 🔌 Plugins Disponíveis

### 1. HTTP Plugin (`pluginhttp`)

**Funcionalidade**: Executa requisições HTTP configuráveis com suporte a retry e timeout.

**Schema**:
```json
{
  "request": {
    "method": "POST|GET|PUT|DELETE|PATCH|HEAD|OPTIONS",
    "url": "https://api.exemplo.com/endpoint",
    "headers": {"Authorization": "Bearer {{.data.token}}"},
    "body": "Dados da requisição"
  },
  "retry": {
    "maxAttempts": 3,
    "delay": 1000
  }
}
```

**Recursos**:
- Suporte a todos os métodos HTTP
- Headers e query parameters customizáveis
- Template engine para dados dinâmicos
- Sistema de retry configurável
- Timeout personalizável

### 2. Google Drive Auth Plugin (`plugingdriveauth`)

**Funcionalidade**: Autentica com Google OAuth2 e obtém tokens de acesso.

**Schema**:
```json
{
  "client_id": "seu-client-id.googleusercontent.com",
  "client_secret": "seu-client-secret",
  "code": "authorization-code-from-oauth-flow",
  "redirect_uri": "http://localhost:8080/callback"
}
```

**Recursos**:
- Implementa fluxo OAuth2 Authorization Code
- Integração com Google APIs
- Tratamento automático de erros de autenticação
- Retorna access_token e refresh_token

### 3. Google Drive Plugin (`plugingdrive`)

**Funcionalidade**: Integração completa com Google Drive para listagem e download de arquivos.

**Schema**:
```json
{
  "credentials": "json-da-conta-de-servico",
  "folderId": "id-da-pasta-no-drive",
  "sharedDriveId": "id-do-shared-drive"
}
```

**Recursos**:
- Autenticação via Service Account
- Listagem de arquivos por pasta
- Download automático de conteúdo
- Suporte a Shared Drives
- Metadados completos (nome, tipo, data de modificação)
- Filtros por tipo de arquivo

**Exemplo de Saída**:
```json
[
  {
    "id": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
    "name": "documento.pdf",
    "mimeType": "application/pdf",
    "content": "conteúdo-do-arquivo-em-base64"
  }
]
```

## 🚀 Instalação e Uso

### Pré-requisitos

- Go 1.23+
- Docker e Docker Compose
- HashiCorp Consul (incluído no docker-compose)

### Quick Start com Docker

```bash
# Clone o repositório
git clone https://github.com/yrn-go/yrn.git
cd yrn

# Inicie todos os serviços
make compose-up

# Verifique se os serviços estão funcionando
curl http://localhost:8500/ui  # Consul UI
curl http://localhost:8080/health  # Agent health
curl http://localhost:8081/health  # Connector health
```

### Desenvolvimento Local

```bash
# Instale dependências
go mod download

# Execute os testes
make test

# Build os binários
make build

# Execute serviços individualmente
go run ./cmd/agent &
go run ./cmd/connector &
go run ./cmd/api &
```

## ⚙️ Configuração

### Variáveis de Ambiente

**Obrigatórias para todos os serviços:**
```bash
SERVICE_NAME=nome-do-servico        # Identificador do serviço no Consul
SERVICE_HOST=localhost              # Host onde o serviço está rodando
SERVICE_PORT=8080                   # Porta do serviço
CONSUL_HTTP_ADDR=localhost:8500     # Endereço do servidor Consul
```

**Específicas do Agent:**
```bash
CONNECTOR_SERVICE_NAME=connector    # Nome do serviço connector no Consul
```

**Opcionais para banco de dados:**
```bash
MONGO_URL=mongodb://localhost:27017      # String de conexão MongoDB
MONGO_DATABASE=yrn_database              # Nome da database MongoDB
REDIS_URL=redis://localhost:6379         # String de conexão Redis
```

### Configuração do Consul

O projeto utiliza Consul para descoberta de serviços. Cada serviço:

1. **Registra-se automaticamente** no Consul na inicialização
2. **Expõe endpoints** de health (`/health`) e schema (`/schema`)
3. **Descobre outros serviços** através de consultas ao Consul
4. **Atualiza status** através de health checks automáticos

## 🛠 Desenvolvimento

### Estrutura do Projeto

```
yrn/
├── cmd/                    # Pontos de entrada dos serviços
│   ├── agent/             # Serviço de descoberta
│   ├── connector/         # Serviço de validação
│   └── api/               # API HTTP
├── pkg/                   # Bibliotecas reutilizáveis
│   ├── ybase/            # Framework core
│   ├── plugin*/          # Implementações de plugins
│   └── ylog/             # Logging
├── internal/             # Lógica interna
│   ├── database/         # Adaptadores de banco
│   └── adapter/          # Adaptadores externos
├── module/               # Módulos de domínio
│   ├── flowmanager/      # Orquestração de workflows
│   └── team/             # Gerenciamento de equipes
└── infra/                # Infraestrutura como código
    └── apps/             # Configurações Kubernetes
```

### Adicionando Novos Plugins

1. **Crie um novo diretório** em `pkg/plugin{nome}`
2. **Implemente a interface** `PluginExecutor`
3. **Defina o JSON Schema** em `schema.json`
4. **Registre o plugin** no sistema

Exemplo básico:

```go
package pluginexemplo

import (
    _ "embed"
    "github.com/yrn-go/yrn/module/flowmanager"
    "github.com/yrn-go/yrn/pkg/yctx"
)

//go:embed schema.json
var Schema []byte

type Executor struct{}

func (e *Executor) Do(ctx *yctx.Context, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (output any, err error) {
    // Validar entrada
    requestData, err := plugincore.ValidateAndGetRequestBody[MySchema](Schema, schemaInputs, previousPluginResponse, responseSharedForAll)
    if err != nil {
        return nil, err
    }

    // Executar lógica do plugin
    result := processData(requestData)

    return result, nil
}
```

### Comandos de Desenvolvimento

```bash
# Construir projeto
make build

# Executar testes
make test

# Linting
make lint

# Construir imagem Docker
make docker-build

# Executar localmente
make run

# Limpar artefatos
make clean
```

## 📚 API Reference

### Agent Service

**Endpoint**: `http://localhost:8080`

- `GET /health` - Health check
- `GET /schema` - Retorna esquema do serviço
- `GET /services` - Lista todos os serviços descobertos

### Connector Service

**Endpoint**: `http://localhost:8081`

- `GET /health` - Health check
- `GET /schema` - Retorna esquema de validação
- `POST /validate` - Valida dados contra schema

### Flow Execution

Workflows são executados através do FlowManager com a seguinte estrutura:

```json
{
  "id": "workflow-001",
  "name": "Processamento de Dados",
  "description": "Workflow para processar dados do Google Drive",
  "tenant": "cliente-001",
  "first_plugin_to_run": "gdrive-auth",
  "plugins": [
    {
      "id": "gdrive-auth",
      "slug": "gdrive-auth",
      "name": "Autenticação Google Drive",
      "schema_input": "{\"client_id\": \"{{.env.GOOGLE_CLIENT_ID}}\"}",
      "next_to_be_executed": ["gdrive-list"]
    },
    {
      "id": "gdrive-list",
      "slug": "google-drive",
      "name": "Listar Arquivos",
      "schema_input": "{\"credentials\": \"{{.data.access_token}}\"}",
      "next_to_be_executed": ["http-post"]
    }
  ]
}
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Padrões de Código

- Siga as convenções do Go (`go fmt`, `go vet`)
- Execute `make lint` antes do commit
- Mantenha cobertura de testes alta
- Documente APIs públicas
- Use semantic versioning

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🔗 Links Úteis

- [Documentação do Consul](https://www.consul.io/docs)
- [JSON Schema Specification](https://json-schema.org/)
- [Google Drive API](https://developers.google.com/drive/api)
- [Gin Web Framework](https://github.com/gin-gonic/gin)

---

**Desenvolvido com ❤️ em Go**