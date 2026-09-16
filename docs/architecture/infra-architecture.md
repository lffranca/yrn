# Arquitetura de Infraestrutura — YRN

> Como o YRN roda, é implantado e observado. `<...>` = a confirmar com o time.

## Runtime e topologia
- Serviços em **containers Docker**. Em dev, `docker-compose.yaml` sobe **Consul (8500)**,
  **agent (8080)** e **connector (8081)**.
- Agentes de execução de fluxo gerenciam **containers do usuário via Docker** (superfície sensível).
- Manifests de deploy em `infra/apps/` (**kustomize**: `deployment.yaml` + `kustomization.yaml`).
- `infra/portainer/` contém stacks (ex.: `stack_n8n.yaml`) — n8n usado como `<referência/integração — confirmar>`.

## Ambientes
| Ambiente | Uso | URL/cluster | Diferenças |
|---|---|---|---|
| dev | docker-compose local | localhost | Consul+agent+connector |
| staging | `<...>` | `<...>` | `<...>` |
| prod | `<...>` | `<...>` | `<...>` |

## CI/CD
- **CI:** `.github/workflows/test-and-build.yaml` — roda `go test ./... -v` e build
  cross-platform (linux/darwin/windows amd64) de agent/connector/api; build e push de
  imagem Docker para **ghcr.io**.
  > Gaps: sem `-race`, sem `govulncheck`/`gosec` como etapa dedicada (gosec só via golangci-lint, que não roda no CI hoje). Considerar adicionar `make lint` ao CI.
- **CD:** release por **tag → atualização de imagem (rollout)** — ver commits
  `chore: update rollout image to vX`. Mecanismo exato (ArgoCD/Portainer): `<confirmar>`.
- Estratégia de release: `<rolling / canary — confirmar>`.

## Configuração e segredos
- Config via **env vars** (`SERVICE_NAME/HOST/PORT`, `CONSUL_HTTP_ADDR`, `CONNECTOR_SERVICE_NAME`).
- **Nunca** segredo no repo. Armazenamento de segredos: `<Vault / Secrets Manager / SOPS — definir>`.

## Observabilidade
- **Logs:** estruturados via `pkg/ylog`. Destino central: `<definir>`.
- **Métricas:** `<Prometheus/OTel — definir SLIs>`.
- **Tracing:** propagação de `context` (`pkg/yctx`); backend `<OTel? — definir>`.
- **Health:** cada serviço expõe `/health`, monitorado pelo Consul.
- **Alertas / SLOs:** `<definir e quem é paginado>`.

## Dados
- Adapters para **MongoDB** e **PostgreSQL** (`internal/database/`). Qual é a fonte de
  verdade de cada domínio, migrations, backup/retenção: `<definir>`.

## Resiliência
- Timeouts, retries idempotentes, circuit breaker e rate limiting entre serviços: `<definir>`.
- Isolamento dos containers de execução de fluxo: `<política de sandbox — definir, é risco de segurança>`.
