# Especificação Funcional: Redis Connection & Key Inspector

## 1. Visão Geral
Esta funcionalidade provê conexão e inspeção de dados Redis no GoScratch. Ela permite cadastrar conexões, navegar por chaves e visualizar valores armazenados em estruturas como strings, hashes, lists, sets e sorted sets.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-09.1: Navegação por Chaves Redis
**Como** desenvolvedor trabalhando com Redis,  
**Quero** conectar em uma instância e listar suas chaves,  
**Para** inspecionar rapidamente o conteúdo armazenado sem recorrer a ferramentas externas.

* **Cenário 09.1.1: Listagem de chaves por padrão**
  * **Given** que existe uma conexão Redis válida configurada
  * **When** acesso o inspetor e informo um padrão de busca
  * **Then** o sistema lista as chaves compatíveis com nome, tipo e TTL quando disponível.

* **Cenário 09.1.2: Inspeção de valor por tipo de dado**
  * **Given** que selecionei uma chave existente
  * **When** solicito a inspeção da chave
  * **Then** o sistema exibe o valor formatado de acordo com o tipo Redis correspondente.

## 3. Requisitos Funcionais (FR)
- **FR-09.1:** Cadastro e reutilização de conexões Redis.
- **FR-09.2:** Listagem de chaves com suporte a filtro por padrão.
- **FR-09.3:** Inspeção formatada de valores por tipo de dado Redis.
