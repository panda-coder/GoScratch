# Desenho Técnico: Database Explorer & Schema Inspector

## Arquitetura do Pacote `pkg/db`

```go
package db

type DriverType string

const (
	DriverSQLite   DriverType = "sqlite"
	DriverPostgres DriverType = "postgres"
	DriverMySQL    DriverType = "mysql"
)

type ConnectionConfig struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Driver  DriverType `json:"driver"`
	ConnStr string     `json:"conn_str"`
}

type ColumnMetadata struct {
	Name         string `json:"name"`
	DataType     string `json:"data_type"`
	IsPrimaryKey bool   `json:"is_primary_key"`
}

type TableMetadata struct {
	Name    string           `json:"name"`
	Columns []ColumnMetadata `json:"columns"`
}

type Explorer interface {
	ListTables(conn ConnectionConfig) ([]TableMetadata, error)
	ListColumns(conn ConnectionConfig, table string) ([]ColumnMetadata, error)
}
```
