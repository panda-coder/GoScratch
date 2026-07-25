# Constituição do Projeto: GoScratch

## 1. Princípios Fundamentais
* **Feedback Instantâneo:** A execução de trechos de código simples deve responder visualmente em menos de 100ms.
* **Inspecionabilidade First-Class:** A função `Dump(v any)` é o cidadão de primeira classe da interface. Ela deve interceptar qualquer estrutura de dados (structs, maps, slices, primitives) e renderizar uma visualização rica e expansível na UI.
* **Simplicidade de Uso:** O usuário não deve ser obrigado a escrever boilerplate (`package main`, `func main()`) para testar trechos rápidos de código.
* **Portabilidade & Binário Único:** O projeto deve compilar em um binário desktop leve e nativo sem dependências de ambientes web (sem Webview, Chromium ou Node.js).

## 2. Tecnologias & Arquitetura Core
* **Linguagem:** Go (utilizando recursos e idioms modernos do Go 1.26+).
* **Interface Gráfica (UI):** **Fyne (fyne.io/fyne/v2)**
  * Layout responsivo com painéis divididos (Editor de Código de um lado / Painel de Saída e Dump do outro).
  * Renderização de componentes customizados do Fyne para exibir as tabelas/árvores geradas pelo `.Dump()`.
* **Execution Engine (Híbrida):**
  * **Engine Primária (Padrão):** **Yaegi** (para execução instantânea, REPL e injeção transparente do `Dump`).
  * **Engine Secundária (Fallback / Modo Avançado):** **`go run` via `os/exec`** com diretório temporário em memória (para suporte a 100% da linguagem, CGO ou módulos `go.mod` complexos).

## 3. Diretrizes de Design & Separação de Conceitos
* **Desacoplamento Rigoroso:**
  * `pkg/runner`: Lógica pura de execução e captura de `stdout`/`stderr`. Não conhece a UI.
  * `pkg/dumper`: Lógica de reflexão e conversão de dados em estruturas intermediárias (ex: nós/árvores/tabelas). Não conhece o Fyne.
  * `pkg/ui`: Camada de apresentação Fyne. Apenas consome o resultado do `dumper` e exibe os widgets correspondentes.
* **Resiliência de UI:** O processo de execução do código do usuário nunca deve travar a *Main Loop* da UI do Fyne. Toda execução de código deve rodar em goroutines separadas com suporte a `context.Context` para cancelamento/timeout.

## 4. Estilo de Código & Qualidade
* Código 100% idiomático em Go, formatado com `gofmt` / `goimports`.
* Tratamento de erros explícito e amigável para o usuário (erros de compilação/execução devem ser destacados visualmente no painel de saída).
* **Padrão de Testes:** Testes unitários são obrigatórios para os pacotes `runner` e `dumper`, devendo utilizar a estrutura de suítes do pacote **`github.com/stretchr/testify/suite`** (`testify/assert` e `testify/require`) para garantir legibilidade, reuso de rotinas de setup/teardown e asserções padronizadas.
