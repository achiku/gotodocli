package txmodel

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	txdb "github.com/achiku/pgtxdb"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib" // pgx
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func init() {
	txdb.Register("txdb", "pgx", "postgres://gotodo_api_test@localhost:5432/gotodo?sslmode=disable")
}

// TestCreateSchema set up test schema
func TestCreateSchema(cfg DBConfig, schema, user string) error {
	poolcfg := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     cfg.User,
			Database: cfg.DBName,
			Port:     5432,
		},
	}
	db, err := pgx.NewConnPool(poolcfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s AUTHORIZATION %s", schema, user))
	if err != nil {
		log.Printf("failed to create test schema: %s", schema)
		return err
	}
	return nil
}

// TestDropSchema set up test schema
func TestDropSchema(cfg DBConfig, schema string) error {
	poolcfg := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     cfg.User,
			Database: cfg.DBName,
			Port:     5432,
		},
	}
	db, err := pgx.NewConnPool(poolcfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", schema))
	if err != nil {
		log.Printf("failed to create test schema: %s", schema)
		return err
	}
	return nil
}

// TestCreateTables create test tables
func TestCreateTables(cfg DBConfig, path string) error {
	orgPwd, _ := os.Getwd()
	defer func() {
		os.Chdir(orgPwd)
	}()

	os.Chdir(path)
	ddl, err := ioutil.ReadFile("./ddl.sql")
	if err != nil {
		return errors.Wrap(err, "failed to read ddl.sql")
	}

	poolcfg := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     cfg.User,
			Database: cfg.DBName,
			Port:     5432,
		},
	}
	db, err := pgx.NewConnPool(poolcfg)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(ddl))
	if err != nil {
		return errors.Wrap(err, "failed to execute ddl.sql")
	}
	return nil
}

// TestSetupTx create tx and cleanup func for test
func TestSetupTx(t *testing.T) (Txer, func()) {
	db, err := sql.Open("txdb", uuid.NewV4().String())
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		tx.Rollback()
		db.Close()
	}
	return tx, cleanup
}

// TestSetupDB create db and cleanup func for test
func TestSetupDB(t *testing.T) (DBer, func()) {
	db, err := sql.Open("txdb", uuid.NewV4().String())
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		db.Close()
	}
	return db, cleanup
}
