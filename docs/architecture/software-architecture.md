# Arquitetura de Software — YRN

> Descreve a arquitetura **real observada** no repo, não um ideal imposto.
> A escolha de estilo está registrada em `docs/adr/0001-architecture-style.md`.

## Estilo: Microsserviços com descoberta de serviços (Consul)
O YRN é um conjunto de serviços Go independentes que se registram e se descobrem via
**Consul**. O framework comum `pkg/ybase` padroniza o ciclo de vida de cada serviço.

> ℹ️ O layout interno (`module/` como domínio, `internal/` como infra/adapters) tem
> **influência hexagonal**, mas a regra de dependência "domínio não importa infra"
> **não é hoje verificada/garantida**. Formalizar (ou não) hexagonal por serviço é uma
> decisão em aberto → candidata a ADR. Não assuma hexagonal estrito ao revisar.

## Serviços (`cmd/`)
- **agent** — descobre serviços e agrega `/schema` de outros serviços.
- **connector** — valida requisições contra JSON schemas predefinidos.
- **api** — API HTTP (Gin) com `/health`.

Cada serviço via `ybase.NewApp()`:
1. Registra no Consul (env vars), 2. expõe `/health`, 3. expõe `/schema` (JSON),
4. configura health check automático.

## Layout de pacotes
```
cmd/<svc>/main.go     composition root de cada serviço (config + ybase + handlers)
pkg/
  ybase/              framework comum: registro Consul, health, schema endpoint
  ylog/  yctx/        logging estruturado e propagação de context
  plugincore/         contrato base de plugin
  pluginhttp/  plugingdrive/  plugingdriveauth/  pluginmapper/   integrações plugáveis
module/               domínio de negócio
  flowmanager/        gestão de fluxos
  tenant/  team/  project/   multi-tenancy e organização
internal/             infra e adapters (não exportável)
  database/{mongodb,postgres}/   adapters de persistência
  producer/           produção de mensagens/eventos
  externalservice/    integrações externas
  test/               helpers de teste
```

## Comunicação entre serviços
- Descoberta dinâmica via Consul (sem URLs hardcoded entre serviços).
- Contrato exposto via endpoint `/schema` — **mudança de schema é mudança de contrato público**:
  versionar e revisar com humano.

## Plugins
- `pkg/plugincore` define o contrato; cada `pkg/plugin<x>` é uma integração concreta.
- Plugin novo = pacote novo implementando o contrato; sem tocar no core.

## Erros
- Wrap com `%w` ao subir; `errors.Is/As` para checar. Mapeie erro→status HTTP num só lugar (boundary Gin).

## Concorrência
- Descoberta de serviços e plugins envolvem goroutines: toda goroutine recebe `context.Context`
  e tem término garantido; estado compartilhado protegido. **Rodar `-race`** (gap atual).

## Config
- 12-factor: tudo por env (`SERVICE_*`, `CONSUL_HTTP_ADDR`), validado no startup. Segredos nunca no código.

## Anti-padrões a evitar
- Hardcode de endpoint de outro serviço (use Consul). Lógica de negócio no handler Gin.
- `context.TODO()` em produção. Pacote `utils`/`common` genérico (vira lixeira).
- Mudar `/schema` sem versionar o contrato.

## Questões em aberto (decidir + virar ADR se estrutural)
- Adotar `-race` no CI/Makefile.
- Formalizar (ou não) regra de dependência hexagonal por serviço.
- Padrão de comunicação entre serviços além de descoberta (REST? mensageria? `internal/producer` sugere fila — definir).
