🚀 GoScratch
O bloco de rascunho de alta performance para desenvolvedores Go.

Teste ideias, execute trechos de código em milissegundos e inspecione dados complexos sem o boilerplate de criar arquivos ou gerenciar projetos temporários.

💡 Sobre o Projeto
O GoScratch é uma ferramenta Desktop leve e nativa (feita com Fyne) projetada para trazer a experiência de scratchpad e REPL instantâneo para o ecossistema Go.

Inspirado no fluxo de trabalho de inspeção rápida do LINQPad, o GoScratch elimina a necessidade de criar pastas temp, configurar arquivos main.go ou lidar com compilações lentas para testar uma única função, validar uma lógica de manipulação de dados ou formatar um JSON.

✨ Principais Funcionalidades
⚡ Feedback Instantâneo (< 100ms): Alimentado pelo interpretador Yaegi, o GoScratch executa seus snippets de código instantaneamente, sem o overhead de compilação do go build.

🔍 Inspecionabilidade First-Class com Dump(): Esqueça o fmt.Printf("%+v\n"). A função global Dump(v any) intercepta structs, maps, slices e primitivos, renderizando tabelas e árvores de dados ricas e expansíveis na interface.

⚙️ Engine Híbrida Inteligente: Executa scripts dinâmicos rapidamente via Yaegi, mas conta com fallback automático para o compilador oficial (go run) quando você precisa de recursos nativos avançados ou pacotes externos do go.mod.

🎨 UI Nativa & Leve: Desenvolvido em Go puro com Fyne, sem Electron, sem Webview e sem consumo excessivo de RAM.

🛡️ Execução Segura & Isolada: Interface responsiva que nunca trava durante a execução de scripts pesados, com suporte a timeouts e cancelamento via contexto.

🛠️ Stack Tecnológica
Linguagem: Go 1.26+

Interface Gráfica: Fyne v2

Interpretador REPL: Yaegi

Testes Unitários: Testify (Suite)

🎯 Para quem é o GoScratch?
Para quem quer testar uma regex, um algoritmo ou uma struct sem poluir o projeto atual.

Para quem precisa inspecionar retornos de APIs ou payloads JSON complexos de forma visual.

Para quem quer um ambiente de aprendizado e prototipagem rápida no dia a dia com Go.
