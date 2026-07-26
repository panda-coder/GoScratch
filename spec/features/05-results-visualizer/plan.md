# Desenho Técnico: Results Visualizer

## Arquitetura do Pacote `pkg/results`

```go
package results

type Row map[string]any

type ResultSet struct {
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
}

type Visualizer interface {
	Load(result ResultSet) error
	SelectedRow() Row
}
```
