# Flow Manager

O Flow Manager é um módulo responsável por gerenciar a execução de fluxos de plugins em uma arquitetura orientada a eventos. Este módulo fornece uma estrutura flexível e robusta para a execução sequencial ou paralela de plugins, com suporte a métricas, tratamento de erros e monitoramento de estado.

## Índice

- [Visão Geral](#visão-geral)
- [Componentes](#componentes)
- [Métricas](#métricas)
- [Uso](#uso)
- [Testes](#testes)
- [Configuração](#configuração)

## Visão Geral

O Flow Manager é projetado para:
- Gerenciar a execução de plugins em um fluxo definido
- Coletar métricas de execução (tempo, memória, CPU)
- Tratar erros e recuperar de panics
- Manter o estado de execução dos plugins
- Permitir compartilhamento de dados entre plugins

## Componentes

### EventManager

O `EventManager` é o componente principal responsável por:
- Registrar plugins no fluxo
- Executar plugins na ordem definida
- Coletar métricas de execução
- Gerenciar o estado dos plugins
- Tratar erros e panics

```go
type EventManager struct {
    pluginManager PluginManager
    statusRepo    PluginStatusRepository
    metrics      map[string]PluginMetrics
}
```

### FlowExecutor

O `FlowExecutor` é responsável por:
- Carregar a definição do fluxo
- Inicializar o EventManager
- Coordenar a execução do fluxo completo

```go
type FlowExecutor struct {
    flowReaderRepository FlowReaderRepository
    pluginManager        PluginManager
    statusRepo           PluginStatusRepository
}
```

### PluginExecutor

Interface que define o contrato para execução de plugins:

```go
type PluginExecutor interface {
    Do(ctx *yctx.Context, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (output any, err error)
}
```

## Métricas

O módulo coleta as seguintes métricas para cada plugin:

- **Tempo de Execução**
  - StartTime: Momento de início
  - EndTime: Momento de término
  - ExecutionTime: Duração total

- **Uso de Recursos**
  - MemoryBefore: Uso de memória antes da execução
  - MemoryAfter: Uso de memória após a execução
  - CPUBefore: Uso de CPU antes da execução
  - CPUAfter: Uso de CPU após a execução

## Uso

### Configuração Básica

```go
// Criar instância do FlowExecutor
executor := NewFlowExecutor(
    flowReaderRepository,
    pluginManager,
    statusRepo,
)

// Executar um fluxo
response, err := executor.Do(ctx, flowId, eventData)
```

### Definição de Plugin

```go
plugin := FlowPlugin{
    Id:               "plugin-id",
    Slug:             "plugin-slug",
    SchemaInput:      `{"field": "value"}`,
    NextToBeExecuted: []string{"next-plugin-id"},
}
```

## Testes

O módulo inclui testes abrangentes que cobrem:

- Execução bem-sucedida de plugins
- Tratamento de erros de repositório
- Coleta de métricas
- Tratamento de erros de plugin
- Recuperação de panics
- Plugins de longa duração

Para executar os testes:

```bash
go test -v ./module/flowmanager/...
```

## Configuração

### Ambiente de Desenvolvimento

1. Instale as dependências:
```bash
make install-deps
```

2. Execute os testes:
```bash
make test
```

3. Execute a análise estática:
```bash
make lint
```

### Integração

Para integrar o Flow Manager em seu projeto:

1. Implemente as interfaces necessárias:
   - PluginExecutor
   - PluginManager
   - PluginStatusRepository

2. Configure o FlowExecutor com suas implementações

3. Defina seus fluxos e plugins

### Boas Práticas

- Sempre utilize o contexto para controle de timeout
- Implemente tratamento de erros adequado
- Monitore as métricas coletadas
- Mantenha os plugins idempotentes
- Documente o schema de entrada/saída dos plugins

## Contribuição

1. Fork o repositório
2. Crie uma branch para sua feature
3. Adicione testes
4. Envie um pull request

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).
