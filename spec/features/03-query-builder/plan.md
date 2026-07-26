# Desenho Técnico: Query Builder

## Arquitetura do Pacote `pkg/querybuilder`

```go
package querybuilder

type SuggestionKind string

const (
	KindKeyword SuggestionKind = "KEYWORD"
	KindTable   SuggestionKind = "TABLE"
	KindColumn  SuggestionKind = "COLUMN"
	KindClause  SuggestionKind = "CLAUSE"
)

type Suggestion struct {
	Label string         `json:"label"`
	Value string         `json:"value"`
	Kind  SuggestionKind `json:"kind"`
}

type Builder interface {
	Suggest(input string) []Suggestion
	Build(parts []string) string
}
```
