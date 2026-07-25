# Especificação Funcional: GoScratch

## 1. Visão Geral & Objetivos
O **GoScratch** é uma aplicação desktop nativa desenvolvida em Go e Fyne inspirada nas melhores funcionalidades do **LINQPad**. Ele oferece um ambiente de scratchpad, REPL instantâneo, navegador de banco de dados e leitor de schemas com suporte a inspeção visual rica de dados.

## 2. Histórias de Usuário & Critérios de Aceite (BDD)

### US-01: Execução Instantânea de Snippets Soltos
**Como** desenvolvedor Go,  
**Quero** escrever trechos curtos de código sem digitar `package main` e `func main()`,  
**Para** testar lógicas e expressões em menos de 100ms.

* **Cenário 01.1: Código sem boilerplate**
  * **Given** que o editor do GoScratch contém `a := 10; b := 20; println(a + b)`
  * **When** eu aciono a execução (`Ctrl+R` ou botão "Executar")
  * **Then** o sistema adiciona o wrapper `package main` e `func main()` de forma transparente
  * **And** exibe `30` no painel de saída em menos de 100ms utilizando a engine Yaegi.

* **Cenário 01.2: Código completo fornecido pelo usuário**
  * **Given** que o editor contém uma declaração explícita de `package main` e `func main()`
  * **When** eu executo o código
  * **Then** o sistema respeita a estrutura fornecida pelo usuário sem aplicar wrappers adicionais.

---

### US-02: Inspeção Rica de Dados com `Dump(v any)`
**Como** desenvolvedor Go,  
**Quero** invocar a função global `Dump(v any)` em qualquer variável ou estrutura de dados,  
**Para** visualizar graficamente tipos complexos (structs, maps, slices, ponteiros) em uma árvore/tabela expansível na UI.

* **Cenário 02.1: Inspeção de Struct com campos aninhados**
  * **Given** que o código contém uma struct `type User struct { Name string; Age int }` e a chamada `Dump(User{Name: "Alice", Age: 30})`
  * **When** a execução é concluída
  * **Then** a aba "Inspeção (Dump)" exibe a árvore de dados navegável com os tipos, nomes de campos e valores.

* **Cenário 02.2: Resiliência contra referências circulares**
  * **Given** que uma estrutura possui um ponteiro que aponta para ela mesma ou cria um ciclo
  * **When** a função `Dump()` é executada
  * **Then** o dumper detecta o ciclo de ponteiros e interrompe a recursão com a notificação `[Circular Reference]`.

---

### US-03: Engine Híbrida e Fallback Inteligente
**Como** desenvolvedor Go,  
**Quero** utilizar recursos nativos avançados ou pacotes externos quando necessário,  
**Para** que meus scripts não fiquem limitados apenas aos recursos suportados pelo interpretador Yaegi.

* **Cenário 03.1: Fallback Automático para `go run`**
  * **Given** que o código utiliza um recurso não suportado ou instável no Yaegi (ex: cgo ou módulo complexo)
  * **When** a engine Yaegi falhar ao interpretar o script
  * **Then** o GoScratch alterna automaticamente para a engine secundária (`go run`)
  * **And** exibe um indicador visual informando que o script rodou no modo compilado (`go run`).

---

### US-04: Resiliência da Interface e Prevenção de Bloqueios
**Como** usuário da aplicação desktop,  
**Quero** que a interface continue responsiva mesmo se meu script entrar em um loop infinito,  
**Para** poder cancelar a execução a qualquer momento sem travar ou fechar o GoScratch.

* **Cenário 04.1: Timeout de Execução**
  * **Given** que um script contém um loop infinito `for {}`
  * **When** a execução atinge 5s (no Yaegi) ou 10s (no `go run`)
  * **Then** a execução é interrompida via contexto e um erro amigável "Timeout de Execução Excedido" é exibido na saída de erros.

---

### US-05: Menu em Árvore na Lateral Esquerda (Sidebar Explorer)
**Como** desenvolvedor,  
**Quero** visualizar um menu em árvore na barra lateral esquerda da aplicação,  
**Para** navegar facilmente entre meus Snippets de código salvos e minhas Conexões de Banco de Dados.

* **Cenário 05.1: Navegação de Snippets Salvos**
  * **Given** que possuo arquivos `.go` salvos no diretório de snippets
  * **When** expando o nó "📂 Meus Snippets" no menu lateral
  * **Then** o sistema lista os arquivos e subpastas existentes
  * **And** ao dar um duplo clique em um arquivo, o código é carregado imediatamente no editor.

---

### US-06: Gerenciador & Leitor de Schemas de Banco de Dados
**Como** desenvolvedor ou DBA,  
**Quero** cadastrar conexões de banco de dados (SQLite, PostgreSQL, MySQL) e expandir seus objetos,  
**Para** inspecionar bancos, esquemas, tabelas e colunas diretamente pela barra lateral.

* **Cenário 06.1: Expansão de Tabela e Colunas**
  * **Given** que tenho uma conexão de banco de dados cadastrada
  * **When** expando o nó da conexão -> Banco -> Tabela "users"
  * **Then** o menu em árvore exibe a lista de colunas com seus respectivos tipos SQL e indicação de Chave Primária (PK).

---

### US-07: Auto-Dump de Consultas SQL (`sql.Rows`)
**Como** desenvolvedor,  
**Quero** passar um objeto `*sql.Rows` ou executar uma query SQL diretamente no `Dump()`,  
**Para** visualizar o resultado da consulta formatado em uma tabela/grid expansível com nome de colunas e dados.

* **Cenário 07.1: Inspeção de `sql.Rows` com `Dump(rows)`**
  * **Given** que executei uma consulta `rows, _ := db.Query("SELECT id, name FROM users")`
  * **When** invoco `Dump(rows)`
  * **Then** o `pkg/dumper` consome os registros e gera um nó de inspeção em formato de tabela com os nomes das colunas como cabeçalhos e as linhas como registros filhos.

* **Cenário 07.2: Duplo Clique em Tabela da Barra Lateral**
  * **Given** que estou navegando na árvore de conexões de banco de dados
  * **When** clico duas vezes em uma tabela (ex: "customers")
  * **Then** o GoScratch insere automaticamente o snippet de conexão e consulta SQL no editor e executa a busca com `Dump()`.

---

## 3. Requisitos Funcionais (FR)

| ID | Requisito | Descrição |
|---|---|---|
| **FR-01** | **Wrapping de Código** | Inserção automática de `package main` e `func main()` caso o código do usuário não possua. |
| **FR-02** | **Função `Dump` Injetada** | Disponibilização global da função `Dump(v any)` sem necessidade de imports explicitados pelo usuário. |
| **FR-03** | **Engine Primária (Yaegi)** | Execução por padrão via interpretador Yaegi para resposta instantânea (< 100ms). |
| **FR-04** | **Engine Secundária (`go run`)** | Fallback para compilação oficial `go run` em diretório temporário quando Yaegi falhar ou for explicitamente selecionado. |
| **FR-05** | **Painel Split & Abas** | UI Fyne dividida com Sidebar em Árvore na Esquerda, Editor no Centro e Painéis com Abas ("Console" e "Inspeção Dump") no Canto Direito. |
| **FR-06** | **Navegador de Dados Dump** | Renderização de structs, maps, slices e tipos primitivos em widgets expansíveis (Tree / Accordion). |
| **FR-07** | **Timeout Configurável** | Cancelamento automático de execuções após 5s (Yaegi) ou 10s (`go run`) usando `context.Context`. |
| **FR-08** | **Destaque de Erros** | Exibição de erros de compilação, sintaxe e tempo de execução em destaque visual na aba Console. |
| **FR-09** | **Botão de Interrupção** | Possibilidade de cancelar manualmente uma execução em andamento. |
| **FR-10** | **Suporte a Atalhos** | Atalho global `Ctrl+Enter` ou `Ctrl+R` para iniciar a execução do código. |
| **FR-11** | **Menu em Árvore Lateral** | Sidebar expansível com visualização em árvore para Snippets e Conexões de Banco. |
| **FR-12** | **Gerenciador de Conexões DB** | Suporte a driver SQLite (modernc sem CGO), PostgreSQL e MySQL. |
| **FR-13** | **Leitor de Schemas SQL** | Obtenção de metadados de tabelas e colunas para exibição na árvore. |
| **FR-14** | **Auto-Dump de `sql.Rows`** | Tratamento especial de `*sql.Rows` no `pkg/dumper` convertendo dados tabulares em `DumpNode`. |
| **FR-15** | **Gerador de Query Rápida** | Duplo clique em tabela da sidebar insere código Go de consulta pronto com `Dump()`. |

---

## 4. Requisitos Não-Funcionais (NFR)

| ID | Requisito | Critério de Aceite |
|---|---|---|
| **NFR-01** | **Latência de Execução** | Resposta visual no Yaegi em menos de **100ms** para snippets básicos. |
| **NFR-02** | **Consumo de Memória** | Aplicação UI (Fyne) consumindo menos de **75MB RAM** em repouso com drivers DB ativos. |
| **NFR-03** | **Independência de Runtimes Web** | Binário desktop nativo único, sem consumo de Electron/Chromium/Node.js. |
| **NFR-04** | **Desacoplamento de Pacotes** | Zero dependências diretas de Fyne nos pacotes `pkg/runner`, `pkg/dumper` e `pkg/db`. |
| **NFR-05** | **Concorrência Segura** | Execução de código e queries 100% isoladas da Main Loop do Fyne através de goroutines dedicadas. |
