package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	_ "modernc.org/sqlite"
)

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

type DBManager interface {
	SaveConnection(cfg ConnectionConfig) error
	ListConnections() ([]ConnectionConfig, error)
	DeleteConnection(id string) error
	GetDB(id string) (*sql.DB, error)
	GetTables(ctx context.Context, id string) ([]TableMetadata, error)
}

type defaultDBManager struct {
	mu         sync.RWMutex
	filePath   string
	activeDBs  map[string]*sql.DB
	configsMap map[string]ConnectionConfig
}

func NewManager() (DBManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	baseDir := filepath.Join(homeDir, ".goscratched")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base dir: %w", err)
	}

	mgr := &defaultDBManager{
		filePath:   filepath.Join(baseDir, "connections.json"),
		activeDBs:  make(map[string]*sql.DB),
		configsMap: make(map[string]ConnectionConfig),
	}

	if err := mgr.loadConfigs(); err != nil {
		mgr.configsMap = make(map[string]ConnectionConfig)
	}

	return mgr, nil
}

func (m *defaultDBManager) loadConfigs() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var configs []ConnectionConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return err
	}

	for _, cfg := range configs {
		m.configsMap[cfg.ID] = cfg
	}
	return nil
}

func (m *defaultDBManager) saveConfigsToFile() error {
	var configs []ConnectionConfig
	for _, cfg := range m.configsMap {
		configs = append(configs, cfg)
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Name < configs[j].Name
	})

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.filePath, data, 0644)
}

func (m *defaultDBManager) SaveConnection(cfg ConnectionConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.configsMap[cfg.ID] = cfg
	return m.saveConfigsToFile()
}

func (m *defaultDBManager) ListConnections() ([]ConnectionConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []ConnectionConfig
	for _, cfg := range m.configsMap {
		list = append(list, cfg)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}

func (m *defaultDBManager) DeleteConnection(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, ok := m.activeDBs[id]; ok {
		db.Close()
		delete(m.activeDBs, id)
	}

	delete(m.configsMap, id)
	return m.saveConfigsToFile()
}

func (m *defaultDBManager) GetDB(id string) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, ok := m.activeDBs[id]; ok {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		db.Close()
		delete(m.activeDBs, id)
	}

	cfg, ok := m.configsMap[id]
	if !ok {
		return nil, fmt.Errorf("connection config not found for id: %s", id)
	}

	driverName := string(cfg.Driver)
	if cfg.Driver == DriverSQLite {
		driverName = "sqlite"
	}

	db, err := sql.Open(driverName, cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	m.activeDBs[id] = db
	return db, nil
}

func (m *defaultDBManager) GetTables(ctx context.Context, id string) ([]TableMetadata, error) {
	db, err := m.GetDB(id)
	if err != nil {
		return nil, err
	}

	m.mu.RLock()
	cfg := m.configsMap[id]
	m.mu.RUnlock()

	switch cfg.Driver {
	case DriverSQLite:
		return getSQLiteTables(ctx, db)
	default:
		return getSQLiteTables(ctx, db)
	}
}

func getSQLiteTables(ctx context.Context, db *sql.DB) ([]TableMetadata, error) {
	rows, err := db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, fmt.Errorf("failed to query sqlite tables: %w", err)
	}
	defer rows.Close()

	var tables []TableMetadata
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}

		cols, _ := getSQLiteTableColumns(ctx, db, name)
		tables = append(tables, TableMetadata{
			Name:    name,
			Columns: cols,
		})
	}

	return tables, nil
}

func getSQLiteTableColumns(ctx context.Context, db *sql.DB, tableName string) ([]ColumnMetadata, error) {
	query := fmt.Sprintf("PRAGMA table_info(%q)", tableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []ColumnMetadata
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull int
		var dfltValue any
		var pk int

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk); err != nil {
			continue
		}

		cols = append(cols, ColumnMetadata{
			Name:         name,
			DataType:     dataType,
			IsPrimaryKey: pk > 0,
		})
	}

	return cols, nil
}
