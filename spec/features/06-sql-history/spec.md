# Especificação Funcional: SQL History

## 1. Visão Geral
Esta funcionalidade provê o histórico de consultas SQL no GoScratch. Ela permite registrar execuções anteriores, revisitar comandos recentes e reaproveitar consultas sem reescrita manual.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-06.1: Reexecução de Consultas Recentes
**Como** desenvolvedor trabalhando iterativamente com SQL,  
**Quero** revisar e reexecutar consultas já utilizadas,  
**Para** acelerar testes, comparações e ajustes em sequência.

* **Cenário 06.1.1: Registro automático após execução**
  * **Given** que executo uma consulta SQL válida
  * **When** a execução é finalizada
  * **Then** o sistema adiciona a consulta ao histórico com data e origem.

* **Cenário 06.1.2: Reaplicação de item do histórico**
  * **Given** que existe ao menos uma consulta registrada no histórico
  * **When** seleciono uma entrada anterior
  * **Then** o sistema reinsere a consulta no editor para nova execução ou edição.

## 3. Requisitos Funcionais (FR)
- **FR-06.1:** Registro automático de consultas executadas.
- **FR-06.2:** Listagem cronológica de entradas do histórico.
- **FR-06.3:** Reaplicação de entradas no editor principal.
