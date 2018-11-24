package trancatemodel

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib" // pgx
	"github.com/pkg/errors"
)

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

// TestSetupDB create db
func TestSetupDB(t *testing.T) (DBer, func()) {
	db, err := sql.Open("pgx", "postgres://gotodo_api_test@localhost:5432/gotodo?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		_, err := db.Exec(`
		truncate table todo cascade;
		truncate table action cascade;
		`)
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
	return db, cleanup
}
