# Variáveis
BINARY_NAME=yrn
GO=go
DOCKER=docker
DOCKER_COMPOSE=docker-compose
GOLANGCI_LINT=$(shell go env GOPATH)/bin/golangci-lint

# Comandos padrão
.PHONY: all build clean test run lint install-lint docker-build docker-run help

all: clean build

# Build
build:
	@echo "Construindo o projeto..."
	$(GO) build -o $(BINARY_NAME) ./cmd/...

# Limpeza
clean:
	@echo "Limpando arquivos..."
	rm -f $(BINARY_NAME)
	$(GO) clean

# Testes
test:
	@echo "Executando testes..."
	$(GO) test -v ./...

# Execução
run:
	@echo "Executando o projeto..."
	$(GO) run ./cmd/...

# Lint
install-lint:
	@echo "Instalando golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2
	@echo "Criando arquivo de configuração do golangci-lint..."
	@echo 'linters-settings:\n  govet:\n    check-shadowing: true\n  gocyclo:\n    min-complexity: 15\n  goconst:\n    min-len: 2\n    min-occurrences: 2\n  misspell:\n    locale: US\n  lll:\n    line-length: 140\n\nlinters:\n  enable:\n    - bodyclose\n    - depguard\n    - errcheck\n    - goconst\n    - gocritic\n    - gocyclo\n    - gofmt\n    - goimports\n    - gosec\n    - gosimple\n    - govet\n    - ineffassign\n    - lll\n    - misspell\n    - nakedret\n    - staticcheck\n    - stylecheck\n    - typecheck\n    - unconvert\n    - unparam\n    - unused\n    - whitespace\n\nissues:\n  exclude-rules:\n    - path: _test\.go\n      linters:\n        - dupl\n        - gocyclo\n        - lll\n        - gosec\n    - path: mock\n      linters:\n        - typecheck\n        - unused\n        - gosec\n    - path: vendor\n      linters:\n        - all\n    - path: pkg/mod\n      linters:\n        - all\n    - linters:\n        - gosec\n      text: "G104"' > .golangci.yml

lint: install-lint
	@echo "Executando análise estática..."
	$(GOLANGCI_LINT) run --timeout=5m --skip-dirs vendor --skip-dirs pkg/mod

# Docker
docker-build:
	@echo "Construindo imagem Docker..."
	$(DOCKER) build -t $(BINARY_NAME) .

docker-run:
	@echo "Executando container Docker..."
	$(DOCKER) run -p 8080:8080 $(BINARY_NAME)

# Docker Compose
compose-up:
	@echo "Iniciando serviços com Docker Compose..."
	$(DOCKER_COMPOSE) up -d

compose-down:
	@echo "Parando serviços com Docker Compose..."
	$(DOCKER_COMPOSE) down

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make all        - Limpa e constrói o projeto"
	@echo "  make build      - Constrói o projeto"
	@echo "  make clean      - Remove arquivos gerados"
	@echo "  make test       - Executa os testes"
	@echo "  make run        - Executa o projeto"
	@echo "  make install-lint - Instala o golangci-lint"
	@echo "  make lint       - Executa análise estática"
	@echo "  make docker-build - Constrói a imagem Docker"
	@echo "  make docker-run   - Executa o container Docker"
	@echo "  make compose-up  - Inicia serviços com Docker Compose"
	@echo "  make compose-down - Para serviços com Docker Compose"
	@echo "  make help       - Mostra esta mensagem de ajuda" 