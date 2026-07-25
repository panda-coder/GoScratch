# Especificação Funcional: Instant REPL & Rich Data Dump

## 1. Visão Geral
Esta funcionalidade provê o núcleo interativo do GoScratch (equivalente ao famoso `Dump()` do LINQPad). Ela permite a execução instantânea de trechos de código em menos de 100ms e a inspeção visual rica de qualquer estrutura de dados em Go.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-01.1: Inspeção Visual de Dados com `Dump(v any)`
**Como** desenvolvedor Go,  
**Quero** chamar a função `Dump(v any)` em qualquer expressão ou variável,  
**Para** navegar graficamente em uma árvore expansível com tipos, campos e valores formatados.

* **Cenário 01.1.1: Inspeção de Struct com Campos Aninhados**
  * **Given** que declaro a struct `type User struct { Name string; Age int }`
  * **When** invoco `Dump(User{Name: "Alice", Age: 30})`
  * **Then** o painel de Inspeção exibe um nó do tipo `KindStruct` com os filhos `Name` e `Age`.

* **Cenário 01.1.2: Prevenção de Referências Circulares**
  * **Given** que uma struct possui um ponteiro auto-referenciável
  * **When** invoco `Dump()` sobre a estrutura
  * **Then** a recursão é interrompida com o aviso `[Circular Reference]` sem estouro de memória.

## 3. Requisitos Funcionais (FR)
- **FR-01.1:** Suporte à reflexão de primitivos, structs, maps, slices, interfaces e ponteiros.
- **FR-01.2:** Tratamento amigável para os tipos `error` e `fmt.Stringer`.
- **FR-01.3:** Renderização em widget de árvore expansível (`widget.Tree`).
