# Desenho Técnico: Redis Connection & Key Inspector

## Arquitetura do Pacote `pkg/redis`

```go
package redis

type ConnectionConfig struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	Password string `json:"password,omitempty"`
	DB       int    `json:"db"`
}

type KeyMetadata struct {
	Name string `json:"name"`
	Type string `json:"type"`
	TTL  int64  `json:"ttl"`
}

type Inspector interface {
	ListKeys(conn ConnectionConfig, pattern string) ([]KeyMetadata, error)
	InspectKey(conn ConnectionConfig, key string) (any, error)
}
```
