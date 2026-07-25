# 🚀 Visão Geral de Features — GoScratch

Esta pasta contém o conjunto de especificações funcionais e técnicas que guiam o desenvolvimento do **GoScratch** sob o fluxo *Spec-Driven Development (SDD)*. 

Cada subdiretório representa uma funcionalidade principal e contém o trio de controle de entrega: `spec.md` (Requisitos), `plan.md` (Arquitetura) e `tasks.md` (Checklist).

---

## 📋 Lista de Features do GoScratch

### 🎯 Core MVP (Fase 1)

#### `01-instant-repl-dumper`
* **Descrição:** Editor de código interativo com execução instantânea e painel de inspeção rica de dados.
* **Mecanismos Principais:**
  * Injeção global transparente da função `Dump(v any)`.
  * Visualizador em árvore/tabela expansível na UI (Fyne) para `structs`, `maps`, `slices` e primitivos.
  * Captura assíncrona de `stdout`, `stderr` e erros de compilação sem travar a Main Loop.

#### `02-database-explorer-sql`
* **Descrição:** Gerenciador de conexões de banco de dados e gerador automático de estruturas de dados.
* **Mecanismos Principais:**
  * Suporte inicial para drivers nativos de Go (PostgreSQL, MySQL, SQLite).
  * Painel lateral de navegação em esquemas, tabelas e colunas.
  * Auto-geração de `structs` Go com tags JSON/SQL a partir de tabelas do banco.
  * Execução rápida de consultas SQL diretamente no editor e renderização dos resultados via `Dump()`.

---

### ⚡ Performance & Execução (Fase 2)

#### `03-hybrid-execution-engine`
* **Descrição:** Motor híbrido de execução que combina velocidade instantânea com 100% de suporte à linguagem Go.
* **Mecanismos Principais:**
  * Engine Primária (Yaegi): Interpretação AST em milissegundos (< 100ms) para scripts leves e testes rápidos.
  * Engine Secundária (`go run` via `os/exec`): Fallback automático ou alternância manual para suportar Generics complexos, CGO ou recursos nativos de versões específicas do Go.
  * Gestão de tempo de execução (*timeout*) e cancelamento gracioso via `context.Context`.

#### `04-package-manager-gomod`
* **Descrição:** Gerenciamento dinâmico de dependências externas sem necessidade de setup manual de projeto.
* **Mecanismos Principais:**
  * Inspecção de bloco `import (...)` no snippet do usuário.
  * Execução em background de `go get` / `go mod tidy` em diretórios temporários na memória (`RAM/tmpfs`).
  * Autocompletar e cache de pacotes baixados do `pkg.go.dev`.

---

### 🖥️ Produtividade & DX (Fase 3)

#### `05-multi-tab-workspace`
* **Descrição:** Ambiente de trabalho com múltiplas abas, histórico de execuções e persistência de trechos de código.
* **Mecanismos Principais:**
  * Sistema de abas independentes (cada aba possui seu próprio contexto de execução e dumper).
  * Auto-salvamento e histórico de *scratchpads* recentes.
  * Gestão de variáveis de ambiente (`.env`) por aba de execução.

#### `06-export-and-code-generators`
* **Descrição:** Utilitários para converter o código rascunhado e seus resultados para uso no projeto final.
* **Mecanismos Principais:**
  * Conversor de payloads (JSON para `struct` Go / SQL Schema para Go GORM/Ent/sqlx).
  * Exportação de dados do `Dump()` para formatos `JSON`, `CSV`, `YAML` e `Markdown Table`.
  * Botão de conversão automática de "Snippet GoScratch" para "Projeto Go Completo" (injetando `package main` e `func main()`).

#### `07-benchmarking-profiler`
* **Descrição:** Medição rápida de performance de código Go e consumo de memória.
* **Mecanismos Principais:**
  * Suporte nativo para funções no estilo `BenchmarkXxx(b *testing.B)`.
  * Exibição visual de tempo por operação (`ns/op`), alocações de memória (`B/op`) e alocações por ciclo (`allocs/op`).