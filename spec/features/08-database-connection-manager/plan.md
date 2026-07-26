# Desenho Técnico: Database Connection Manager

## Arquitetura do Pacote `pkg/connections`

```go
package connections

type Connection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	DSN      string `json:"dsn"`
	IsActive bool   `json:"is_active"`
}

type Manager interface {
	Save(conn Connection) error
	List() ([]Connection, error)
	SetActive(id string) error
	Test(conn Connection) error
}
```
