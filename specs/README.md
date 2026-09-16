# Specs

> Uma spec por feature: **o quê** e **o porquê** antes do **como**.
> 1 spec = 1 branch = 1 PR. A spec é a fonte que o plano e as tasks referenciam.

## Formato (`NNNN-nome-da-feature.md`)
```markdown
# <Feature>

## Problema
<que dor, para quem, por que agora.>

## Decisão
<abordagem escolhida + 1 alternativa descartada e o motivo.>

## Requisitos
- R1 ...
- R2 ...

## Critérios de aceite
<viram testes — caminho feliz + edge + falha.>

## Fora de escopo
<o que NÃO entra — mantém a fatia pequena.>
```

Gere via skill `plan-feature` ("planeje X"). Decisões irreversíveis na spec → também viram ADR
em `docs/adr/`. Para o YRN, mudanças em `/schema` de um serviço são contrato público — sinalize.
