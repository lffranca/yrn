# Regras do Projeto

> O "porquê" e as fronteiras do YRN. Quem entra no time lê isto primeiro.
> Seções com `<...>` dependem de decisão do time — não invente, pergunte.

## Propósito
Plataforma de automação **low-code/no-code** que permite criar fluxos lógicos e
executá-los em **agentes remotos** gerenciados via Docker. Inspirada no n8n,
diferencia-se pela orquestração distribuída de agentes em ambientes isolados.

## Escopo
- **Faz parte:** orquestração de fluxos, descoberta de serviços (Consul), validação
  de schemas (connector), execução em agentes, plugins de integração (HTTP, Google Drive).
- **NÃO faz parte:** `<delimitar — ex.: UI yrn-admin vive em outro repo? billing? auth de usuário final?>`

## Glossário do domínio
- **Flow / Fluxo** — sequência lógica de passos automatizados (`module/flowmanager`).
- **Agent** — serviço que descobre outros serviços e agrega schemas (`cmd/agent`).
- **Connector** — valida requisições contra JSON schemas (`cmd/connector`).
- **Tenant** — unidade de multi-tenancy (`module/tenant`).
- **Plugin** — integração plugável (`pkg/plugin*`): http, gdrive, mapper, core.
- `<Adicionar: Project, Team, Resource, Credential conforme o negócio define.>`

## Padrões de código (Go)
- `gofmt` + `goimports` obrigatórios. Lint: `.golangci.yml` (line-length 140; `gosec`,
  `gocyclo`, `errcheck`, `govet` shadow, etc. habilitados).
- Nomes: pacotes curtos, minúsculos, sem underscore; exporte o mínimo; receivers consistentes.
- Avalie a **stdlib antes** de adicionar dependência; nova dep precisa de justificativa no PR.
- Arquitetura: ver `software-architecture.md`.

## Política de testes
- Cubra caminho feliz + edge + falha. Tempo/IO injetados para testabilidade.
- **Adotar `-race`** (hoje ausente — ver questão em aberto). Concorrência sem teste de race é risco.
- Domínio (`module/*`) rápido e table-driven; adapters (`internal/database/*`, plugins) com integração.
- Piso de cobertura: `<definir — foque em valor, não em número>`.

## Erros e observabilidade
- Erros wrapped com `%w`; sentinelas no domínio. Logs estruturados (`pkg/ylog`) com correlation id.
- **Sem PII em log.** Propagação de `context` (`pkg/yctx`) nas chamadas entre serviços.

## Pull Requests
- Pequenos (idealmente < ~400 linhas). Draft cedo. Conventional Commits.
- Branch: `feat/ fix/ chore/ refactor/ test/ docs/`. 1 spec = 1 branch = 1 PR.
- CI **verde obrigatório**. Aprovação humana sempre.

## Segurança (não se delega)
- Auth, autorização, input e segredos → **review humano**.
- Execução de fluxos em agentes Docker é superfície sensível: isolamento e validação de
  comandos antes de tocar no Docker daemon. `gosec` no lint.

## Definition of Ready / Done
- **DoR:** spec aprovada, critérios de aceite claros, dependências conhecidas.
- **DoD:** testes verdes, lint limpo, doc/contrato (schema) atualizado, PR revisado e CI verde.

## Versionamento e releases
- SemVer; Conventional Commits alimentam o changelog. Release por **tag → rollout de imagem**
  (ver histórico `chore: update rollout image to vX`). CD via kustomize/ArgoCD — ver `infra-architecture.md`.

## Decisões (ADRs)
Decisões estruturantes/irreversíveis → `docs/adr/`. Estilo arquitetural em `0001-architecture-style.md`.

## Stakeholders / donos
`<Quem decide arquitetura, contrato de serviços e infra; quem aprova decisões irreversíveis.>`
