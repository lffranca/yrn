# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build & Test
- `make build` - Build the project binaries
- `make test` - Run all tests with verbose output
- `make run` - Run the project (builds all cmd entries)
- `make clean` - Clean build artifacts

### Code Quality
- `make lint` - Run golangci-lint (auto-installs if needed)
- `make install-lint` - Install golangci-lint v1.55.2 and generate config

### Docker & Services
- `make docker-build` - Build Docker image
- `make docker-run` - Run container on port 8080
- `make compose-up` - Start all services (Consul + Agent + Connector)
- `make compose-down` - Stop all services

### Individual Service Commands
- `go run ./cmd/agent` - Run agent service
- `go run ./cmd/connector` - Run connector service
- `go run ./cmd/api` - Run API service

## Architecture Overview

YRN is a distributed service orchestration platform built in Go using a microservices architecture with Consul for service discovery.

### Core Components

**Services (`cmd/`):**
- **Agent** (`cmd/agent/main.go`) - Service discovery and schema aggregation service that queries other services for their schemas
- **Connector** (`cmd/connector/main.go`) - JSON schema-based service that validates requests against predefined schemas
- **API** (`cmd/api/main.go`) - Basic HTTP API service with health endpoint

**Shared Libraries (`pkg/`):**
- **ybase** - Core service framework providing Consul registration, health checks, and schema endpoints
- **ylog** - Logging utilities
- **yctx** - Context handling
- **plugin*** - Various plugin implementations (HTTP, Google Drive, etc.)

**Business Logic (`internal/`):**
- **producer** - Message/event production
- **database** - Database adapters (MongoDB, PostgreSQL)
- **adapter** - External service adapters
- **externalservice** - External service integrations

**Domain Modules (`module/`):**
- **flowmanager** - Workflow management
- **team** - Team management
- **tenant** - Multi-tenancy
- **project** - Project management

### Service Discovery Pattern

All services use the `ybase.NewApp()` pattern which:
1. Registers service with Consul using environment variables
2. Exposes `/health` for health checks
3. Exposes `/schema` endpoint returning JSON schema
4. Configures automatic health checking via Consul

### Environment Variables

Required for all services:
- `SERVICE_NAME` - Service identifier for Consul
- `SERVICE_HOST` - Host address for service
- `SERVICE_PORT` - Port for service to listen on
- `CONSUL_HTTP_ADDR` - Consul server address

Agent-specific:
- `CONNECTOR_SERVICE_NAME` - Name of connector service to discover

### Docker Compose Setup

The system runs as three services:
- **consul** (port 8500) - Service registry and discovery
- **agent** (port 8080) - Service discovery agent
- **connector** (port 8081) - Schema validation service

### Code Style

- Uses golangci-lint with 140 character line limit
- Follows standard Go project layout with `cmd/`, `internal/`, `pkg/` structure
- Services are built using Gin web framework
- All services include Consul integration via `pkg/ybase`

### Testing

Run tests with `make test` which executes `go test -v ./...` across all packages.