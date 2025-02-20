# P0001 - YRN: Plataforma de Automação Low-Code/No-Code

## 📌 Resumo
O **YRN** é um sistema de criação de fluxos lógicos no modelo **low-code/no-code**, permitindo a automação de tarefas e integração de serviços. Inspirado no **n8n**, o YRN é desenvolvido em **Golang** e se diferencia pelo gerenciamento distribuído de agentes para orquestração de fluxos dentro de ambientes isolados (Docker).

## 🎯 Objetivos
- Criar uma plataforma **modular** e **extensível** para automação.
- Possibilitar a execução de fluxos diretamente em **agentes remotos** gerenciados via Docker.
- Oferecer uma interface intuitiva para criação e gerenciamento de fluxos.

## 🏗 Arquitetura do Projeto
### 📌 Componentes principais:
1. **yrn-admin (React.js)**
    - Interface web para gerenciar fluxos, credenciais, agentes e recursos.

2. **yrn-api (Golang - REST API)**
    - Backend responsável pelo gerenciamento de fluxos e agentes.

3. **yrn-agent (Golang - Dockerized)**
    - Aplicação que gerencia containers do usuário conforme instruções do `yrn-api`.

### 📌 Fluxo de Operação
1. O usuário configura um fluxo no **yrn-admin**.
2. O **yrn-api** recebe a requisição, valida e define as ações.
3. Se necessário, o **yrn-api** instrui um **yrn-agent** remoto a executar comandos em containers.
4. O **yrn-agent** interage com o Docker para iniciar/parar serviços conforme necessário.
5. O resultado da execução é retornado ao **yrn-api** e, consequentemente, ao **yrn-admin**.

## 🚀 Plano de Desenvolvimento
| Mês  | Tarefa                        | Status  |
|------|--------------------------------|---------|
| M1   | Definição da arquitetura       | 🔄 Em andamento |
| M2   | Protótipo inicial do yrn-admin | 🛠 Em desenvolvimento |
| M3   | API básica do yrn-api          | ⏳ Planejado |
| M4   | Implementação do yrn-agent     | ⏳ Planejado |
| M5   | Testes e melhorias             | ⏳ Planejado |

## ⚠️ Riscos e Desafios
- Segurança na execução de fluxos distribuídos.
- Escalabilidade e gerenciamento de múltiplos agentes.
- Garantia de consistência entre os agentes e a API.

## ✅ Critérios de Sucesso
- O sistema permite a criação e execução de fluxos distribuídos via Docker.
- A interface `yrn-admin` permite gerenciar fluxos sem complexidade técnica.
- A API `yrn-api` gerencia os agentes de forma eficiente e segura.
- O `yrn-agent` pode ser registrado em múltiplos ambientes sem dependências externas complexas.
