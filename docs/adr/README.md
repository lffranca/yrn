# Architecture Decision Records (ADR)

> Cada decisão **irreversível ou cara de reverter** vira um registro aqui. ADR é
> histórico: não se apaga, se **supera** (status `Superseded by NNNN`).

## Quando criar um ADR
- Escolha de tecnologia estruturante (banco, fila, framework, descoberta de serviços).
- Contrato público / schema (`/schema`) que outros serviços dependem.
- Decisão de infra com custo de reversão (topologia, multi-tenancy, kustomize/rollout).

## Formato (`NNNN-titulo-curto.md`)
```markdown
# NNNN — Título
- Status: Proposed | Accepted | Superseded by MMMM
- Data: YYYY-MM-DD
- Decisores: <quem>

## Contexto
<forças em jogo, restrições.>

## Decisão
<o que foi decidido, no imperativo.>

## Alternativas consideradas
<opções e por que foram descartadas.>

## Consequências
<trade-offs positivos e negativos resultantes.>
```

## Índice
- [0001 — Microsserviços com descoberta via Consul](./0001-architecture-style.md)

> Nota: `docs/decision-records/` é um stub vazio anterior — consolidar aqui ou remover.
