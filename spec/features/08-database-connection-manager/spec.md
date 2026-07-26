# Especificação Funcional: Database Connection Manager

## 1. Visão Geral
Esta funcionalidade provê o gerenciamento de conexões de banco de dados no GoScratch. Ela permite cadastrar, testar, selecionar e reutilizar conexões em recursos como exploração SQL e construção de consultas.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-08.1: Cadastro e Seleção de Conexões
**Como** desenvolvedor trabalhando com múltiplos bancos,  
**Quero** cadastrar e selecionar conexões reutilizáveis,  
**Para** alternar rapidamente entre ambientes sem reconfiguração manual a cada uso.

* **Cenário 08.1.1: Cadastro de uma nova conexão**
  * **Given** que preciso acessar um banco ainda não configurado
  * **When** informo nome, driver e string de conexão válidos
  * **Then** o sistema persiste a nova conexão para reutilização futura.

* **Cenário 08.1.2: Definição de conexão ativa**
  * **Given** que existem múltiplas conexões salvas
  * **When** seleciono uma conexão como ativa
  * **Then** os recursos SQL passam a utilizá-la como padrão nas execuções seguintes.

## 3. Requisitos Funcionais (FR)
- **FR-08.1:** Cadastro e persistência de conexões nomeadas.
- **FR-08.2:** Teste de conexão antes de uso ou salvamento.
- **FR-08.3:** Seleção de conexão ativa para recursos dependentes de banco.
