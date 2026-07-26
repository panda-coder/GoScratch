# Especificação Funcional: Query Builder

## 1. Visão Geral
Esta funcionalidade provê uma experiência assistida de construção de consultas SQL no GoScratch. Ela permite compor consultas com sugestões contextuais, reduzir erros de digitação e acelerar a exploração de dados.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-03.1: Montagem Assistida de Consultas SQL
**Como** desenvolvedor que explora bancos de dados,  
**Quero** receber sugestões de palavras-chave, tabelas e colunas durante a digitação,  
**Para** montar consultas válidas com mais rapidez e menos erro manual.

* **Cenário 03.1.1: Sugestão de Tabelas após `FROM`**
  * **Given** que existe uma conexão ativa com metadados disponíveis
  * **When** digito `SELECT * FROM ` no editor de consulta
  * **Then** o sistema exibe sugestões de tabelas compatíveis com o contexto atual.

* **Cenário 03.1.2: Sugestão de Colunas após seleção de tabela**
  * **Given** que escolhi uma tabela para a consulta
  * **When** inicio a cláusula `SELECT`
  * **Then** o sistema sugere colunas pertencentes à tabela selecionada.

## 3. Requisitos Funcionais (FR)
- **FR-03.1:** Suporte a sugestões contextuais de palavras-chave SQL.
- **FR-03.2:** Suporte a sugestões de tabelas e colunas com base na conexão ativa.
- **FR-03.3:** Geração de consulta final em formato editável no editor.
