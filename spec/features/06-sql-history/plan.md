# Desenho Técnico: SQL History

## Arquitetura do Pacote `pkg/history`

```go
package history

type Entry struct {
	Query     string `json:"query"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

type Store interface {
	Append(entry Entry) error
	List() ([]Entry, error)
	Clear() error
}
```
