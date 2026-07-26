# Desenho Técnico: Object Inspector

## Arquitetura do Pacote `pkg/inspector`

```go
package inspector

type Node struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Value    string  `json:"value"`
	Children []*Node `json:"children,omitempty"`
}

type Inspector interface {
	Inspect(v any) *Node
}
```
