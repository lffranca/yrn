# 0001 — Estilo arquitetural: Microsserviços com descoberta via Consul

- **Status:** Accepted (documenta o estado atual do repo)
- **Data:** 2026-06-12
- **Decisores:** `<time/tech lead — confirmar>`

## Contexto
O YRN orquestra fluxos de automação executados em agentes distribuídos. Serviços
precisam se localizar dinamicamente, expor seus schemas e ser implantados/escalados
de forma independente.

## Decisão
Adotar **microsserviços** Go independentes (`cmd/agent`, `cmd/connector`, `cmd/api`),
com **Consul** para registro e descoberta de serviços e o framework comum **`pkg/ybase`**
padronizando registro, `/health` e `/schema`. Integrações externas entram como **plugins**
(`pkg/plugin*`) sobre `pkg/plugincore`. Detalhes em `../architecture/software-architecture.md`.

## Alternativas consideradas
- **Monólito único:** rejeitado — não atende execução distribuída em agentes remotos.
- **Service mesh / k8s DNS sem Consul:** possível no futuro; Consul já está em uso e
  cobre descoberta + health hoje.

## Consequências
- ✅ Serviços evoluem e escalam de forma independente; descoberta dinâmica sem URLs fixas.
- ✅ Contrato explícito por serviço via `/schema`.
- ⚠️ Complexidade operacional (Consul, rede, observabilidade distribuída).
- ⚠️ Mudança de `/schema` é mudança de contrato público → versionar e revisar.

## Em aberto (podem virar ADRs próprios)
- Formalizar regra de dependência **hexagonal** por serviço (hoje não verificada).
- Padrão de comunicação inter-serviço (REST vs. mensageria; `internal/producer` sugere fila).
- Fonte de verdade entre **MongoDB** e **PostgreSQL**.

## Superseded by
Nenhum.
