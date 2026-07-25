# Checklist de Tarefas: GoScratch

Este plano de tarefas divide a implementação do GoScratch em etapas pequenas, isoladas e testáveis, organizadas por pacotes.

---

## Fase 1: Motor de Inspeção (`pkg/dumper`)

- [x] **1.1. Estruturas Base & Contratos (`pkg/dumper`)**
  - Definir o enum `NodeKind` (`KindPrimitive`, `KindStruct`, `KindMap`, `KindSlice`, etc.).
  - Definir a struct `DumpNode` e a interface `Dumper`.

- [x] **1.2. Refletor de Primitivos e Tipos Simples**
  - Implementar reflexão básica para `int`, `float`, `string`, `bool` e `error`.
  - Escrever testes unitários em `pkg/dumper/dumper_test.go` usando `testify/suite`.

- [x] **1.3. Refletor de Coleções (Slices, Arrays, Maps)**
  - Implementar navegação recursiva por elementos de slices/arrays e pares chave-valor de maps.
  - Testar casos de mapas vazios, nulos e aninhados em `dumper_test.go`.

- [x] **1.4. Refletor de Structs e Ponteiros com Prevenção de Ciclos**
  - Implementar inspeção de campos exportados e não-exportados de structs.
  - Implementar mapa de ponteiros visitados para detectar referências circulares (`KindCircular`).
  - Testar structs complexas e referências circulares em `dumper_test.go`.

---

## Fase 2: Motor de Execução (`pkg/runner`)

- [x] **2.1. Estruturas Base & Contrato `Runner`**
  - Definir `ExecutionResult`, enum `Mode` (`ModeYaegi`, `ModeGoRun`) e a interface `Runner`.

- [x] **2.2. Detector & Injetor de Boilerplate (Code Wrapper)**
  - Implementar função utilitária `WrapCode(code string) string` que detecta ausência de `package main` / `func main()` e adiciona o wrapper automaticamente.
  - Testar cenários com e sem boilerplate.

- [x] **2.3. Integração com Interpretador Yaegi (`ModeYaegi`)**
  - Configurar `interp.New()` com captura de `stdout` e `stderr` via `bytes.Buffer`.
  - Injetar a função nativa `Dump(v any)` no escopo do Yaegi integrando com `pkg/dumper`.
  - Adicionar suporte a cancelamento por `context.Context`.
  - Escrever testes unitários em `pkg/runner/runner_test.go` usando `testify/suite`.

- [x] **2.4. Engine Secundária (`ModeGoRun` - Fallback)**
  - Implementar execução via `os/exec.CommandContext` chamando `go run` em diretório temporário.
  - Implementar mecanismo de fallback automático: tentar Yaegi primeiro; se falhar por limitação de linguagem, rodar `go run`.
  - Testar fallback em `runner_test.go`.

---

## Fase 3: Interface Gráfica Fyne (`pkg/ui`)

- [x] **3.1. Estrutura de Janela & Layout Split**
  - Configurar janela principal do Fyne v2 com `container.NewHSplit`.
  - Adicionar editor multi-linha com fonte monospace do lado esquerdo.

- [x] **3.2. Painel de Saída & Abas (Console & Dump)**
  - Construir `container.NewAppTabs` no lado direito.
  - Criar aba "Console" exibindo `stdout` e `stderr` com estilo visual diferenciado para erros.
  - Criar aba "Inspeção (Dump)" usando widget `widget.Tree` customizado para renderizar a árvore de `DumpNode`.

- [x] **3.3. Concorrência & Indicadores de Estado**
  - Executar a chamada de `runner.Execute` em goroutines separadas para não travar a UI.
  - Atualizar componentes da UI com `fyne.Do()`.
  - Adicionar barra de status exibindo tempo de execução (ms), engine utilizada (`YAEGI` vs `GORUN`) e botão de cancelamento.
  - Adicionar atalho `Ctrl+Enter` / `Ctrl+R` para disparo de execução.

---

## Fase 4: Integração Final & Empacotamento

- [x] **4.1. Teste de Integração End-to-End**
  - Executar trechos de código reais no GoScratch (ex: manipulação de JSON, structs pesadas, loops com timeout).
- [x] **4.2. Build do Binário Nativo**
  - Validar compilação limpa com `go build -o goscratched main.go` e verificar tamanho e consumo de memória.

---

## Fase 5: Módulo de Banco de Dados & Auto-Dump de SQL (`pkg/db`)

- [x] **5.1. Suporte a `*sql.Rows` no Dumper (`pkg/dumper`)**
  - Adicionar inspeção automática de `*sql.Rows` convertendo o resultado tabular de queries para `DumpNode`.
  - Escrever testes unitários em `pkg/dumper/dumper_test.go` com SQLite em memória.

- [x] **5.2. Gerenciador de Conexões & Schemas (`pkg/db`)**
  - Criar `ConnectionConfig`, `DBManager` e persistência de conexões em `~/.goscratched/connections.json`.
  - Implementar leitor de metadados `GetTables` e `GetColumns` para SQLite, PostgreSQL e MySQL.
  - Escrever testes unitários em `pkg/db/db_test.go`.

---

## Fase 6: Biblioteca de Snippets Salvos (`pkg/snippets`)

- [x] **6.1. Gerenciador de Arquivos (`pkg/snippets`)**
  - Implementar leitura e gravação de arquivos `.go` no diretório `~/.goscratched/snippets/`.
  - Escrever testes unitários em `pkg/snippets/snippets_test.go`.

---

## Fase 7: Menu em Árvore Lateral Fyne (`pkg/ui/sidebar.go`)

- [x] **7.1. Árvore Explorer (`pkg/ui/sidebar.go`)**
  - Criar componente de barra lateral esquerda com `widget.Tree` exibindo:
    - **📂 Snippets Salvos** (arquivos e subpastas).
    - **🗄️ Conexões de Banco de Dados** (Conexões -> Tabelas -> Colunas com PKs).
  - Adicionar botão de "Nova Conexão" na barra lateral.

- [x] **7.2. Interatividade & Atalhos**
  - Clique duplo em Snippet: Carrega código no editor.
  - Clique duplo em Tabela: Inserção e execução automática de snippet SQL com `Dump(rows)`.
- [x] **7.3. Teste de Integração e Compilação Final**
  - Rodar `go test ./...` e `go build -o goscratched main.go`.
