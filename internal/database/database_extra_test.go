package database

import (
	"testing"
	"time"
)

func TestMySQLFullLifecycle(t *testing.T) {
	db := NewMySQL()
	if db.Name() != "mysql" {
		t.Errorf("expected 'mysql', got %s", db.Name())
	}

	if err := db.Connect(&Config{Host: "localhost", Port: 3306, Database: "test"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}

	result, err := db.Exec("CREATE TABLE t (id INT)")
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if result.RowsAffected != 1 {
		t.Errorf("expected 1 row affected, got %d", result.RowsAffected)
	}

	rows, err := db.Query("SELECT * FROM t")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if rows == nil {
		t.Error("expected non-nil rows")
	}

	row, err := db.QueryRow("SELECT * FROM t WHERE id = 1")
	if err != nil {
		t.Fatalf("queryrow: %v", err)
	}
	if row == nil {
		t.Error("expected non-nil row")
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if _, err := tx.Exec("INSERT INTO t VALUES (1)"); err != nil {
		t.Fatalf("tx exec: %v", err)
	}
	txRows, err := tx.Query("SELECT * FROM t")
	if err != nil {
		t.Fatalf("tx query: %v", err)
	}
	if txRows == nil {
		t.Error("expected non-nil tx rows")
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("tx commit: %v", err)
	}

	if err := db.Migrate([]Migration{
		{Version: 1, Name: "init", Up: "CREATE TABLE IF NOT EXISTS _m(id INT)"},
		{Version: 2, Name: "add_col", Up: "ALTER TABLE _m ADD COLUMN name TEXT"},
	}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if db.MigrationVersion() != 2 {
		t.Errorf("expected version 2, got %d", db.MigrationVersion())
	}

	if err := db.Rollback(1); err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if db.MigrationVersion() != 1 {
		t.Errorf("expected version 1 after rollback, got %d", db.MigrationVersion())
	}

	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestMySQLNotConnected(t *testing.T) {
	db := NewMySQL()

	if _, err := db.Exec("SELECT 1"); err == nil {
		t.Error("expected error")
	}
	if _, err := db.Query("SELECT 1"); err == nil {
		t.Error("expected error")
	}
	if _, err := db.QueryRow("SELECT 1"); err == nil {
		t.Error("expected error")
	}
	if _, err := db.Begin(); err == nil {
		t.Error("expected error")
	}
	if err := db.Migrate(nil); err == nil {
		t.Error("expected error")
	}
	if err := db.Rollback(0); err == nil {
		t.Error("expected error")
	}
	if err := db.Ping(); err == nil {
		t.Error("expected error")
	}
}

func TestMySQLTransactionRollback(t *testing.T) {
	db := NewMySQL()
	db.Connect(&Config{})

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if _, err := tx.Exec("INSERT INTO t VALUES (1)"); err != nil {
		t.Fatalf("tx exec: %v", err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatalf("tx rollback: %v", err)
	}
}

func TestSQLiteFullLifecycle(t *testing.T) {
	db := NewSQLite()
	if db.Name() != "sqlite" {
		t.Errorf("expected 'sqlite', got %s", db.Name())
	}

	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}

	result, err := db.Exec("CREATE TABLE t (id INT)")
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if result.RowsAffected != 1 {
		t.Errorf("expected 1 row affected, got %d", result.RowsAffected)
	}

	rows, err := db.Query("SELECT * FROM t")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if rows == nil {
		t.Error("expected non-nil rows")
	}

	row, err := db.QueryRow("SELECT * FROM t WHERE id = 1")
	if err != nil {
		t.Fatalf("queryrow: %v", err)
	}
	if row == nil {
		t.Error("expected non-nil row")
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if _, err := tx.Exec("INSERT INTO t VALUES (1)"); err != nil {
		t.Fatalf("tx exec: %v", err)
	}
	txRows, err := tx.Query("SELECT * FROM t")
	if err != nil {
		t.Fatalf("tx query: %v", err)
	}
	if txRows == nil {
		t.Error("expected non-nil tx rows")
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("tx commit: %v", err)
	}

	if err := db.Migrate([]Migration{
		{Version: 1, Name: "init", Up: "CREATE TABLE IF NOT EXISTS _m(id INT)"},
		{Version: 2, Name: "add_col", Up: "ALTER TABLE _m ADD COLUMN name TEXT"},
	}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if db.MigrationVersion() != 2 {
		t.Errorf("expected version 2, got %d", db.MigrationVersion())
	}

	if err := db.Rollback(1); err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if db.MigrationVersion() != 1 {
		t.Errorf("expected version 1 after rollback, got %d", db.MigrationVersion())
	}

	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestSQLiteNotConnected(t *testing.T) {
	db := NewSQLite()

	if _, err := db.Exec("SELECT 1"); err == nil {
		t.Error("expected error")
	}
	if _, err := db.Query("SELECT 1"); err == nil {
		t.Error("expected error")
	}
	if _, err := db.QueryRow("SELECT 1"); err == nil {
		t.Error("expected error")
	}
	if _, err := db.Begin(); err == nil {
		t.Error("expected error")
	}
	if err := db.Migrate(nil); err == nil {
		t.Error("expected error")
	}
	if err := db.Rollback(0); err == nil {
		t.Error("expected error")
	}
	if err := db.Ping(); err == nil {
		t.Error("expected error")
	}
}

func TestSQLiteTransactionRollback(t *testing.T) {
	db := NewSQLite()
	db.Connect(&Config{Database: ":memory:"})

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if _, err := tx.Exec("INSERT INTO t VALUES (1)"); err != nil {
		t.Fatalf("tx exec: %v", err)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatalf("tx rollback: %v", err)
	}
}

func TestPostgreSQLRollback(t *testing.T) {
	db := NewPostgreSQL()
	db.Connect(&Config{})

	if err := db.Migrate([]Migration{
		{Version: 1, Name: "v1", Up: "CREATE TABLE t1(id INT)"},
		{Version: 2, Name: "v2", Up: "CREATE TABLE t2(id INT)"},
		{Version: 3, Name: "v3", Up: "CREATE TABLE t3(id INT)"},
	}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if db.MigrationVersion() != 3 {
		t.Fatalf("expected version 3, got %d", db.MigrationVersion())
	}

	if err := db.Rollback(1); err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if db.MigrationVersion() != 1 {
		t.Errorf("expected version 1 after rollback, got %d", db.MigrationVersion())
	}
}

func TestPoolEmptyGet(t *testing.T) {
	pool := NewPool(5, 2, time.Hour)
	got := pool.Get()
	if got != nil {
		t.Error("expected nil from empty pool")
	}
}

func TestPoolOverflow(t *testing.T) {
	pool := NewPool(2, 1, time.Hour)

	conn1 := NewPostgreSQL()
	conn1.Connect(&Config{})
	pool.Put(conn1)

	conn2 := NewPostgreSQL()
	conn2.Connect(&Config{})
	pool.Put(conn2)

	if pool.Size() != 2 {
		t.Errorf("expected size 2, got %d", pool.Size())
	}

	conn3 := NewPostgreSQL()
	conn3.Connect(&Config{})
	pool.Put(conn3)

	if pool.Size() != 2 {
		t.Errorf("expected size 2 after overflow, got %d", pool.Size())
	}
}

func TestManagerGetNotFound(t *testing.T) {
	m := NewManager()
	_, ok := m.Get("nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

func TestManagerListEmpty(t *testing.T) {
	m := NewManager()
	names := m.List()
	if len(names) != 0 {
		t.Errorf("expected empty list, got %d", len(names))
	}
}

func TestManagerConnectAllMissingConfig(t *testing.T) {
	m := NewManager()
	pg := NewPostgreSQL()
	m.Register("pg", pg)

	configs := map[string]*Config{
		"other": {Host: "localhost"},
	}
	if err := m.ConnectAll(configs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestManagerCloseAllWithCloseError(t *testing.T) {
	m := NewManager()
	pg := NewPostgreSQL()
	m.Register("pg", pg)

	if err := m.CloseAll(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMySQLMigrationVersion(t *testing.T) {
	db := NewMySQL()
	db.Connect(&Config{})
	if db.MigrationVersion() != 0 {
		t.Errorf("expected version 0, got %d", db.MigrationVersion())
	}
}

func TestSQLiteMigrationVersion(t *testing.T) {
	db := NewSQLite()
	db.Connect(&Config{Database: ":memory:"})
	if db.MigrationVersion() != 0 {
		t.Errorf("expected version 0, got %d", db.MigrationVersion())
	}
}
