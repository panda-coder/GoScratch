# Especificação Funcional: Snippet Library

## 1. Visão Geral
Esta funcionalidade provê uma biblioteca reutilizável de snippets no GoScratch. Ela permite salvar, localizar e reutilizar trechos de código recorrentes sem precisar recriá-los manualmente a cada sessão.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-04.1: Reutilização Rápida de Snippets
**Como** desenvolvedor Go,  
**Quero** salvar snippets reutilizáveis com título e categorização,  
**Para** inserir rapidamente trechos comuns no editor principal.

* **Cenário 04.1.1: Salvamento de um snippet novo**
  * **Given** que escrevi um trecho útil no editor
  * **When** escolho a ação de salvar como snippet
  * **Then** o sistema persiste o conteúdo com título e metadados básicos.

* **Cenário 04.1.2: Busca de snippet por termo**
  * **Given** que existem snippets previamente cadastrados
  * **When** digito um termo de busca
  * **Then** o sistema lista apenas os snippets compatíveis com o filtro informado.

## 3. Requisitos Funcionais (FR)
- **FR-04.1:** Suporte ao cadastro de snippets com título e conteúdo.
- **FR-04.2:** Suporte à busca e listagem de snippets salvos.
- **FR-04.3:** Inserção de snippets no editor principal.
