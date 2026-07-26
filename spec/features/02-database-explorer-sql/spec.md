# Especificação Funcional: Database Explorer & Schema Inspector

## 1. Visão Geral
Esta funcionalidade provê a exploração de conexões SQL e a inspeção de schema no GoScratch. Ela permite navegar por bancos, tabelas e colunas em uma estrutura visual, facilitando a descoberta de metadados antes da escrita de consultas.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-02.1: Navegação por Tabelas e Colunas
**Como** desenvolvedor trabalhando com bancos relacionais,  
**Quero** expandir a árvore de uma conexão ativa para visualizar tabelas e colunas,  
**Para** compreender rapidamente o schema disponível antes de escrever consultas.

* **Cenário 02.1.1: Visualização de tabelas da conexão ativa**
  * **Given** que existe uma conexão válida selecionada
  * **When** abro o explorador de banco de dados
  * **Then** o sistema lista as tabelas disponíveis para navegação.

* **Cenário 02.1.2: Visualização de colunas com tipo e chave primária**
  * **Given** que selecionei uma tabela no explorador
  * **When** expando o nó correspondente
  * **Then** o sistema exibe as colunas com seus tipos e indicação de chave primária quando aplicável.

## 3. Requisitos Funcionais (FR)
- **FR-02.1:** Suporte à listagem de tabelas a partir de uma conexão ativa.
- **FR-02.2:** Suporte à listagem de colunas com metadados de tipo e chave primária.
- **FR-02.3:** Exibição hierárquica em árvore para navegação de schema.
