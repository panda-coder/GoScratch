package db_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/panda-coder/GoScratch/pkg/db"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type DBTestSuite struct {
	suite.Suite
	mgr db.DBManager
}

func (s *DBTestSuite) SetupTest() {
	mgr, err := db.NewManager()
	s.Require().NoError(err)
	s.mgr = mgr
}

func (s *DBTestSuite) TestSaveAndListConnections() {
	cfg := db.ConnectionConfig{
		ID:      "sqlite_test",
		Name:    "SQLite Test DB",
		Driver:  db.DriverSQLite,
		ConnStr: ":memory:",
	}

	err := s.mgr.SaveConnection(cfg)
	s.NoError(err)

	list, err := s.mgr.ListConnections()
	s.NoError(err)
	s.NotEmpty(list)

	found := false
	for _, c := range list {
		if c.ID == "sqlite_test" && c.Name == "SQLite Test DB" {
			found = true
			break
		}
	}
	s.True(found, "expected sqlite_test connection to be found in list")
}

func (s *DBTestSuite) TestGetTablesAndColumns() {
	tempDir := s.T().TempDir()
	dbFile := filepath.Join(tempDir, "test.db")

	rawDB, err := sql.Open("sqlite", dbFile)
	s.Require().NoError(err)
	_, err = rawDB.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, email TEXT NOT NULL)")
	s.Require().NoError(err)
	rawDB.Close()

	cfg := db.ConnectionConfig{
		ID:      "sqlite_file",
		Name:    "SQLite File DB",
		Driver:  db.DriverSQLite,
		ConnStr: dbFile,
	}
	s.Require().NoError(s.mgr.SaveConnection(cfg))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tables, err := s.mgr.GetTables(ctx, "sqlite_file")
	s.NoError(err)
	s.Require().Len(tables, 1)
	s.Equal("users", tables[0].Name)

	cols := tables[0].Columns
	s.Require().Len(cols, 2)

	s.Equal("id", cols[0].Name)
	s.True(cols[0].IsPrimaryKey)

	s.Equal("email", cols[1].Name)
	s.False(cols[1].IsPrimaryKey)
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
