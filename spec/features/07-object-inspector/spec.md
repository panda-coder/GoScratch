# Especificação Funcional: Object Inspector

## 1. Visão Geral
Esta funcionalidade provê a inspeção navegável de objetos no GoScratch. Ela permite explorar estruturas complexas em árvore, entendendo campos, tipos e valores de forma visual e interativa.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-07.1: Navegação em Estruturas Complexas
**Como** desenvolvedor analisando estruturas em memória,  
**Quero** expandir e recolher objetos complexos em uma árvore visual,  
**Para** compreender rapidamente a composição interna de structs, coleções e valores aninhados.

* **Cenário 07.1.1: Expansão de struct aninhada**
  * **Given** que um objeto possui campos compostos e aninhados
  * **When** solicito a inspeção do objeto
  * **Then** o sistema exibe a hierarquia navegável com nós para cada campo relevante.

* **Cenário 07.1.2: Visualização de tipo e valor do nó**
  * **Given** que um nó da árvore representa um campo inspecionado
  * **When** visualizo esse nó na interface
  * **Then** o sistema apresenta nome, tipo e valor formatado do elemento correspondente.

## 3. Requisitos Funcionais (FR)
- **FR-07.1:** Suporte à inspeção de structs, mapas, slices e valores primitivos.
- **FR-07.2:** Exibição hierárquica em árvore navegável.
- **FR-07.3:** Exibição de nome, tipo e valor para cada nó inspecionado.
