# Desenho Técnico: Snippet Library

## Arquitetura do Pacote `pkg/snippets`

```go
package snippets

type Snippet struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Code     string   `json:"code"`
	Tags     []string `json:"tags,omitempty"`
	Language string   `json:"language"`
}

type Repository interface {
	Save(snippet Snippet) error
	List() ([]Snippet, error)
	Find(query string) ([]Snippet, error)
}
```
