# AGENTS.md

> Contexto portável do projeto, lido por Claude Code / OpenCode / Droid / Codex.
> Detalhes de arquitetura em `docs/architecture/`. Decisões em `docs/adr/`.

## O que é este projeto
**YRN** — plataforma distribuída de automação low-code/no-code (inspirada no n8n),
escrita em Go. Orquestra fluxos lógicos executados em **agentes remotos** sobre Docker.
Serviços se descobrem via **Consul**; o framework comum é `pkg/ybase`. Visão geral
em `docs/proposals/P0001-yrn-overview.md`.

## Serviços (`cmd/`)
- **agent** (`cmd/agent`) — descoberta de serviços e agregação de schemas.
- **connector** (`cmd/connector`) — validação de requisições contra JSON schemas.
- **api** (`cmd/api`) — API HTTP básica com `/health`.

## Comandos (reais deste repo)
- Testes:   `make test`   → `go test -v ./...`
- Build:    `make build`  → `go build -o <bin> ./cmd/...`
- Lint:     `make lint`   → `golangci-lint run` (v1.55.2, line-length 140, `gosec` incluso)
- Format:   `gofmt -w .` / `goimports -w .`
- Run:      `make run` · serviços individuais: `go run ./cmd/{agent,connector,api}`
- Docker:   `make docker-build` · `make compose-up` (Consul + agent + connector)

> ⚠️ Hoje `go test` roda **sem** `-race`. Recomenda-se adicionar `-race` no CI/Makefile
> (há concorrência real: descoberta de serviços, plugins). Ver questões em aberto no review.

## Variáveis de ambiente (todos os serviços)
`SERVICE_NAME` · `SERVICE_HOST` · `SERVICE_PORT` · `CONSUL_HTTP_ADDR`
Agent: `CONNECTOR_SERVICE_NAME`.

## Convenção de branch
Sempre a partir do `main` atualizado. Nunca commitar direto no `main`.
```
feat/<escopo>-<curto>   fix/<escopo>-<curto>
chore/<curto>   refactor/<curto>   test/<curto>   docs/<curto>
```
Uma spec = uma branch = um PR. Cada task pequena = um commit.

## Convenção de commit (Conventional Commits)
`feat:` `fix:` `test:` `refactor:` `chore:` `docs:` — alimenta changelog/semver e o
fluxo de release (rollout por tag). Mensagem no imperativo, escopo opcional: `feat(api): ...`.

## Fluxo de PR
1. Branch + commits pequenos.
2. `make test` e `make lint` verdes localmente.
3. `gh pr create --draft` com corpo apontando a spec e o resultado de QA/sec.
4. Review (code-review-backend / go-review / go-qa). CI roda os mesmos checks.

## Convenções que o review COBRA (não são sugestões)
1. Erros explícitos e wrapped (`%w`). 2. Concorrência com término garantido por `context`.
3. Schema/contrato entre serviços é interface pública — mudança requer review humano.
4. Segurança (auth/input/segredos) sempre com review humano.

## Limites para agentes de IA
Decisões irreversíveis (schema de serviço, contrato Consul, infra/kustomize) → ADR + review humano.
Segurança nunca se delega. Não aplicar fix sem entender a causa raiz.
