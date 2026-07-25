# Especificação Funcional: Database Connection Manager & Schema Explorer

## 1. Visão Geral
Gerenciamento de conexões SQL (SQLite, PostgreSQL, MySQL) e exibição em árvore na barra lateral esquerda (Explorer) permitindo navegar em bancos, tabelas e colunas com indicação de chaves primárias.

## 2. Histórias de Usuário (BDD)

### US-02.1: Navegação por Metadados da Tabela
**Como** desenvolvedor,  
**Quero** expender o nó de uma conexão de banco na barra lateral,  
**Para** visualizar todas as tabelas e suas colunas com seus tipos de dados.

* **Cenário 02.1.1: Visualização de Colunas de Tabela**
  * **Given** que selecionei uma conexão de banco de dados ativa
  * **When** expando o nó da tabela "users"
  * **Then** o sistema exibe a lista de colunas `id (INTEGER, PK)`, `name (TEXT)`, `email (TEXT)`.
